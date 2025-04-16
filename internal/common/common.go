package common

const (
	MsgUserWithNameAlreadyExists = "user with that name already exists"
	MsgUserDontHaveOldCookie     = "user dont have old cookie"
	MsgExpireOldCookieSuccess    = "successfully expired old sesssion cookie"

	COOKIE_DAYS_LIMIT        = 3
	COOKIE_EXPIRED_LAST_YEAR = -1
	REVIEWS_PER_PAGE         = 20
)

var ALLOWED_IMAGE_TYPES = map[string]bool{
	".svg":  true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
}

const (
	MB       = 1 << 20
	LIMIT_MB = 5
)

const (
	CSRF_TOKEN_NAME = "csrf_token"
)
