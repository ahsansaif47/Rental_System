package Settings

import (
	utils "Rental_System/Utils"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func deleteUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	isEmail := false
	if strings.Contains(email, "@") {
		isEmail = true
	}

	type Response struct {
		Resp   string
		Status int
	}

	if isEmail {
		userEmail, err := utils.ConnStr.Query(`SELECT email from users where email = $1`, email)
		if err != nil {
			log.Fatal("Error querying user: ", err)
		} else {
			defer userEmail.Close()

			if rowsErr := utils.Rows_iteration_error_check(userEmail); rowsErr != nil {
				log.Fatal("Error iterating rows: ", rowsErr)
			} else {
				var foundEmail string
				for userEmail.Next() {
					err := userEmail.Scan(&foundEmail)
					if err != nil {
						log.Println("Error setting value: ", err)
					}
				}

				if foundEmail == email {
					_, err := utils.ConnStr.Query(`DELETE FROM user WHERE email = $1`, foundEmail)
					if err != nil {
						log.Fatal("Error running delete user query: ", err)
					} else {
						resp := Response{Resp: "User deleted", Status: 200}
						json.NewEncoder(w).Encode(resp)
					}
				}
			}
		}
	}
}
