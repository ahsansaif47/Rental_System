package main

import (
	"Rental_System/Auth"
	settings "./Settings"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/register_user", Auth.Register_user).Methods("POST")
	router.HandleFunc("/login", Auth.Login).Methods("GET")
	router.HandleFunc("/delete_user", settings.deleteUser).Methods("GET")

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", router)

}
