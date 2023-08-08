package Auth

import (
	utils "Rental_System/Utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	emailUname := r.FormValue("email_uname")
	password := r.FormValue("password")

	fmt.Println("email_uname is: ", emailUname)
	fmt.Println("Password is: ", password)
	type loginResp struct {
		Resp   string
		Status int
	}

	resp := loginResp{
		Resp:   "",
		Status: 200,
	}

	if utils.ConnErr == nil {
		encryptedPass, err := utils.Encrypt(password)
		if err == nil {
			fmt.Println("No Error till encryption")
			if hashStatus := utils.Compare_Encryption(password, encryptedPass); hashStatus {
				user_query, err := utils.ConnStr.Prepare(`SELECT * FROM users WHERE (email = $1) OR (name = $1) AND password = $2`)
				if err != nil {
					log.Println("Error preparing query: ", err.Error())
				} else {
					users, err := user_query.Query(emailUname, encryptedPass)
					fmt.Println("Querying user")
					if err == nil {
						rowsIter_err := utils.Rows_iteration_error_check(users)
						if rowsIter_err == nil {
							usersCount := utils.Count_rows(users)
							fmt.Println("User count is: ", usersCount)
							if usersCount > 0 {
								fmt.Println("In here")
								resp.Resp = "Found User"
								resp.Status = 200
								json_resp, _ := json.Marshal(resp)
								_, err := w.Write(json_resp)
								if err != nil {
									log.Println("Error writing response: ", err.Error())
								}
								fmt.Println("Found user")
							} else {
								resp.Resp = "User Not Found!"
								resp.Status = 404
								json_resp, _ := json.Marshal(resp)
								_, err := w.Write(json_resp)
								if err != nil {
									log.Println("Error writing response: ", err.Error())
								}
								fmt.Println("User not found!")
							}
						} else {
							log.Println("Error iterating rows: ", rowsIter_err.Error())
						}
					} else {
						log.Println("Error Querying: ", err.Error())
					}
				}
			}
		} else {
			log.Println("Error encrypting password: ", err.Error())
		}

	} else {
		log.Fatal("PostgreSQL Connection Error: ", utils.ConnErr.Error())
	}
}
