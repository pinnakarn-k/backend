package excel

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	tableGreenHex = "006937"
	lightGrayHex  = "969696"
)

type ExportFile struct {
	Title       string
	Subtitle    string
	GeneratedBy string

	InfoItems []InfoItem
	Columns   []Column
	Rows      [][]any
	Remark    []string
}

type InfoItem struct {
	LabelTH string
	LabelEN string
	Value   string
}

type Column struct {
	Header string
	Align  string // L, C, R

	// Fixed Excel column width.
	// Use only when a column must have a specific width.
	//
	// Example:
	// Width: 20
	//
	// If Width > 0, WidthRatio is ignored.
	Width float64

	// Relative width ratio.
	// Recommended for most reports.
	//
	// Example:
	// Date   = 15
	// Stock  = 20
	// Amount = 25
	//
	// Framework will automatically distribute
	// available report width based on ratio.
	WidthRatio float64
}

type Config struct {
	SheetName string
	FontName  string

	HeaderImagePath string
	HeaderRows      int

	// Approx A4 portrait printable width in pixels.
	A4PrintableWidthPx int

	OpenPassword string
}

func DefaultConfig() Config {
	return Config{
		SheetName:          "Report",
		FontName:           "TH Sarabun New",
		HeaderRows:         4,
		A4PrintableWidthPx: 746,
	}
}

func Generate(file ExportFile) ([]byte, error) {
	return GenerateWithConfig(file, DefaultConfig())
}

func GenerateWithConfig(file ExportFile, cfg Config) ([]byte, error) {
	cfg = normalizeConfig(cfg)

	f := excelize.NewFile()
	defer f.Close()

	defaultSheet := f.GetSheetName(0)
	if err := f.SetSheetName(defaultSheet, cfg.SheetName); err != nil {
		return nil, fmt.Errorf("set sheet name: %w", err)
	}

	sheet := cfg.SheetName

	lastColName, err := excelize.ColumnNumberToName(max(1, len(file.Columns)))
	if err != nil {
		return nil, err
	}

	if err := setupA4Page(f, sheet); err != nil {
		return nil, err
	}

	if err := setupColumnWidths(f, sheet, file, cfg); err != nil {
		return nil, err
	}

	if err := drawHeaderImage(f, sheet, file, cfg, lastColName); err != nil {
		return nil, err
	}

	currentRow := cfg.HeaderRows + 2

	if file.Title != "" {
		cell := fmt.Sprintf("A%d", currentRow)
		endCell := fmt.Sprintf("%s%d", lastColName, currentRow)

		_ = f.MergeCell(sheet, cell, endCell)
		_ = f.SetCellValue(sheet, cell, file.Title)

		style, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold:   true,
				Size:   18,
				Family: cfg.FontName,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		if err != nil {
			return nil, err
		}

		_ = f.SetCellStyle(sheet, cell, endCell, style)
		currentRow += 2
	}

	if len(file.InfoItems) > 0 {
		nextRow, err := drawInfoBlock(f, sheet, file, cfg, currentRow, lastColName)
		if err != nil {
			return nil, err
		}
		currentRow = nextRow + 1
	}

	tableStartRow := currentRow

	if err := drawTable(f, sheet, file, cfg, tableStartRow); err != nil {
		return nil, err
	}

	remarkStartRow := tableStartRow + len(file.Rows) + 2
	if err := drawRemark(f, sheet, file, cfg, remarkStartRow, lastColName); err != nil {
		return nil, err
	}

	if err := setupPrintFooter(f, sheet); err != nil {
		return nil, err
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write excel buffer: %w", err)
	}

	data := buf.Bytes()

	if cfg.OpenPassword != "" {
		return encryptExcel(data, cfg.OpenPassword)
	}

	return data, nil
}

func normalizeConfig(cfg Config) Config {
	if cfg.SheetName == "" {
		cfg.SheetName = "Report"
	}
	if cfg.FontName == "" {
		cfg.FontName = "TH Sarabun New"
	}
	if cfg.HeaderRows <= 0 {
		cfg.HeaderRows = 4
	}
	if cfg.A4PrintableWidthPx <= 0 {
		cfg.A4PrintableWidthPx = 746
	}

	return cfg
}

