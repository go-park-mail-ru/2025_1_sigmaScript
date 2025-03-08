package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GET(t *testing.T) {
	rr, req := getResponseRequest(t, "GET", "/collections/", nil)
	GetCollections(rr, req)
	assertHeaders(t, http.StatusOK, rr)
	var got mocks.Collections
	expected := mocks.MainPageCollections
	err := json.NewDecoder(rr.Body).Decode(&got)
	assert.NoError(t, err, errs.ErrParseJSON)
	assert.True(t, reflect.DeepEqual(got, expected))
}
