package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/AlexNov03/AuthService/delivery/sessiondelivery"
	"github.com/AlexNov03/AuthService/delivery/taskdelivery"
	"github.com/AlexNov03/AuthService/middleware"
	sessionrepo "github.com/AlexNov03/AuthService/repository/session"
	taskrepo "github.com/AlexNov03/AuthService/repository/task"
	userrepo "github.com/AlexNov03/AuthService/repository/user"
	"github.com/AlexNov03/AuthService/server"
	sessionuc "github.com/AlexNov03/AuthService/usecase/session"
	taskuc "github.com/AlexNov03/AuthService/usecase/task"
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
	taskRepo := taskrepo.NewTaskRepo(db)

	sessionUC := sessionuc.NewSessionUsecase(sessionRepo)
	userUC := useruc.NewUserUsecase(userRepo)
	taskUC := taskuc.NewTaskUsecase(taskRepo)

	sessionDeliv := sessiondelivery.NewSessionDelivery(sessionUC, userUC)
	taskDeliv := taskdelivery.NewTaskDelivery(taskUC)

	middleware := middleware.NewMiddleware(sessionUC)

	handler := server.NewHandler(sessionDeliv, taskDeliv, middleware)

	server := server.NewServer(handler)

	server.ListenAndServe()
}
