package emailtemplate

import (
	"os"
	"strings"
)

func RenderHTMLFile(
	path string,
	values map[string]string,
) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	html := string(data)

	for key, value := range values {
		placeholder := "{" + key + "}"

		html = strings.ReplaceAll(
			html,
			placeholder,
			value,
		)
	}

	return html, nil
}

// how to use
// html, err := emailtemplate.RenderHTMLFile(
// 	"templates/hd_email.html",
// 	map[string]string{
// 		"x": "คุณสมชาย",
// 		"y": "REF001",
// 		"customer_name": "คุณสมชาย",
// 	},
// )
// if err != nil {
// 	return err
// }
