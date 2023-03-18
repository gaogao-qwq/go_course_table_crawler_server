package main

import (
	"course_table_server/internal/app/server/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", handlers.LoginHandler)
	log.Fatal(http.ListenAndServe(":56789", mux))
}
