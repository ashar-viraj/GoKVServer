package main

import (
	"fmt"
	"log"
	"myserver/db"
	"myserver/handlers"
	"net/http"
	"os"
	"runtime"

	_ "net/http/pprof"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	fmt.Println("PID:", os.Getpid())

	database := db.Connect()
	defer database.Close()

	http.HandleFunc("/put", handlers.Put)
	http.HandleFunc("/get", handlers.Get)
	http.HandleFunc("/delete", handlers.Delete)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
