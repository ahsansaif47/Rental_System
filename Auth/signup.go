package Auth

import (
	utils "Rental_System/Utils"
	"encoding/json"
	"log"
	"net/http"
)

func Register_user(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	type registration_response struct {
		RegResponse string
		Status      int
	}

	// connStr, err := utils.Connect_postgres()
	if utils.ConnErr != nil {
		log.Fatal("PostgreSQL Connection Error: ", utils.ConnErr.Error())
	} else {
		insert_user_query, err := utils.ConnStr.Prepare(`Insert into users(email, password) Values($1, $2)`)
		if err != nil {
			log.Println("Error preparing query: ", err.Error())
		}
		encryptedPass, err := utils.Encrypt(password)

		if err != nil {
			log.Println("Error encrypting password: ", err.Error())
		} else {
			if hashStatus := utils.Compare_Encryption(password, encryptedPass); hashStatus {
				_, err = insert_user_query.Exec(email, encryptedPass)
				if err != nil {
					if utils.Unique_constraint_violation_check(err) {
						reg_resp := registration_response{
							RegResponse: "User Already Resigtered",
							Status:      200,
						}
						json_data, _ := json.Marshal(reg_resp)
						_, err = w.Write(json_data)
						if err != nil {
							log.Println("Error Writing Response: ", err.Error())
						}
						return
					} else {
						log.Fatal("Error adding new user: ", err.Error())
						return
					}
				} else {
					reg_resp := registration_response{
						RegResponse: "User Resigtered",
						Status:      200,
					}
					json_data, _ := json.Marshal(reg_resp)
					_, err := w.Write(json_data)
					if err != nil {
						log.Println("Error Writing Response: ", err.Error())
					}
					return
				}
			}
		}
	}
}
