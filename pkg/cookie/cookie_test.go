package cookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	mocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/service/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/cookie"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExpireOldSessionCookie_NoCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionSvc := mocks.NewMockSessionRepositoryInterface(ctrl)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	dummyCookie := &config.Cookie{
		SessionName: "session_id",
		HTTPOnly:    true,
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
		Path:        "/",
	}

	err := cookie.ExpireOldSessionCookie(rec, req, dummyCookie, mockSessionSvc)
	assert.NoError(t, err)

	respCookies := rec.Result().Cookies()
	found := false
	for _, c := range respCookies {
		if c.Name == dummyCookie.SessionName {
			found = true
			assert.True(t, c.Expires.Before(time.Now()))
		}
	}

	assert.False(t, found)
}

func TestExpireOldSessionCookie_WithCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionSvc := mocks.NewMockSessionRepositoryInterface(ctrl)

	expectedSessionID := "old_session_value"
	mockSessionSvc.
		EXPECT().
		DeleteSession(gomock.Any(), expectedSessionID).
		Return(nil).
		Times(1)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: expectedSessionID,
	})

	cfg := &config.Cookie{
		SessionName: "session_id",
		HTTPOnly:    true,
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
		Path:        "/",
	}

	err := cookie.ExpireOldSessionCookie(rec, req, cfg, mockSessionSvc)
	assert.NoError(t, err)

	respCookies := rec.Result().Cookies()
	var expiredCookie *http.Cookie
	for _, c := range respCookies {
		if c.Name == cfg.SessionName {
			expiredCookie = c
			break
		}
	}
	assert.NotNil(t, expiredCookie)
	assert.True(t, expiredCookie.Expires.Before(time.Now()))
}

func TestPreparedNewCookie(t *testing.T) {
	cfg := &config.Cookie{
		SessionName: "session_id",
		HTTPOnly:    true,
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
		Path:        "/",
	}
	newSessionID := "new_session"

	c := cookie.PreparedNewCookie(cfg, newSessionID)

	assert.Equal(t, cfg.SessionName, c.Name)
	assert.Equal(t, newSessionID, c.Value)
	assert.Equal(t, cfg.HTTPOnly, c.HttpOnly)
	assert.Equal(t, cfg.Secure, c.Secure)
	assert.Equal(t, cfg.SameSite, c.SameSite)
	assert.Equal(t, cfg.Path, c.Path)

	expectedExpire := time.Now().AddDate(0, 0, cookie.CookieDaysLimit)
	assert.InDelta(t, expectedExpire.Unix(), c.Expires.Unix(), 5)
}

func TestPreparedExpiredCookie(t *testing.T) {
	cfg := &config.Cookie{
		SessionName: "session_id",
		HTTPOnly:    true,
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
		Path:        "/",
	}
	c := cookie.PreparedExpiredCookie(cfg)

	assert.Equal(t, cfg.SessionName, c.Name)
	assert.Equal(t, "", c.Value)
	assert.Equal(t, cfg.HTTPOnly, c.HttpOnly)
	assert.Equal(t, cfg.Secure, c.Secure)
	assert.Equal(t, cfg.SameSite, c.SameSite)
	assert.Equal(t, cfg.Path, c.Path)

	assert.True(t, c.Expires.Before(time.Now()))
}
