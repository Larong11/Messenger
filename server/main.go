package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"server/application/use_cases/user"
	infr_http "server/infrastructure/http"
	"server/infrastructure/http/handlers"
	"server/infrastructure/persistence"
)

var ctx = context.Background()

func main() {

	pool, err := pgxpool.New(ctx, "postgres://<username>:<password>@localhost:5432/gotodo")
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
