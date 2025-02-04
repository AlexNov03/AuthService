package useruc

import (
	"github.com/AlexNov03/AuthService/models"
	"github.com/google/uuid"
)

type UserRepo interface {
	GetUserInfo(userID string) (*models.UserInfo, error)
	AddUserInfo(userID string, userInfo *models.UserRegData)
	FindUser(loginData *models.UserLoginData) (*models.UserInfo, error)
}

type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(userRepo UserRepo) *UserUsecase {
	return &UserUsecase{repo: userRepo}
}

func (uc *UserUsecase) GetUserInfo(userID string) (*models.UserInfo, error) {
	res, err := uc.repo.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *UserUsecase) AddUserInfo(userInfo *models.UserRegData) string {
	userID := uuid.NewString()
	uc.repo.AddUserInfo(userID, userInfo)
	return userID
}

func (uc *UserUsecase) FindUser(loginData *models.UserLoginData) (*models.UserInfo, error) {
	info, err := uc.repo.FindUser(loginData)
	if err != nil {
		return nil, err
	}
	return info, nil
}
