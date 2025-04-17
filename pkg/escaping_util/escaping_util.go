package escapingutil

import (
	"errors"
	"html"
)

const (
	DEFAULT_TEXT_MAX_LENGTH = 500
)

var (
	ErrorMaxLength = errors.New("text exeeds the allowed maximum length")
)

func ValidateInputTextData(textData string, textMaxLength ...int) (string, error) {
	maxLength := 1
	if len(textMaxLength) > 0 && textMaxLength[0] > 0 {
		maxLength = textMaxLength[0]
	}

	if len(textData) > maxLength {
		return "", ErrorMaxLength
	}

	return html.EscapeString(textData), nil
}
