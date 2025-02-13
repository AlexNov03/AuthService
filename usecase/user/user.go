package useruc

import (
	"errors"
	"net/http"

	"github.com/AlexNov03/AuthService/errors/internalerr"
	"github.com/AlexNov03/AuthService/models"
	"github.com/google/uuid"
)

type UserRepo interface {
	AddUser(userID string, userInfo *models.UserRegData) error
	GetUserByID(userID string) (*models.UserInfo, error)
	GetUserByEmail(email string) (*models.UserInfo, error)
	GetUserByLoginData(loginData *models.UserLoginData) (*models.UserInfo, error)
}

type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(userRepo UserRepo) *UserUsecase {
	return &UserUsecase{repo: userRepo}
}

func (uc *UserUsecase) SignUp(regData *models.UserRegData) (*models.UserInfo, error) {
	_, err := uc.repo.GetUserByEmail(regData.Email)

	if err == nil {
		return nil, internalerr.NewInternalError(http.StatusConflict, "this email is already used")
	}

	internalError := &internalerr.InternalError{}
	if ok := errors.As(err, &internalError); ok {
		if internalError.Code != http.StatusNotFound {
			return nil, internalError
		}
	}

	userID := uuid.NewString()
	err = uc.repo.AddUser(userID, regData)

	if err != nil {
		return nil, err
	}
	return &models.UserInfo{ID: userID, FirstName: regData.FirstName,
		LastName: regData.LastName, Email: regData.Email}, nil
}

func (uc *UserUsecase) Login(loginData *models.UserLoginData) (*models.UserInfo, error) {
	info, err := uc.repo.GetUserByLoginData(loginData)
	if err != nil {
		return nil, err
	}
	return info, nil
}
