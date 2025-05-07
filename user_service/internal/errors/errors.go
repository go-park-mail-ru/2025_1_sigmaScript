package errs

import "errors"

// main
const (
	ErrLoadConfig  = "Error loading config"
	ErrStartServer = "Error starting server"
	ErrShutdown    = "Error shutting down"
)

// config
const (
	ErrInitializeConfig  = "Error initializing config"
	ErrUnmarshalConfig   = "Error unmarshalling config"
	ErrReadConfig        = "Error reading config"
	ErrReadEnvironment   = "Error reading .env file"
	ErrGetDirectory      = "Error getting directory"
	ErrDirectoryNotFound = "Error finding directory"
)

// handlers
const (
	ErrParseJSON                     = "Error parsing JSON"
	ErrParseJSONShort                = "parse_json_error"
	ErrAlreadyExists                 = "Already exists"
	ErrAlreadyExistsShort            = "already_exists"
	ErrPasswordsMismatch             = "Passwords mismatch"
	ErrPasswordsMismatchShort        = "passwords_mismatch"
	ErrBcrypt                        = "Error hashing password"
	ErrSendJSON                      = "Error sending JSON"
	ErrIncorrectLogin                = "user with this login does not exist"
	ErrIncorrectPassword             = "provided password is incorrect"
	ErrIncorrectLoginOrPassword      = "Incorrect login or password"
	ErrIncorrectLoginOrPasswordShort = "not_found"
	ErrMsgGenerateSession            = "error generating session ID"
	ErrMsgGenerateSessionShort       = "generate_session_error"
	ErrUnauthorized                  = "Unauthorized"
	ErrUnauthorizedShort             = "unauthorized"
	ErrMsgSessionNotExists           = "session does not exist"
	ErrMsgSessionNotExistsShort      = "not_exists"
	ErrInvalidPassword               = "Invalid password"
	ErrInvalidPasswordShort          = "invalid_password"
	ErrSomethingWentWrong            = "something went wrong"
	ErrBadPayload                    = "bad payload"
	ErrInvalidLogin                  = "Invalid login"
	ErrInvalidLoginShort             = "invalid_login"
	ErrLengthLogin                   = "Login length must be 2-17 chars"
	ErrLengthLoginShort              = "length_login"
	ErrEmptyLogin                    = "Empty login"
	ErrEmptyLoginShort               = "empty_login"
	ErrNotFoundShort                 = "not_found"
	ErrMsgGenerateCSRFToken          = "error generating CSRF token"
	ErrMsgBadCSRFToken               = "bad csrf token"
	ErrParseForm                     = "error parsing request form"
	ErrParseFormShort                = "parse_form_error"
)

// jsonutil
const (
	ErrEncodeJSON      = "Error encoding JSON"
	ErrEncodeJSONShort = "encode_json_error"
	ErrCloseBody       = "Error closing body"
)

// validation/auth
const (
	ErrPasswordTooShort = "Password too short"
	ErrPasswordTooLong  = "Password too long"
	ErrEmptyPassword    = "Empty password"
)

// tests
const (
	ErrWrongHeaders      = "Wrong headers"
	ErrWrongResponseCode = "Wrong response code"
	ErrCookieEmpty       = "Cookie is empty"
	ErrCookieHttpOnly    = "Cookie HttpOnly flag is not set"
	ErrSessionCreated    = "Session should not have been created"
	ErrCookieExpire      = "Cookie must expire"
)

// session
const (
	ErrMsgNegativeSessionIDLength = "negative session ID length"
	ErrMsgLengthTooShort          = "length too short"
	ErrMsgLengthTooLong           = "length too long"
	ErrMsgFailedToGetSession      = "failed to get session"
)

// user
const (
	ErrMsgOnlyAllowedImageFormats = "only SVG, PNG, JPG, JPEG, and WebP are allowed"
)

// error types
var (
	ErrPersonNotFound = errors.New("person by this id not found")
	ErrMovieNotFound  = errors.New("movie by this id not found")

	ErrCollectionNotExist = errors.New("collection does not exist")

	ErrGenerateSession  = errors.New(ErrMsgGenerateSession)
	ErrSessionNotExists = errors.New(ErrMsgSessionNotExists)

	ErrInvalidFileType = errors.New("invalid_file_type")
)
