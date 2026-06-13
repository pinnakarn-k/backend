package main

import (
	"os"
	"os/signal"
	"syscall"

	"backend/internal/app"
	"backend/internal/config"
	"backend/internal/logger"
	redisinfra "backend/internal/redis"
)

func main() {

	// test excels
	// file := excel.ExportFile{
	// 	SheetName: "Portfolio",

	// 	Headers: []string{
	// 		"Account No",
	// 		"Customer",
	// 		"Symbol",
	// 		"Quantity",
	// 		"Price",
	// 		"Amount",
	// 	},

	// 	Rows: [][]any{
	// 		{
	// 			"123-6",
	// 			"สมชาย ใจดี",
	// 			"AOT",
	// 			100,
	// 			62.50,
	// 			6250.00,
	// 		},
	// 		{
	// 			"456-7",
	// 			"สมหญิง รักดี",
	// 			"PTT",
	// 			50,
	// 			30.00,
	// 			1500.00,
	// 		},
	// 		{
	// 			"789-0",
	// 			"บริษัท ABC จำกัด",
	// 			"CPALL",
	// 			1000,
	// 			55.25,
	// 			55250.00,
	// 		},
	// 	},
	// }

	// data, err := excel.Generate(file)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = os.WriteFile(
	// 	"report.xlsx",
	// 	data,
	// 	0644,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("report.xlsx generated")
	// test excel

	// test excel 2
	// rows := make([][]any, 0)

	// samples := [][]any{
	// 	{"02/01/2567", "AOT", "ซื้อ", 100, "62.50", "6,250.00"},
	// 	{"02/01/2567", "PTT", "ขาย", 200, "34.25", "6,850.00"},
	// 	{"03/01/2567", "DELTA", "ซื้อ", 50, "78.75", "3,937.50"},
	// 	{"03/01/2567", "BDMS", "ซื้อ", 150, "28.50", "4,275.00"},
	// 	{"04/01/2567", "CPALL", "ขาย", 100, "55.25", "5,525.00"},
	// 	{"04/01/2567", "KBANK", "ซื้อ", 80, "130.00", "10,400.00"},
	// 	{"05/01/2567", "SCB", "ขาย", 120, "102.00", "12,240.00"},
	// 	{"05/01/2567", "ADVANC", "ซื้อ", 30, "221.00", "6,630.00"},
	// 	{"08/01/2567", "TRUE", "ขาย", 300, "7.90", "2,370.00"},
	// 	{"08/01/2567", "AOT", "ซื้อ", 50, "63.00", "3,150.00"},
	// }

	// for i := 0; i < 80; i++ {
	// 	rows = append(rows, samples[i%len(samples)])
	// }

	// data, err := excel.GenerateWithConfig(
	// 	excel.ExportFile{
	// 		Title:       "Transaction History",
	// 		Subtitle:    "รายการเคลื่อนไหวย้อนหลัง",
	// 		GeneratedBy: "system",

	// 		InfoItems: []excel.InfoItem{
	// 			{LabelTH: "ประเภทรายการ :", LabelEN: "Transaction Type", Value: "ซื้อ-ขายหลักทรัพย์"},
	// 			{LabelTH: "เลขบัญชีซื้อขายหลักทรัพย์ :", LabelEN: "Account No", Value: "123-4-56789-0"},
	// 			{LabelTH: "ช่วงเวลา :", LabelEN: "Period", Value: "01/01/2567 - 31/01/2567"},
	// 			{LabelTH: "เงื่อนไขเพิ่มเติม :", LabelEN: "Additional Conditions", Value: "ไม่มี"},
	// 		},

	// 		// Columns: []excel.Column{
	// 		// 	{
	// 		// 		Header: "Account No",
	// 		// 		Width: 18,
	// 		// 	},
	// 		// 	{
	// 		// 		Header: "Customer Name",
	// 		// 		WidthRatio: 40,
	// 		// 	},
	// 		// 	{
	// 		// 		Header: "Remark",
	// 		// 		WidthRatio: 60,
	// 		// 	},
	// 		// }

	// 		Columns: []excel.Column{
	// 			{Header: "วันที่ / Date", WidthRatio: 15, Align: "C"},
	// 			{Header: "หลักทรัพย์ / Stock", WidthRatio: 20, Align: "L"},
	// 			{Header: "ซื้อ-ขาย/ Buy-Sell", WidthRatio: 14, Align: "C"},
	// 			{Header: "จำนวนหุ้น/ Volume", WidthRatio: 15, Align: "R"},
	// 			{Header: "ราคา/ Price", WidthRatio: 13, Align: "R"},
	// 			{Header: "จำนวนเงิน/ Amount", WidthRatio: 23, Align: "R"},
	// 		},

	// 		Rows: rows,

	// 		Remark: []string{
	// 			"* จำนวนเงิน คือ มูลค่าซื้อหรือขายสุทธิ ที่ถูกคำนวณค่าธรรมเนียมต่างๆ และภาษีมูลค่าเพิ่มเรียบร้อยแล้ว",
	// 			"Amount represents the net transaction value after fees and Value-Added Tax have been applied.",
	// 		},
	// 	},
	// 	excel.Config{
	// 		SheetName:          "Report",
	// 		FontName:           "TH Sarabun New",
	// 		HeaderImagePath:    "assets/header.png",
	// 		HeaderRows:         4,
	// 		A4PrintableWidthPx: 746,
	// 		OpenPassword:       "123456",
	// 	},
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := os.WriteFile("report.xlsx", data, 0644); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("report.xlsx generated")
	// test excel 2

	// test pdf 2
	// rows := make([][]any, 0)

	// samples := [][]any{
	// 	{"02/01/2567", "x", "ซื้อ", 100, "62.50", "6,250.00"},
	// 	{"02/01/2567", "x", "ขาย", 200, "34.25", "6,850.00"},
	// 	{"03/01/2567", "x", "ซื้อ", 50, "78.75", "3,937.50"},
	// 	{"03/01/2567", "x", "ซื้อ", 150, "28.50", "4,275.00"},
	// 	{"04/01/2567", "x", "ขาย", 100, "55.25", "5,525.00"},
	// 	{"04/01/2567", "x", "ซื้อ", 80, "130.00", "10,400.00"},
	// 	{"05/01/2567", "x", "ขาย", 120, "102.00", "12,240.00"},
	// 	{"05/01/2567", "x", "ซื้อ", 30, "221.00", "6,630.00"},
	// 	{"08/01/2567", "x", "ขาย", 300, "7.90", "2,370.00"},
	// 	{"08/01/2567", "x", "ซื้อ", 50, "63.00", "3,150.00"},
	// }

	// for i := 0; i < 80; i++ {
	// 	rows = append(rows, samples[i%len(samples)])
	// }

	// data, err := pdf.GenerateWithConfig(
	// 	pdf.ExportFile{
	// 		Title:       "Transaction History",
	// 		Subtitle:    "รายการเคลื่อนไหวย้อนหลัง",
	// 		GeneratedBy: "system",

	// 		InfoItems: []pdf.InfoItem{
	// 			{LabelTH: "ประเภทรายการ :", LabelEN: "Transaction Type", Value: "ซื้อ-ขายหลักทรัพย์"},
	// 			{LabelTH: "เลขบัญชีซื้อขายหลักทรัพย์ :", LabelEN: "Account No", Value: "123-4-56789-0"},
	// 			{LabelTH: "ช่วงเวลา :", LabelEN: "Period", Value: "01/01/2567 - 31/01/2567"},
	// 			{LabelTH: "เงื่อนไขเพิ่มเติม :", LabelEN: "Additional Conditions", Value: "ไม่มี"},
	// 		},

	// 		Columns: []pdf.Column{
	// 			{Header: "วันที่ / Date", WidthRatio: 15, Align: "C"},
	// 			{Header: "หลักทรัพย์ / Stock", WidthRatio: 20, Align: "L"},
	// 			{Header: "ซื้อ-ขาย/ Buy-Sell", WidthRatio: 14, Align: "C"},
	// 			{Header: "จำนวนหุ้น/ Volume", WidthRatio: 15, Align: "R"},
	// 			{Header: "ราคา/ Price", WidthRatio: 13, Align: "R"},
	// 			{Header: "จำนวนเงิน/ Amount", WidthRatio: 23, Align: "R"},
	// 		},

	// 		//Columns: []pdf.Column{
	// 		// 	{Header: "Account No", Width: 28, Align: "C"},
	// 		// 	{Header: "Customer Name", WidthRatio: 40, Align: "L"},
	// 		// 	{Header: "Amount", WidthRatio: 20, Align: "R"},
	// 		// }

	// 		Rows: rows,

	// 		Remark: []string{
	// 			"* จำนวนเงิน คือ มูลค่าซื้อหรือขายสุทธิ ที่ถูกคำนวณค่าธรรมเนียมต่างๆ และภาษีมูลค่าเพิ่มเรียบร้อยแล้ว",
	// 			"Amount represents the net transaction value after fees and Value-Added Tax have been applied.",
	// 		},
	// 	},
	// 	pdf.Config{
	// 		FontPath:        "fonts/THSarabunNew.ttf",
	// 		FontName:        "THSarabun",
	// 		HeaderImagePath: "assets/header.png",
	// 		HeaderImageType: "PNG",
	// 		FooterText:      "footer text",
	// 		Orientation:     "P",
	// 		PageSize:        "A4",
	// 		Unit:            "mm",
	// 		OpenPassword:    "123456",
	// 	},
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := os.WriteFile("report.pdf", data, 0644); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("report.pdf generated")
	// test pdf 2

	// // test pdf
	// rows := make([][]any, 0)

	// for i := 1; i <= 120; i++ {
	// 	rows = append(rows, []any{
	// 		i,
	// 		"สมชาย ใจดี",
	// 		"123-6",
	// 		"AOT",
	// 		100,
	// 		6250.00,
	// 	})
	// }

	// data, err := pdf.GenerateWithConfig(
	// 	pdf.ExportFile{
	// 		Title:       "Historical Data",
	// 		Subtitle:    "รายงานข้อมูลย้อนหลัง",
	// 		GeneratedBy: "system",
	// 		Headers: []string{
	// 			"ลำดับ",
	// 			"ชื่อลูกค้า",
	// 			"เลขที่บัญชี",
	// 			"หลักทรัพย์",
	// 			"จำนวน",
	// 			"มูลค่า",
	// 		},
	// 		Rows: rows,
	// 	},
	// 	pdf.Config{
	// 		FontPath:        "fonts/THSarabunNew.ttf",
	// 		FontName:        "THSarabun",
	// 		HeaderImagePath: "assets/header.png",
	// 		HeaderImageType: "PNG",
	// 	},
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := os.WriteFile("report.pdf", data, 0644); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("report.pdf generated")
	// // test pdf

	cfg := config.Load()

	logger := logger.New(
		cfg.Service,
		cfg.Env,
	)

	redisClient, err := redisinfra.New(
		redisinfra.Config{
			Enabled:  cfg.RedisEnabled,
			Host:     cfg.RedisHost,
			Port:     cfg.RedisPort,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		},
	)
	if err != nil {
		logger.Error(
			"failed to connect redis",
			"error", err,
		)

		os.Exit(1)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error(
				"failed to close redis",
				"error", err,
			)
		}
	}()

	fiberApp := app.New(logger, redisClient)

	go func() {
		if err := fiberApp.Listen(cfg.ListenAddress()); err != nil {
			logger.Error(
				"failed to start server",
				"error", err,
			)

			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	logger.Info("shutting down server")

	if err := fiberApp.Shutdown(); err != nil {
		logger.Error(
			"failed to shutdown server",
			"error", err,
		)

		os.Exit(1)
	}

	logger.Info("server stopped")
}
