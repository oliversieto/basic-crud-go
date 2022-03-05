package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	fmt.Println("Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
