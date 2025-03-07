package errs

// main
const (
  ErrLoadConfig  = "Error loading config"
  ErrStartServer = "Error starting server"
  ErrShutdown    = "Error shutting down"
)

// config
const (
  ErrInitializeConfig = "Error initializing config"
  ErrUnmarshalConfig  = "Error unmarshalling config"
  ErrReadConfig       = "Error reading config"
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
  ErrBcryptShort                   = "bcrypt_error"
  ErrSendJSON                      = "Error sending JSON"
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
  ErrEasyPassword     = "Easy password"
)
