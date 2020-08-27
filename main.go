package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/emanuelcondo/api-todo/db"
	routers "github.com/emanuelcondo/api-todo/routers"
)

func main() {
	//Load environmenatal variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDB()

	router := routers.Routes()
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, router)) // note, the port is usually gotten from the environment
}
