package main

import (
	sessiondelivery "github.com/AlexNov03/AuthService/delivery/session"
	sessionrepo "github.com/AlexNov03/AuthService/repository/session"
	userrepo "github.com/AlexNov03/AuthService/repository/user"
	"github.com/AlexNov03/AuthService/server"
	sessionuc "github.com/AlexNov03/AuthService/usecase/session"
	useruc "github.com/AlexNov03/AuthService/usecase/user"
)

func main() {
	sessionRepo := sessionrepo.NewSessionRepo()
	userRepo := userrepo.NewUserRepo()

	sessionUC := sessionuc.NewSessionUsecase(sessionRepo)
	userUC := useruc.NewUserUsecase(userRepo)

	sessionDeliv := sessiondelivery.NewSessionDelivery(sessionUC, userUC)

	handler := server.NewHandler(sessionDeliv)

	server := server.NewServer(handler)

	server.ListenAndServe()
}
