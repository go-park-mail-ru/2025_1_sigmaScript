package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/common"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/messages"
	mocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery/mocks" // Import the generated mocks
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/delivery/http/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCSRFHandler_CreateCSRF_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dummyCookie := &config.Cookie{
		SessionName: "session_id",
		HTTPOnly:    true,
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
		Path:        "/",
	}
	ctx := config.WrapCookieContext(context.Background(), dummyCookie)

	mockSessionSvc := mocks.NewMockSessionServiceInterface(ctrl)

	updateReq := dto.UpdateUserRequest{
		Username:            "newusername",
		OldPassword:         "oldpassword",
		NewPassword:         "newpassword",
		RepeatedNewPassword: "newpassword",
		Avatar:              "newavatar.png",
	}
	jsonReq, err := json.Marshal(updateReq)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/auth/csrf-token", bytes.NewReader(jsonReq))
	req = req.WithContext(ctx)

	req.AddCookie(&http.Cookie{Name: "session_id", Value: "oldsession"})
	rec := httptest.NewRecorder()

	mockSessionSvc.
		EXPECT().
		GetSession(gomock.Any(), "oldsession").
		Return("oldusername", nil).
		Times(1)

	handler := NewCSRFHandler(ctx, mockSessionSvc)
	handler.CreateCSRFTokenHandler(rec, req)

	res := rec.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp ds.Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, messages.SuccessfulNewCSRFToken, resp.Message)

	found := false
	headerValue := res.Header.Get(HeaderCSRFToken)
	if len(headerValue) == common.CSRF_TOKEN_LENGTH {
		found = true
	}
	assert.True(t, found, "expected X-CSRF-Token Header value to be set as 64 byte sha256 string, got '%v' instead", headerValue)
}
