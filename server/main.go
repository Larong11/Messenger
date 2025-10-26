package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"server/application/use_cases/user"
	"server/infrastructure/persistence"
	infr_http "server/interface/http"
	"server/interface/http/handlers"
)

var ctx = context.Background()

func main() {

	pool, err := pgxpool.New(ctx, "postgres://postgres:135790@localhost:5432/messenger_db?sslmode=disable")
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database!")

	userRepo := persistence.NewPostgresUserRepository(pool)
	registerUserUseCases := user.NewRegisterUserUseCases(userRepo)
	userHandler := handlers.NewUserHandler(registerUserUseCases)
	router := infr_http.NewRouter(userHandler)
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
