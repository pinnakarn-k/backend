package excel

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExportFile struct {
	SheetName string
	Headers   []string
	Rows      [][]any
}

func Generate(file ExportFile) ([]byte, error) {
	if file.SheetName == "" {
		file.SheetName = "Report"
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	defaultSheet := "Sheet1"

	if err := f.SetSheetName(defaultSheet, file.SheetName); err != nil {
		return nil, fmt.Errorf("set sheet name: %w", err)
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create header style: %w", err)
	}

	for colIndex, header := range file.Headers {
		cell, err := excelize.CoordinatesToCellName(colIndex+1, 1)
		if err != nil {
			return nil, fmt.Errorf("header cell name: %w", err)
		}

		if err := f.SetCellValue(file.SheetName, cell, header); err != nil {
			return nil, fmt.Errorf("set header value: %w", err)
		}

		if err := f.SetCellStyle(file.SheetName, cell, cell, headerStyle); err != nil {
			return nil, fmt.Errorf("set header style: %w", err)
		}
	}

	for rowIndex, row := range file.Rows {
		excelRow := rowIndex + 2

		for colIndex, value := range row {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, excelRow)
			if err != nil {
				return nil, fmt.Errorf("data cell name: %w", err)
			}

			if err := f.SetCellValue(file.SheetName, cell, value); err != nil {
				return nil, fmt.Errorf("set data value: %w", err)
			}
		}
	}

	for colIndex := range file.Headers {
		colName, err := excelize.ColumnNumberToName(colIndex + 1)
		if err != nil {
			return nil, fmt.Errorf("column name: %w", err)
		}

		_ = f.SetColWidth(file.SheetName, colName, colName, 20)
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write excel buffer: %w", err)
	}

	return buffer.Bytes(), nil
}

func GenerateReader(file ExportFile) (*bytes.Reader, error) {
	data, err := Generate(file)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

// ตัวอย่างใช้กับเมนู HD
// rows := make([][]any, 0, len(items))

// for _, item := range items {
// 	rows = append(rows, []any{
// 		item.AccountNo,
// 		item.CustomerName,
// 		item.Symbol,
// 		item.Quantity,
// 		item.Price,
// 		item.Amount,
// 	})
// }

// excelBytes, err := excel.Generate(excel.ExportFile{
// 	SheetName: "Historical Data",
// 	Headers: []string{
// 		"เลขที่บัญชี",
// 		"ชื่อลูกค้า",
// 		"หลักทรัพย์",
// 		"จำนวน",
// 		"ราคา",
// 		"มูลค่า",
// 	},
// 	Rows: rows,
// })
// if err != nil {
// 	return err
// }

// ส่ง download ผ่าน Fiber
// c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
// c.Set("Content-Disposition", `attachment; filename="report.xlsx"`)

// return c.Send(excelBytes)