func setupA4Page(f *excelize.File, sheet string) error {
	orientation := "portrait"
	paperSize := 9 // A4
	fitToWidth := 1
	fitToHeight := 0

	if err := f.SetPageLayout(sheet, &excelize.PageLayoutOptions{
		Orientation: &orientation,
		Size:        &paperSize,
		FitToWidth:  &fitToWidth,
		FitToHeight: &fitToHeight,
	}); err != nil {
		return fmt.Errorf("set page layout: %w", err)
	}

	left := 0.25
	right := 0.25
	top := 0.35
	bottom := 0.35
	header := 0.2
	footer := 0.2

	return f.SetPageMargins(sheet, &excelize.PageLayoutMarginsOptions{
		Left:   &left,
		Right:  &right,
		Top:    &top,
		Bottom: &bottom,
		Header: &header,
		Footer: &footer,
	})
}

func drawHeaderImage(
	f *excelize.File,
	sheet string,
	file ExportFile,
	cfg Config,
	lastColName string,
) error {
	if cfg.HeaderImagePath == "" || len(file.Columns) == 0 {
		return nil
	}

	imgFile, err := os.Open(cfg.HeaderImagePath)
	if err != nil {
		return fmt.Errorf("open header image: %w", err)
	}
	defer imgFile.Close()

	imgCfg, _, err := image.DecodeConfig(imgFile)
	if err != nil {
		return fmt.Errorf("decode header image: %w", err)
	}

	targetWidthPx := reportWidthPixels(file, cfg)
	scale := float64(targetWidthPx) / float64(imgCfg.Width)
	targetHeightPt := float64(imgCfg.Height) * scale * 0.75

	rowHeight := targetHeightPt / float64(cfg.HeaderRows)
	for r := 1; r <= cfg.HeaderRows; r++ {
		if err := f.SetRowHeight(sheet, r, rowHeight); err != nil {
			return err
		}
	}

	endCell := fmt.Sprintf("%s%d", lastColName, cfg.HeaderRows)
	if err := f.MergeCell(sheet, "A1", endCell); err != nil {
		return err
	}

	if err := f.AddPicture(sheet, "A1", cfg.HeaderImagePath, &excelize.GraphicOptions{
		OffsetX:         0,
		OffsetY:         0,
		ScaleX:          scale,
		ScaleY:          scale,
		LockAspectRatio: true,
	}); err != nil {
		return fmt.Errorf("add header image: %w", err)
	}

	return nil
}

func drawInfoBlock(
	f *excelize.File,
	sheet string,
	file ExportFile,
	cfg Config,
	startRow int,
	lastColName string,
) (int, error) {
	labelStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12, Family: cfg.FontName},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
			WrapText:   true,
		},
	})
	if err != nil {
		return 0, err
	}

	valueStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12, Family: cfg.FontName},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
	})
	if err != nil {
		return 0, err
	}

	row := startRow

	for _, item := range file.InfoItems {
		labelCell := fmt.Sprintf("A%d", row)
		valueCell := fmt.Sprintf("B%d", row)

		_ = f.SetCellValue(sheet, labelCell, item.LabelTH+"\n"+item.LabelEN)
		_ = f.SetCellValue(sheet, valueCell, item.Value)
		_ = f.SetRowHeight(sheet, row, 28)
		_ = f.SetCellStyle(sheet, labelCell, labelCell, labelStyle)
		_ = f.SetCellStyle(sheet, valueCell, valueCell, valueStyle)

		row++
	}

	printDateLabelCell := fmt.Sprintf("%s%d", lastColName, startRow)
	printDateValueCell := fmt.Sprintf("%s%d", lastColName, startRow+1)

	_ = f.SetCellValue(sheet, printDateLabelCell, "วันที่พิมพ์ / Print Date")
	_ = f.SetCellValue(sheet, printDateValueCell, time.Now().Format("02/01/2006 15:04:05"))

	rightStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12, Family: cfg.FontName},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
	})
	if err != nil {
		return 0, err
	}

	_ = f.SetCellStyle(sheet, printDateLabelCell, printDateValueCell, rightStyle)

	return row, nil
}

