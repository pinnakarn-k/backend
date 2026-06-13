package pdf

import (
	"bytes"
	"fmt"

	"github.com/phpdave11/gofpdf"
)

const (
	tableGreenR = 0
	tableGreenG = 105
	tableGreenB = 55
)

type ExportFile struct {
	Title       string
	Subtitle    string
	GeneratedBy string

	InfoItems []InfoItem

	Columns []Column
	Rows    [][]any

	Remark []string
}

type InfoItem struct {
	LabelTH string
	LabelEN string
	Value   string
}

type Column struct {
	Header string
	Align  string // L, C, R

	// Fixed PDF column width in current unit, normally "mm".
	// Use only when a column needs a specific width.
	// If Width > 0, WidthRatio is ignored.
	Width float64

	// Recommended for most reports.
	// Framework distributes remaining table width by ratio.
	WidthRatio float64
}

type Config struct {
	FontPath string
	FontName string

	HeaderImagePath string
	HeaderImageType string // PNG, JPG

	FooterText string

	Orientation string
	PageSize    string
	Unit        string

	OpenPassword string
}

func DefaultConfig() Config {
	return Config{
		FontPath:        "fonts/THSarabunNew.ttf",
		FontName:        "THSarabun",
		HeaderImagePath: "",
		HeaderImageType: "PNG",
		FooterText:      "footer text",
		Orientation:     "P",
		PageSize:        "A4",
		Unit:            "mm",
	}
}

func Generate(file ExportFile) ([]byte, error) {
	return GenerateWithConfig(file, DefaultConfig())
}

func GenerateWithConfig(file ExportFile, cfg Config) ([]byte, error) {
	cfg = normalizeConfig(cfg)

	p := gofpdf.New(cfg.Orientation, cfg.Unit, cfg.PageSize, "")
	p.AddUTF8Font(cfg.FontName, "", cfg.FontPath)
	p.SetFont(cfg.FontName, "", 14)

	p.SetMargins(10, 10, 10)
	p.SetAutoPageBreak(false, 10)

	p.AddPage()

	drawFirstPageHeader(p, file, cfg)
	drawInfoBlock(p, file, cfg)
	drawTable(p, file, cfg)
	drawRemark(p, file, cfg)
	drawFooter(p, cfg)

	var buf bytes.Buffer
	if err := p.Output(&buf); err != nil {
		return nil, fmt.Errorf("output pdf: %w", err)
	}

	pdfBytes := buf.Bytes()

	if cfg.OpenPassword != "" {
		return encryptPDF(pdfBytes, cfg.OpenPassword)
	}

	return pdfBytes, nil
}

func normalizeConfig(cfg Config) Config {
	if cfg.FontName == "" {
		cfg.FontName = "THSarabun"
	}
	if cfg.FontPath == "" {
		cfg.FontPath = "fonts/THSarabunNew.ttf"
	}
	if cfg.HeaderImageType == "" {
		cfg.HeaderImageType = "PNG"
	}
	if cfg.FooterText == "" {
		cfg.FooterText = "footer text"
	}
	if cfg.Orientation == "" {
		cfg.Orientation = "P"
	}
	if cfg.PageSize == "" {
		cfg.PageSize = "A4"
	}
	if cfg.Unit == "" {
		cfg.Unit = "mm"
	}

	return cfg
}

func drawFirstPageHeader(
	p *gofpdf.Fpdf,
	file ExportFile,
	cfg Config,
) {
	if cfg.HeaderImagePath != "" {
		pageWidth, _ := p.GetPageSize()
		left, _, right, _ := p.GetMargins()

		imageWidth := pageWidth - left - right
		imageHeight := 32.0

		p.ImageOptions(
			cfg.HeaderImagePath,
			left,
			10,
			imageWidth,
			imageHeight,
			false,
			gofpdf.ImageOptions{
				ImageType: cfg.HeaderImageType,
				ReadDpi:   true,
			},
			0,
			"",
		)

		p.SetY(10 + imageHeight + 6)
		return
	}

	p.SetTextColor(0, 0, 0)

	if file.Title != "" {
		p.SetFont(cfg.FontName, "", 22)
		p.CellFormat(0, 10, file.Title, "", 1, "L", false, 0, "")
	}

	if file.Subtitle != "" {
		p.SetFont(cfg.FontName, "", 16)
		p.CellFormat(0, 8, file.Subtitle, "", 1, "L", false, 0, "")
	}

	p.Ln(4)
}

func drawInfoBlock(
	p *gofpdf.Fpdf,
	file ExportFile,
	cfg Config,
) {
	if len(file.InfoItems) == 0 {
		return
	}

	p.SetTextColor(0, 0, 0)

	labelX := 18.0
	valueX := 60.0
	y := p.GetY() + 4

	for i, item := range file.InfoItems {
		yy := y + float64(i*13)

		p.SetFont(cfg.FontName, "", 12)
		p.SetXY(labelX, yy)
		p.CellFormat(38, 5, item.LabelTH, "", 0, "R", false, 0, "")

		p.SetFont(cfg.FontName, "", 10)
		p.SetXY(labelX, yy+5)
		p.CellFormat(38, 4, item.LabelEN, "", 0, "R", false, 0, "")

		p.SetFont(cfg.FontName, "", 12)
		p.SetXY(valueX, yy)
		p.CellFormat(75, 5, item.Value, "", 0, "L", false, 0, "")
	}

	rightLabelX := 145.0
	rightValueX := 170.0

	p.SetFont(cfg.FontName, "", 12)
	p.SetXY(rightLabelX, y+26)
	p.CellFormat(22, 5, "วันที่พิมพ์ :", "", 0, "R", false, 0, "")

	p.SetFont(cfg.FontName, "", 10)
	p.SetXY(rightLabelX, y+31)
	p.CellFormat(22, 4, "Print Date", "", 0, "R", false, 0, "")

	p.SetFont(cfg.FontName, "", 12)
	p.SetXY(rightValueX, y+26)
	p.CellFormat(35, 5, "18/05/2567 10:30:45", "", 0, "L", false, 0, "")

	p.SetY(y + float64(len(file.InfoItems)*13) + 6)
}

