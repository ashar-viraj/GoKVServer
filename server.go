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

	http.HandleFunc("/create", handlers.Create)
	http.HandleFunc("/read", handlers.Read)
	http.HandleFunc("/update", handlers.Update)
	http.HandleFunc("/delete", handlers.Delete)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
