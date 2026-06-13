package pdf

import (
	"bytes"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func encryptPDF(data []byte, openPassword string) ([]byte, error) {
	if openPassword == "" {
		return data, nil
	}

	reader := bytes.NewReader(data)

	var encrypted bytes.Buffer

	ownerPassword := openPassword + "_owner"

	conf := model.NewAESConfiguration(
		openPassword,
		ownerPassword,
		256,
	)

	conf.Permissions = model.PermissionsAll

	if err := api.Encrypt(reader, &encrypted, conf); err != nil {
		return nil, fmt.Errorf("encrypt pdf: %w", err)
	}

	return encrypted.Bytes(), nil
}
