package main

import (
	"database/sql"
	"fmt"
	"os"

	sessiondelivery "github.com/AlexNov03/AuthService/delivery/session"
	sessionrepo "github.com/AlexNov03/AuthService/repository/session"
	userrepo "github.com/AlexNov03/AuthService/repository/user"
	"github.com/AlexNov03/AuthService/server"
	sessionuc "github.com/AlexNov03/AuthService/usecase/session"
	useruc "github.com/AlexNov03/AuthService/usecase/user"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error while loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=db1 sslmode=%s", dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("error while connecting to database")
	}
	defer db.Close()

	sessionRepo := sessionrepo.NewSessionRepo()

	userRepo := userrepo.NewUserRepo(db)

	sessionUC := sessionuc.NewSessionUsecase(sessionRepo)
	userUC := useruc.NewUserUsecase(userRepo)

	sessionDeliv := sessiondelivery.NewSessionDelivery(sessionUC, userUC)

	handler := server.NewHandler(sessionDeliv)

	server := server.NewServer(handler)

	server.ListenAndServe()
}
