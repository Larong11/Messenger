package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/application/use_cases/user"
	"server/infrastructure/email/smtp"
	"server/infrastructure/persistence"
	"server/interface/http/handlers"
	infrhttp "server/interface/router"
	"server/interface/websocket"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	dbURL := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database!")

	smtpConfig := smtp.Config{
		Host:     "smtp.gmail.com", // or your SMTP host
		Port:     "587",            // 587 for TLS, 465 for SSL
		Username: "your-email@gmail.com",
		Password: "your-app-password", // Use app password for Gmail
		From:     "your-email@gmail.com",
	}
	emailService := smtp.NewEmailService(smtpConfig)

	userRepo := persistence.NewPostgresUserRepository(pool)
	registerUserUseCases := user.NewRegisterUserUseCases(userRepo, emailService)
	userHandler := handlers.NewUserHandler(registerUserUseCases)
	WSHandler := websocket.NewWSHandler()
	router := infrhttp.NewRouter(userHandler, WSHandler)
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