func drawTable(
	p *gofpdf.Fpdf,
	file ExportFile,
	cfg Config,
) {
	if len(file.Columns) == 0 {
		return
	}

	colWidths := calculateColumnWidths(p, file.Columns)

	drawTableHeader(p, file.Columns, colWidths, cfg)

	for _, row := range file.Rows {
		rowHeight := 7.0

		if shouldAddPage(p, rowHeight) {
			p.AddPage()
			p.SetY(18)
			drawTableHeader(p, file.Columns, colWidths, cfg)
		}

		drawTableRow(p, row, file.Columns, colWidths, rowHeight, cfg)
	}
}

func drawTableHeader(
	p *gofpdf.Fpdf,
	columns []Column,
	colWidths []float64,
	cfg Config,
) {
	p.SetFont(cfg.FontName, "", 13)
	p.SetFillColor(tableGreenR, tableGreenG, tableGreenB)
	p.SetTextColor(255, 255, 255)
	p.SetDrawColor(120, 150, 130)

	for i, col := range columns {
		p.CellFormat(colWidths[i], 12, col.Header, "1", 0, "C", true, 0, "")
	}

	p.Ln(-1)

	p.SetTextColor(0, 0, 0)
	p.SetDrawColor(220, 220, 220)
}

func drawTableRow(
	p *gofpdf.Fpdf,
	row []any,
	columns []Column,
	colWidths []float64,
	rowHeight float64,
	cfg Config,
) {
	p.SetFont(cfg.FontName, "", 12)
	p.SetTextColor(0, 0, 0)
	p.SetDrawColor(220, 220, 220)

	for i := range columns {
		value := ""

		if i < len(row) {
			value = fmt.Sprint(row[i])
		}

		align := columns[i].Align
		if align == "" {
			align = "C"
		}

		p.CellFormat(colWidths[i], rowHeight, value, "1", 0, align, false, 0, "")
	}

	p.Ln(-1)
}

func calculateColumnWidths(
	p *gofpdf.Fpdf,
	columns []Column,
) []float64 {
	pageWidth, _ := p.GetPageSize()
	left, _, right, _ := p.GetMargins()

	usableWidth := pageWidth - left - right

	widths := make([]float64, len(columns))

	totalFixedWidth := 0.0
	totalRatio := 0.0

	for _, col := range columns {
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

	remainingWidth := usableWidth - totalFixedWidth
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	for i, col := range columns {
		if col.Width > 0 {
			widths[i] = col.Width
			continue
		}

		ratio := col.WidthRatio
		if ratio <= 0 {
			ratio = 1
		}

		if totalRatio <= 0 {
			widths[i] = remainingWidth / float64(len(columns))
		} else {
			widths[i] = remainingWidth * ratio / totalRatio
		}
	}

	return widths
}

func drawRemark(
	p *gofpdf.Fpdf,
	file ExportFile,
	cfg Config,
) {
	if len(file.Remark) == 0 {
		return
	}

	p.Ln(4)

	p.SetFont(cfg.FontName, "", 10)
	p.SetTextColor(150, 150, 150)

	for _, remark := range file.Remark {
		if shouldAddPage(p, 8) {
			p.AddPage()
		}

		p.MultiCell(
			0,
			5,
			remark,
			"",
			"L",
			false,
		)
	}

	p.SetTextColor(0, 0, 0)
}

func shouldAddPage(
	p *gofpdf.Fpdf,
	nextRowHeight float64,
) bool {
	_, pageHeight := p.GetPageSize()
	_, _, _, bottom := p.GetMargins()

	footerSpace := 20.0

	return p.GetY()+nextRowHeight > pageHeight-bottom-footerSpace
}

func drawFooter(
	p *gofpdf.Fpdf,
	cfg Config,
) {
	totalPages := p.PageCount()

	for page := 1; page <= totalPages; page++ {
		p.SetPage(page)

		pageWidth, pageHeight := p.GetPageSize()
		left, _, right, _ := p.GetMargins()

		footerHeight := 9.0
		footerY := pageHeight - 16
		footerWidth := pageWidth - left - right

		p.SetFillColor(230, 230, 230)
		p.Rect(left, footerY, footerWidth, footerHeight, "F")

		p.SetFont(cfg.FontName, "", 10)
		p.SetTextColor(90, 90, 90)
		p.SetXY(left, footerY)
		p.CellFormat(footerWidth, footerHeight, cfg.FooterText, "", 0, "C", false, 0, "")

		p.SetFont(cfg.FontName, "", 11)
		p.SetTextColor(0, 0, 0)
		p.SetXY(left, footerY)
		p.CellFormat(
			footerWidth-4,
			footerHeight,
			fmt.Sprintf("Page %d of %d", page, totalPages),
			"",
			0,
			"R",
			false,
			0,
			"",
		)
	}
}
