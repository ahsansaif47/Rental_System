package Auth

import (
	utils "Rental_Sys/Utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func check_user_registration(email string) bool {
	query := fmt.Sprintf("Select email from users where email = %s", email)
	fmt.Print("Final query is: ", query)
	emails, err := utils.Execute_query(query)

	if err != nil {
		if emails != nil {
			count := 0
			if emails.Next() {
				count++
			}
			if count >= 1 {
				return true
			}
		}
		return false
	}
	return false
}

func Register_user(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	type registration_response struct {
		response string
		status   int
	}

	if check_user_registration(email) {
		reg_resp := registration_response{
			response: "Already registered",
			status:   200,
		}

		json_data, _ := json.Marshal(reg_resp)
		_, err := w.Write(json_data)
		if err != nil {
			log.Println(err.Error())
		}
		return
	} else {
		connStr, err := utils.Connect_postgres()
		if err != nil {
			log.Fatal("Connection Error: ", err.Error())
		} else {
			insert_user_query, err := connStr.Prepare(`Insert into users(email, password) Values($1, $2)`)
			if err != nil {
				log.Println("Error preparing query: ", err.Error())
			}
			_, err = insert_user_query.Exec(email, password)
			if err != nil {
				log.Fatal("Error adding new user: ", err.Error())
			}
		}
	}
}
