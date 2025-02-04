package sessiondelivery

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sessionrepo "github.com/AlexNov03/AuthService/repository/session"
	userrepo "github.com/AlexNov03/AuthService/repository/user"
	sessionuc "github.com/AlexNov03/AuthService/usecase/session"
	useruc "github.com/AlexNov03/AuthService/usecase/user"
)

func TestApi(t *testing.T) {
	testCases := []struct {
		Name         string
		Code         int
		Type         string
		Method       string
		RequestBody  string
		ResponseBody string
	}{
		{
			Name:         "test1",
			Type:         "login",
			Method:       http.MethodPost,
			RequestBody:  `{"email":"alexnov03@mail.ru", "password":"12345"}`,
			Code:         404,
			ResponseBody: "{\"error\":\"repository: No user info found by this email and password\"}\n",
		},
		{
			Name:         "test2",
			Type:         "signup",
			Method:       http.MethodPost,
			RequestBody:  `{"firstname":"alex", "secondname":"novak", "email":"alexnov03@mail.ru", "password":"12345"}`,
			Code:         200,
			ResponseBody: "{}\n",
		},
		{
			Name:         "test3",
			Type:         "login",
			Method:       http.MethodPost,
			RequestBody:  `{"email":"alexnov03@mail.ru", "password":"12345"}`,
			Code:         200,
			ResponseBody: "{}\n",
		},
		{
			Name:         "test4",
			Type:         "login",
			Method:       http.MethodPost,
			RequestBody:  `{"email":"alexnov3@mail.ru", "password":"12345"}`,
			Code:         404,
			ResponseBody: "{\"error\":\"repository: No user info found by this email and password\"}\n",
		},
	}

	sessionRepo := sessionrepo.NewSessionRepo()
	userRepo := userrepo.NewUserRepo()

	sessionUC := sessionuc.NewSessionUsecase(sessionRepo)
	userUC := useruc.NewUserUsecase(userRepo)

	sessionDelivery := NewSessionDelivery(sessionUC, userUC)

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			body := strings.NewReader(tt.RequestBody)
			request := httptest.NewRequest(tt.Method, "http://localhost/"+tt.Type, body)
			recorder := httptest.NewRecorder()
			switch tt.Type {
			case "login":
				sessionDelivery.Login(recorder, request)
			case "logout":
				sessionDelivery.LogOut(recorder, request)
			case "signup":
				sessionDelivery.SignUp(recorder, request)
			}

			responseBody, err := io.ReadAll(recorder.Body)
			if err != nil {
				t.Errorf("%v", err)
			}

			if recorder.Code != tt.Code {
				t.Errorf("statusCode is not valid, wanted: %v, having: %v", tt.Code, recorder.Code)
			}

			if string(responseBody) != tt.ResponseBody {
				t.Errorf("responseBody is not valid, wanted: %v, having: %v ", tt.ResponseBody, string(responseBody))
			}

		})
	}
}
