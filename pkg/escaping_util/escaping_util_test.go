package escapingutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInputTextData_ValidData(t *testing.T) {

	tests := []struct {
		inputText   string
		correctText string
	}{
		{"ab", "ab"},
		{"abcdefghijklmnopqr", "abcdefghijklmnopqr"},
		{"user123", "user123"},
		{"user-name", "user-name"},
		{"user_name", "user_name"},
		{"UserName", "UserName"},
		{"User_123-Name", "User_123-Name"},
		{"123456", "123456"},
		{"--__--", "--__--"},
		{" ab ", "ab"},
		{"\t  abcdefghijklmnopqr\n", "abcdefghijklmnopqr"},
		{"<script>", "&lt;script&gt;"},
	}

	for _, dummyData := range tests {

		validatedData, err := ValidateInputTextData(dummyData.inputText)
		assert.NoError(t, err)

		assert.Equal(t, dummyData.correctText, validatedData)
	}
}

func TestValidateInputTextData_NotValidData(t *testing.T) {

	tests := []struct {
		inputText   string
		correctText string
	}{
		{"abcdefghijklmnopqrabcdefghijklmnopqrabcdefghijklmnopqrabcdefghijklmnopqr",
			"abcdefghijklmnopqrabcdefghijklmnopqrabcdefghijklmnopqrabcdefghijklmnopqr"},
		{"                    ",
			"                    "},
	}

	for _, dummyData := range tests {

		validatedData, err := ValidateInputTextData(dummyData.inputText, int64(10))
		assert.Error(t, err)

		assert.NotEqual(t, dummyData.correctText, validatedData)
	}
}
