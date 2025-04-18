package escapingutil

import (
	"errors"
	"html"
	"strings"
)

const (
	DEFAULT_TEXT_MAX_LENGTH = int64(500)
)

var (
	ErrorMaxLength            = errors.New("text exeeds the allowed maximum length")
	ErrorEmptyOrInvalidString = errors.New("text string is null or contains only invalid symbols")
)

func ValidateInputTextData(textData string, textMaxLength ...int64) (string, error) {
	maxLength := DEFAULT_TEXT_MAX_LENGTH
	if len(textMaxLength) > 0 && textMaxLength[0] > 0 {
		maxLength = textMaxLength[0]
	}

	if int64(len(textData)) > maxLength {
		return "", ErrorMaxLength
	}
	trimmedDataString := html.EscapeString(strings.TrimSpace(textData))

	if len(trimmedDataString) < 1 {
		return "", ErrorEmptyOrInvalidString
	}

	return trimmedDataString, nil
}
