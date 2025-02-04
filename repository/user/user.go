package userrepo

import (
	"net/http"
	"sync"

	interr "github.com/AlexNov03/AuthService/errors/internalerr"
	"github.com/AlexNov03/AuthService/models"
)

// *****************************************************
// *												   *
// *     Mock БД на структурах для тестирования        *
// *  Так делать нельзя и скоро здесь будет Postgres   *
// *                                                   *
// *****************************************************

type UserRepo struct {
	Storage map[string]*models.UserRegData
	mu      sync.RWMutex
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		Storage: make(map[string]*models.UserRegData),
		mu:      sync.RWMutex{},
	}
}

func (ur *UserRepo) GetUserInfo(userID string) (*models.UserInfo, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()
	if val, ok := ur.Storage[userID]; ok {
		return &models.UserInfo{
			FirstName: val.FirstName,
			LastName:  val.LastName,
			Email:     val.Email,
		}, nil
	}
	return nil, interr.NewInternalError(http.StatusNotFound, "repository: No user info found by this userID")
}

func (ur *UserRepo) AddUserInfo(userID string, userInfo *models.UserRegData) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	ur.Storage[userID] = userInfo
}

func (ur *UserRepo) FindUser(loginData *models.UserLoginData) (*models.UserInfo, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()
	for key, value := range ur.Storage {
		if value.Email == loginData.Email && value.Password == loginData.Password {
			return &models.UserInfo{
				ID:        key,
				FirstName: value.FirstName,
				LastName:  value.LastName,
				Email:     value.Email}, nil
		}
	}
	return nil, interr.NewInternalError(http.StatusNotFound, "repository: No user info found by this email and password")
}
