package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Film struct {
	Position string `json:"position"`
	Title    string `json:"title"`
}

type TestCase struct {
	Name string `json:"name"`
	Data []Film `json:"data"`
}

func TestSendOK(t *testing.T) {
	tests := []TestCase{
		{
			Name: "simple 1",
			Data: []Film{
				{
					Position: "top 1",
					Title:    "Pulp Fiction",
				},
				{
					Position: "top 2",
					Title:    "Spider-Man",
				},
			},
		},
		{
			Name: "simple 2",
			Data: []Film{
				{
					Position: "top 1",
					Title:    "Star Wars",
				},
				{
					Position: "top 2",
					Title:    "Inside Out",
				},
			},
		},
	}
	rr := httptest.NewRecorder()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := SendJSON(t.Context(), rr, test.Data)
			require.NoError(t, err)

			var got []Film
			expected := test.Data
			err = json.NewDecoder(rr.Body).Decode(&got)
			require.NoError(t, err, errs.ErrParseJSON)
			assert.ElementsMatch(t, got, expected)
		})
	}
}

func TestSendFail(t *testing.T) {
	ch := make(chan int)
	rr := httptest.NewRecorder()
	t.Run("json encode error", func(t *testing.T) {
		err := SendJSON(t.Context(), rr, ch)
		fmt.Println(err.Error())
		require.Error(t, err)
	})
}

func TestReadOK(t *testing.T) {
	tests := []TestCase{
		{
			Name: "simple 1",
			Data: []Film{
				{
					Position: "top 1",
					Title:    "Pulp Fiction",
				},
				{
					Position: "top 2",
					Title:    "Spider-Man",
				},
			},
		},
		{
			Name: "simple 2",
			Data: []Film{
				{
					Position: "top 1",
					Title:    "Star Wars",
				},
				{
					Position: "top 2",
					Title:    "Inside Out",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			jsonData, err := json.Marshal(test.Data)
			require.NoError(t, err)
			jsonReader := bytes.NewReader(jsonData)
			req := httptest.NewRequest("POST", "/", jsonReader)

			var got []Film
			err = ReadJSON(req, &got)
			expected := test.Data
			require.NoError(t, err)
			assert.ElementsMatch(t, got, expected)
		})
	}
}

func TestReadFail(t *testing.T) {
	test := struct {
		Number int
	}{
		Number: 2034,
	}
	t.Run("json decode error", func(t *testing.T) {
		jsonData, err := json.Marshal(test)
		require.NoError(t, err)
		jsonReader := bytes.NewReader(jsonData)
		req := httptest.NewRequest("POST", "/", jsonReader)

		var got []Film
		err = ReadJSON(req, &got)
		require.Error(t, err)
	})
}

func TestSendError(t *testing.T) {
	errMsg := "error happened"
	rr := httptest.NewRecorder()
	t.Run("json encode error", func(t *testing.T) {
		SendError(t.Context(), rr, 200, errMsg, errMsg)

		var result ErrorResponse
		err := json.NewDecoder(rr.Result().Body).Decode(&result)

		assert.NoError(t, err)
		assert.Equal(t, errMsg, result.Error)
		assert.Equal(t, errMsg, result.Message)
	})
}
