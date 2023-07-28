package Auth

import (
	utils "Rental_System/Utils"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func validate_email_server(dName string) bool {
	domains := []string{"gmail", "yahoo", "hormail"}
	for i := 0; i < len(domains); i++ {
		if domains[i] == dName {
			return true
		}
	}
	return false
}

func validate_email(email string) bool {
	email_regex := "[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-z]+.[a-z]{3,4}"
	regex, err := regexp.Compile(email_regex)
	if err != nil {
		log.Println("Error compiling regex: ", err.Error())
		return false
	}
	return regex.MatchString(email)
}

func Register_user(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	type registration_response struct {
		RegResponse string
		Status      int
	}

	is_email := validate_email(email)
	if !is_email {
		invalid_email_resp := registration_response{
			RegResponse: "Incorrect email format",
			Status:      200,
		}
		json_data, _ := json.Marshal(invalid_email_resp)
		_, err := w.Write(json_data)
		if err != nil {
			log.Println("Error writing response: ", err.Error())
		}
	} else {
		mailServer_dot_com := strings.Split(email, "@")
		domain := strings.Split(mailServer_dot_com[1], ".")[0]
		isSerevr_valid := validate_email_server(domain)
		if !isSerevr_valid {
			invalidServer_resp := registration_response{
				RegResponse: "Invalid email server",
				Status:      200,
			}
			json_data, _ := json.Marshal(invalidServer_resp)
			_, err := w.Write(json_data)
			if err != nil {
				log.Println("Error writing response: ", err.Error())
			}
		} else {
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
	}

}
