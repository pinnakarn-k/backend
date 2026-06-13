package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func encryptExcel(data []byte, openPassword string) ([]byte, error) {
	if openPassword == "" {
		return data, nil
	}

	encrypted, err := excelize.Encrypt(data, &excelize.Options{
		Password: openPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("encrypt excel: %w", err)
	}

	return encrypted, nil
}
