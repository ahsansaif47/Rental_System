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
	confirm_pass := r.FormValue("confirmation")
	user_name := r.FormValue("uname")

	type registration_response struct {
		RegResponse string
		Status      int
	}
	if user_name != "" && password != "" && email != "" && confirm_pass != "" {
		isEmail_valid := validate_email(email)
		if !isEmail_valid {
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
					if confirm_pass == password {
						encryptedPass, err := utils.Encrypt(password)
						if err != nil {
							log.Println("Error encrypting password: ", err.Error())
						} else {
							insert_user_query, err := utils.ConnStr.Prepare(`Insert into users(email, password) Values($1, $2)`)
							if err != nil {
								log.Println("Error preparing query: ", err.Error())
							}
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
					} else {
						pass_resp := registration_response{
							RegResponse: "Passwords do not match",
							Status:      200,
						}
						json_data, _ := json.Marshal(pass_resp)
						_, err := w.Write(json_data)
						if err != nil {
							log.Println("Error writing response: ", err.Error())
						}
					}
				}
			}
		}
	} else {
		emptyResp := registration_response{
			RegResponse: "Either of your field is empty",
			Status:      200,
		}
		json_data, _ := json.Marshal(emptyResp)
		_, err := w.Write(json_data)
		if err != nil {
			log.Println("Error writing response: ", err.Error())
		}
	}
}
