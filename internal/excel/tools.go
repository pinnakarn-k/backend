package excel

// import "backend/internal/excel"

// func (s *Service) generateExcel(
// 	req ExportRequest,
// 	columns []excel.Column,
// 	rows [][]any,
// 	fileName string,
// ) (FileResult, error) {
// 	data, err := excel.GenerateWithConfig(
// 		excel.ExportFile{
// 			Title:       "Transaction History",
// 			Subtitle:    "รายการเคลื่อนไหวย้อนหลัง",
// 			GeneratedBy: "system",
// 			InfoItems: []excel.InfoItem{
// 				{LabelTH: "ประเภทรายการ :", LabelEN: "Transaction Type", Value: "ซื้อ-ขายหลักทรัพย์"},
// 				{LabelTH: "เลขบัญชีซื้อขายหลักทรัพย์ :", LabelEN: "Account No", Value: req.AccountNo},
// 				{LabelTH: "ช่วงเวลา :", LabelEN: "Period", Value: req.FromDate + " - " + req.ToDate},
// 				{LabelTH: "เงื่อนไขเพิ่มเติม :", LabelEN: "Additional Conditions", Value: "-"},
// 			},
// 			Columns: columns,
// 			Rows:    rows,
// 			Remark: []string{
// 				"* จำนวนเงิน คือ มูลค่าซื้อหรือขายสุทธิ ที่ถูกคำนวณค่าธรรมเนียมต่างๆ และภาษีมูลค่าเพิ่มเรียบร้อยแล้ว",
// 				"Amount represents the net transaction value after fees and Value-Added Tax have been applied.",
// 			},
// 		},
// 		excel.Config{
// 			SheetName:          "Report",
// 			FontName:           "TH Sarabun New",
// 			HeaderImagePath:    "assets/header.png",
// 			HeaderRows:         4,
// 			A4PrintableWidthPx: 746,
// 			OpenPassword:       "123456",
// 		},
// 	)
// 	if err != nil {
// 		return FileResult{}, err
// 	}

// 	return FileResult{
// 		FileName:    fileName,
// 		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
// 		Bytes:       data,
// 	}, nil
// }
