package main

import (
	"Rental_Sys/Auth"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/register_user", Auth.Register_user).Methods("POST")

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
