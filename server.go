package main

import (
	"fmt"
	"log"
	"myserver/db"
	"myserver/handlers"
	"net/http"
)

func main() {
	database := db.Connect()
	defer database.Close()

	http.HandleFunc("/put", handlers.Put)
	http.HandleFunc("/get", handlers.Get)
	http.HandleFunc("/delete", handlers.Delete)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
