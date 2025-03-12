package errs

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
  ErrGenerateSession               = "Error generating session ID"
  ErrGenerateSessionShort          = "generate_session_error"
  ErrUnauthorized                  = "Unauthorized"
  ErrUnauthorizedShort             = "unauthorized"
  ErrSessionNotExists              = "Session does not exist"
  ErrSessionNotExistsShort         = "not_exists"
  ErrInvalidPassword               = "Invalid password"
  ErrInvalidPasswordShort          = "invalid_password"
  ErrSomethingWentWrong            = "something went wrong"
  ErrBadPayload                    = "bad payload"
  ErrInvalidEmail                  = "Invalid email"
  ErrInvalidEmailShort             = "invalid_email"
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
  ErrNegativeSessionIDLength = "Negative session ID length"
  ErrLengthTooShort          = "Length too short"
  ErrLengthTooLong           = "Length too long"
)
