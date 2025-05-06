package dto

import "time"

// SrvCreateCookie contains userID
type SrvCreateCookie struct {
	UserID string
}

// SrvDeleteCookie contains Token
type SrvDeleteCookie struct {
	Token string
}

// SrvCheckCookie contains Token
type SrvCheckCookie struct {
	Token string
}

// SrvSuccessResp contains check flag
type SrvSuccessResp struct {
	Success bool
}

// Cookie contains cookie data
type Cookie struct {
	Name    string
	UserID  string
	TokenID string
	Expiry  time.Time
}

// Token contains cookie token data
type Token struct {
	UserID  string
	TokenID string
	Expiry  time.Time
}
