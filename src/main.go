package main

import (
	"finances_manager_go/database"
	"finances_manager_go/handlers"
	"finances_manager_go/middlewares"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	database.Setup()
	api := &handlers.APIEnv{
		DB: database.GetDB(),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/add", middlewares.Auth(api.AddIncome))
	mux.HandleFunc("/signup", api.CreateUser)
	http.ListenAndServe(":"+port, mux)
}
