package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	deliveryCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/delivery"
	repoCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/repository"
	serviceCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/service"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GET(t *testing.T) {
	rr, req := getResponseRequest(t, "GET", "/collections/", nil)

	collectionRepo := repoCollection.NewCollectionRepository(&mocks.MainPageCollections)
	collectionService := serviceCollection.NewCollectionService(collectionRepo)
	collectionHandler := deliveryCollection.NewCollectionHandler(collectionService)

	collectionHandler.GetMainPageCollections(rr, req)
	assertHeaders(t, http.StatusOK, rr)
	var got mocks.Collections
	expected := mocks.MainPageCollections
	err := json.NewDecoder(rr.Body).Decode(&got)
	assert.NoError(t, err, errs.ErrParseJSON)
	assert.True(t, reflect.DeepEqual(got, expected))
}
