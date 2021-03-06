package main

import (
	"basic-crud/user"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", user.Create).Methods(http.MethodPost)
	router.HandleFunc("/users", user.GetAll).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", user.GetOne).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", user.Update).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", user.Delete).Methods(http.MethodDelete)

	fmt.Println("Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