func drawTable(
	f *excelize.File,
	sheet string,
	file ExportFile,
	cfg Config,
	startRow int,
) error {
	if len(file.Columns) == 0 {
		return nil
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Family: cfg.FontName,
			Color:  "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{tableGreenHex},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: borderStyle("BFBFBF"),
	})
	if err != nil {
		return err
	}

	bodyStyles := map[string]int{}

	for _, align := range []string{"L", "C", "R"} {
		styleID, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size:   12,
				Family: cfg.FontName,
			},
			Alignment: &excelize.Alignment{
				Horizontal: alignToExcel(align),
				Vertical:   "center",
			},
			Border: borderStyle("D9D9D9"),
		})
		if err != nil {
			return err
		}
		bodyStyles[align] = styleID
	}

	for i, col := range file.Columns {
		cell, err := excelize.CoordinatesToCellName(i+1, startRow)
		if err != nil {
			return err
		}

		_ = f.SetCellValue(sheet, cell, col.Header)
		_ = f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	_ = f.SetRowHeight(sheet, startRow, 28)

	for r, row := range file.Rows {
		excelRow := startRow + 1 + r

		for c := range file.Columns {
			cell, err := excelize.CoordinatesToCellName(c+1, excelRow)
			if err != nil {
				return err
			}

			value := ""
			if c < len(row) {
				value = fmt.Sprint(row[c])
			}

			align := file.Columns[c].Align
			if align == "" {
				align = "C"
			}

			styleID, ok := bodyStyles[align]
			if !ok {
				styleID = bodyStyles["C"]
			}

			_ = f.SetCellValue(sheet, cell, value)
			_ = f.SetCellStyle(sheet, cell, cell, styleID)
		}

		_ = f.SetRowHeight(sheet, excelRow, 22)
	}

	return nil
}

func drawRemark(
	f *excelize.File,
	sheet string,
	file ExportFile,
	cfg Config,
	startRow int,
	lastColName string,
) error {
	if len(file.Remark) == 0 {
		return nil
	}

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: cfg.FontName,
			Color:  lightGrayHex,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "top",
			WrapText:   true,
		},
	})
	if err != nil {
		return err
	}

	for i, remark := range file.Remark {
		row := startRow + i
		startCell := fmt.Sprintf("A%d", row)
		endCell := fmt.Sprintf("%s%d", lastColName, row)

		_ = f.MergeCell(sheet, startCell, endCell)
		_ = f.SetCellValue(sheet, startCell, remark)
		_ = f.SetCellStyle(sheet, startCell, endCell, style)
		_ = f.SetRowHeight(sheet, row, 20)
	}

	return nil
}

func setupPrintFooter(f *excelize.File, sheet string) error {
	return f.SetHeaderFooter(sheet, &excelize.HeaderFooterOptions{
		OddFooter:  "&RPage &P of &N",
		EvenFooter: "&RPage &P of &N",
	})
}

func setupColumnWidths(
	f *excelize.File,
	sheet string,
	file ExportFile,
	cfg Config,
) error {
	widths := calculateColumnWidths(file, cfg)

	for i, width := range widths {
		colName, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return err
		}

		if width < 8 {
			width = 8
		}

		if err := f.SetColWidth(sheet, colName, colName, width); err != nil {
			return err
		}
	}

	return nil
}

func calculateColumnWidths(file ExportFile, cfg Config) []float64 {
	widths := make([]float64, len(file.Columns))

	totalAvailableWidth := pixelsToExcelTotalWidth(cfg.A4PrintableWidthPx)

	totalFixedWidth := 0.0
	totalRatio := 0.0

	for _, col := range file.Columns {
		if col.Width > 0 {
			totalFixedWidth += col.Width
			continue
		}

		ratio := col.WidthRatio
		if ratio <= 0 {
			ratio = 1
		}

		totalRatio += ratio
	}

	remainingWidth := totalAvailableWidth - totalFixedWidth
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	for i, col := range file.Columns {
		if col.Width > 0 {
			widths[i] = col.Width
			continue
		}

		ratio := col.WidthRatio
		if ratio <= 0 {
			ratio = 1
		}

		widths[i] = remainingWidth * ratio / totalRatio
	}

	return widths
}

func reportWidthPixels(file ExportFile, cfg Config) int {
	widths := calculateColumnWidths(file, cfg)

	totalExcelWidth := 0.0
	for _, w := range widths {
		totalExcelWidth += w
	}

	return excelTotalWidthToPixels(totalExcelWidth)
}

func pixelsToExcelTotalWidth(px int) float64 {
	if px <= 0 {
		return 100
	}

	return float64(px-5) / 7
}

func excelTotalWidthToPixels(width float64) int {
	if width <= 0 {
		return 0
	}

	return int(width*7 + 5)
}

func borderStyle(color string) []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: color, Style: 1},
		{Type: "top", Color: color, Style: 1},
		{Type: "right", Color: color, Style: 1},
		{Type: "bottom", Color: color, Style: 1},
	}
}

func alignToExcel(align string) string {
	switch align {
	case "L":
		return "left"
	case "R":
		return "right"
	default:
		return "center"
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
