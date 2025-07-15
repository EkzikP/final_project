package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"go1f/pkg/api"
	"go1f/pkg/db"
	"go1f/pkg/server"
	"log"
	"os"
)

func main() {

	_ = godotenv.Load()
	todoPassword := os.Getenv("TODO_PASSWORD")
	api.PASS = todoPassword

	todoDbfile := os.Getenv("TODO_DBFILE")
	if todoDbfile == "" {
		todoDbfile = "scheduler.db"
	}

	DB, err := sql.Open("sqlite", todoDbfile)
	if err != nil {
		log.Fatal(err)
		return
	}
	db.DB = DB
	defer DB.Close()

	err = db.Init(todoDbfile)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Запуск WEB сервера
	todoPort := os.Getenv("TODO_PORT")
	if todoPort == "" {
		todoPort = "7540"
	}
	err = startServer(todoPort)
	if err != nil {
		log.Fatal(err)
	}

}

func startServer(port string) error {
	err := server.StartServer(port)
	if err != nil {
		return fmt.Errorf("Ошибка при запуске сервера: %w", err)
	}
	return nil
}
