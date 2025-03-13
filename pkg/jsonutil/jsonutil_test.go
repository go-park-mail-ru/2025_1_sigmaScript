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
	Header string `json:"Header"`
	Body   string `json:"Beader"`
}

type Test struct {
	Name string `json:"Neader"`
	Data []Film `json:"Data"`
}

func TestSendOK(t *testing.T) {
	tests := []Test{
		{
			Name: "simple 1",
			Data: []Film{
				{
					Header: "top 1",
					Body:   "Pulp Fiction",
				},
				{
					Header: "top 2",
					Body:   "Spider-Man",
				},
			},
		},
		{
			Name: "simple 2",
			Data: []Film{
				{
					Header: "top 1",
					Body:   "Star Wars",
				},
				{
					Header: "top 2",
					Body:   "Inside Out",
				},
			},
		},
	}
	rr := httptest.NewRecorder()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := SendJSON(rr, test.Data)
			require.NoError(t, err)

			var got []Film
			expected := test.Data
			err = json.NewDecoder(rr.Body).Decode(&got)
			require.NoError(t, err, errs.ErrParseJSON)
			assert.True(t, CompareFilmSlice(got, expected))
		})
	}
}

func TestSendFail(t *testing.T) {
	ch := make(chan int)
	rr := httptest.NewRecorder()
	t.Run("json encode error", func(t *testing.T) {
		err := SendJSON(rr, ch)
		fmt.Println(err.Error())
		require.Error(t, err)
	})
}

func TestReadOK(t *testing.T) {
	tests := []Test{
		{
			Name: "simple 1",
			Data: []Film{
				{
					Header: "top 1",
					Body:   "Pulp Fiction",
				},
				{
					Header: "top 2",
					Body:   "Spider-Man",
				},
			},
		},
		{
			Name: "simple 2",
			Data: []Film{
				{
					Header: "top 1",
					Body:   "Star Wars",
				},
				{
					Header: "top 2",
					Body:   "Inside Out",
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
			assert.True(t, CompareFilmSlice(got, expected))
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

func CompareFilmSlice(first, second []Film) bool {
	if len(first) != len(second) {
		return false
	}
	for ind := range first {
		if first[ind].Header != second[ind].Header || first[ind].Body != second[ind].Body {
			return false
		}
	}
	return true
}
