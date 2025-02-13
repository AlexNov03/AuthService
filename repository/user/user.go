package userrepo

import (
	"database/sql"
	"errors"
	"net/http"

	interr "github.com/AlexNov03/AuthService/errors/internalerr"
	"github.com/AlexNov03/AuthService/models"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (ur *UserRepo) AddUser(userID string, userInfo *models.UserRegData) error {

	_, err := ur.DB.Exec(`INSERT INTO "user" (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)`,
		userInfo.FirstName, userInfo.LastName, userInfo.Email, userInfo.Password)
	if err != nil {
		return interr.NewInternalError(http.StatusInternalServerError, err.Error())
	}
	return nil

}

func (ur *UserRepo) GetUserByID(userID string) (*models.UserInfo, error) {

	userInfo := &models.UserInfo{}
	row := ur.DB.QueryRow(`SELECT user_id, first_name, last_name, email FROM "user" WHERE user_id=$1`, userID)

	err := row.Scan(&userInfo.ID, &userInfo.FirstName, &userInfo.LastName, &userInfo.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, interr.NewInternalError(http.StatusNotFound, "no user info found by this userID")
	}
	if err != nil {
		return nil, interr.NewInternalError(http.StatusInternalServerError, err.Error())
	}

	return userInfo, nil

}

func (ur *UserRepo) GetUserByEmail(email string) (*models.UserInfo, error) {
	userInfo := &models.UserInfo{}
	row := ur.DB.QueryRow(`SELECT user_id, first_name, last_name, email FROM "user" WHERE email=$1`, email)
	err := row.Scan(&userInfo.ID, &userInfo.FirstName, &userInfo.LastName, &userInfo.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, interr.NewInternalError(http.StatusNotFound, "no user info found by this email")
	}
	if err != nil {
		return nil, interr.NewInternalError(http.StatusInternalServerError, err.Error())
	}

	return userInfo, nil
}

func (ur *UserRepo) GetUserByLoginData(loginData *models.UserLoginData) (*models.UserInfo, error) {

	userInfo := &models.UserInfo{}
	row := ur.DB.QueryRow(`SELECT user_id, first_name, last_name, email FROM "user" WHERE email=$1 AND password=$2`,
		loginData.Email, loginData.Password)
	err := row.Scan(&userInfo.ID, &userInfo.FirstName, &userInfo.LastName, &userInfo.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, interr.NewInternalError(http.StatusNotFound, "no user info found by this email and password")
	}
	if err != nil {
		return nil, interr.NewInternalError(http.StatusInternalServerError, err.Error())
	}

	return userInfo, nil
}
