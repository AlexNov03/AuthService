package sessionuc

type SessionRepo interface {
	AddSession(sessionID string, userID string)
	GetSessionUser(sessionID string) (string, error)
	DeleteSession(sessionID string)
}

type SessionUsecase struct {
	repo SessionRepo
}

func NewSessionUsecase(sessionRepo SessionRepo) *SessionUsecase {
	return &SessionUsecase{repo: sessionRepo}
}

func (uc *SessionUsecase) AddSession(sessionID string, userID string) {
	uc.repo.AddSession(sessionID, userID)
}

func (uc *SessionUsecase) GetSessionUser(sessionID string) (string, error) {
	res, err := uc.repo.GetSessionUser(sessionID)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (uc *SessionUsecase) DeleteSession(sessionID string) {
	uc.repo.DeleteSession(sessionID)
}
