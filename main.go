package main

import (
	"course_table_server/internal/app/server/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/semester-list", handlers.SemesterListHandler)
	mux.HandleFunc("/course-table", handlers.CourseTableHandler)
	fmt.Println("Listen and serve on port :56789")
	log.Fatal(http.ListenAndServe(":56789", mux))
}
