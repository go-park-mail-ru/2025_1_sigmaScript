package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getResponseRequest(t *testing.T, method, target string, data any) (*httptest.ResponseRecorder, *http.Request) {
	var req *http.Request
	jsonData, err := json.Marshal(data)
	require.NoError(t, err, errs.ErrParseJSON, errs.ErrEncodeJSON)
	jsonReader := bytes.NewReader(jsonData)
	if method == "GET" {
		req = httptest.NewRequest(method, target, nil)
	} else {
		req = httptest.NewRequest(method, target, jsonReader)
	}
	rr := httptest.NewRecorder()
	return rr, req
}

func assertHeaders(t *testing.T, code int, rr *httptest.ResponseRecorder) {
	assert.Equal(t, code, rr.Code, errs.ErrWrongResponseCode)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), errs.ErrWrongHeaders)
}
