package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"go1f/pkg/api"
	"go1f/pkg/db"
	"go1f/pkg/server"
	"os"
)

func main() {

	_ = godotenv.Load()
	TODO_PASSWORD := os.Getenv("TODO_PASSWORD")
	api.PASS = TODO_PASSWORD

	TODO_DBFILE := os.Getenv("TODO_DBFILE")
	if TODO_DBFILE == "" {
		TODO_DBFILE = "scheduler.db"
	}

	DB, err := sql.Open("sqlite", TODO_DBFILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.DB = DB
	defer DB.Close()

	err = db.Init(TODO_DBFILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Запуск WEB сервера
	TODO_PORT := os.Getenv("TODO_PORT")
	if TODO_PORT == "" {
		TODO_PORT = "7540"
	}
	str := startServer(TODO_PORT)
	if str != "" {
		fmt.Println(str)
	}

}

func startServer(port string) string {
	err := server.StartServer(port)
	if err != nil {
		return fmt.Sprintf("Ошибка при запуске сервера: %s", err.Error())
	}
	return ""
}
