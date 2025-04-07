package repository

const (
	noData = ""
)

// type AuthRepository struct {
// 	// username --> hashedPass
// 	users synccredmap.SyncCredentialsMap
// 	// sessionID --> username
// 	sessions synccredmap.SyncCredentialsMap
// }

// func NewAuthRepository() *AuthRepository {
// 	return &AuthRepository{
// 		users:    *synccredmap.NewSyncCredentialsMap(),
// 		sessions: *synccredmap.NewSyncCredentialsMap(),
// 	}
// }

// func (r *AuthRepository) GetSession(ctx context.Context, sessionID string) (string, error) {
// 	logger := log.Ctx(ctx)

// 	username, ok := r.sessions.Load(sessionID)
// 	if !ok {
// 		logger.Error().Err(errors.Wrap(errs.ErrSessionNotExists, errs.ErrMsgFailedToGetSession)).Msg(errs.ErrMsgSessionNotExists)
// 		return noData, errs.ErrSessionNotExists
// 	}
// 	return username, nil
// }

// func (r *AuthRepository) DeleteSession(ctx context.Context, sessionID string) error {
// 	logger := log.Ctx(ctx)

// 	_, ok := r.sessions.Load(sessionID)
// 	if !ok {
// 		logger.Error().Err(errors.Wrap(errs.ErrSessionNotExists, errs.ErrMsgFailedToGetSession)).Msg(errs.ErrMsgSessionNotExists)
// 		return errs.ErrSessionNotExists
// 	}

// 	r.sessions.Delete(sessionID)
// 	return nil
// }

// func (r *AuthRepository) StoreSession(ctx context.Context, newSessionID, login string) error {
// 	r.sessions.Store(newSessionID, login)
// 	return nil
// }

// func (r *AuthRepository) CreateUser(ctx context.Context, login, hashedPass string) error {
// 	logger := log.Ctx(ctx)

// 	if _, exists := r.users.Load(login); exists {
// 		logger.Error().Err(errors.New(errs.ErrAlreadyExists)).Msg(common.MsgUserWithNameAlreadyExists)
// 		return errors.New(errs.ErrAlreadyExists)
// 	}
// 	r.users.Store(login, string(hashedPass))
// 	return nil
// }

// func (r *AuthRepository) GetUser(ctx context.Context, login string) (hashedPass string, errRepo error) {
// 	logger := log.Ctx(ctx)

// 	hashedPass, exists := r.users.Load(login)
// 	if exists {
// 		return hashedPass, nil
// 	}
// 	err := errors.New(errs.ErrIncorrectLogin)
// 	logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(err.Error())
// 	return noData, err
// }
