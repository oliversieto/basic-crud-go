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
	router.HandleFunc("/users", user.CreateUser).Methods("POST")
	fmt.Println("Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
