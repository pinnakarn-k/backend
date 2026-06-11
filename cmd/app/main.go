package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/internal/app"
	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/pdf"
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

	// test pdf
	rows := make([][]any, 0)

	for i := 1; i <= 120; i++ {
		rows = append(rows, []any{
			i,
			"สมชาย ใจดี",
			"123-6",
			"AOT",
			100,
			6250.00,
		})
	}

	data, err := pdf.GenerateWithConfig(
		pdf.ExportFile{
			Title:       "Historical Data",
			Subtitle:    "รายงานข้อมูลย้อนหลัง",
			GeneratedBy: "system",
			Headers: []string{
				"ลำดับ",
				"ชื่อลูกค้า",
				"เลขที่บัญชี",
				"หลักทรัพย์",
				"จำนวน",
				"มูลค่า",
			},
			Rows: rows,
		},
		pdf.Config{
			FontPath:        "fonts/THSarabunNew.ttf",
			FontName:        "THSarabun",
			HeaderImagePath: "assets/header.png",
			HeaderImageType: "PNG",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("report.pdf", data, 0644); err != nil {
		log.Fatal(err)
	}

	log.Println("report.pdf generated")
	// test pdf

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
