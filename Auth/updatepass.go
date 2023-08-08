package Auth

import (
	utils "Rental_System/Utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func UpdatePass(w http.ResponseWriter, r *http.Request) {
	old_pass := r.FormValue("old_pass")
	new_pass := r.FormValue("new_pass")
	confirm_pass := r.FormValue("confirm_pass")

	fmt.Println("Old Password is: ", old_pass)
	fmt.Println("New password is: ", new_pass)
	fmt.Println("Confirmation is: ", confirm_pass)

	type UpdatePassResponse struct {
		Response string
		Status   int
	}

	resp := UpdatePassResponse{
		Response: "",
		Status:   200,
	}

	if utils.ConnErr == nil {
		pass_query, err := utils.ConnStr.Prepare(`SELECT password from users WHERE password = $1`)
		curr_pass, err := pass_query.Query(pass_query)
		fmt.Println("Querying user")

		if err == nil {
			pass := ""
			for curr_pass.Next() {
				curr_pass.Scan(&pass)
			}

			old_encryption, err := utils.Encrypt(old_pass)
			if err == nil {
				if old_encryption == pass {
					if new_pass == confirm_pass {
						new_pass_encryption, err := utils.Encrypt(new_pass)
						if err == nil {
							udpate_pass_query, err := utils.ConnStr.Prepare(`UPDATE users SET password = $1 WHERE password = $2`)
							if err == nil {
								_, err := udpate_pass_query.Query(new_pass_encryption)
								if err == nil {
									resp.Response = "Password Updated Sucessfully"
									resp.Status = 200
									json_resp, err := json.Marshal(resp)
									if err == nil {
										w.Write(json_resp)
									} else {
										resp.Response = "Error marshalling response"
										resp.Status = 500
									}
								} else {
									log.Println("Error updating password: ", err)
									resp.Response = "Error updating password"
									resp.Status = 500
									json_resp, err := json.Marshal(resp)
									if err == nil {
										w.Write(json_resp)
									} else {
										resp.Response = "Error marshalling response"
										resp.Status = 500
									}
								}
							} else {
								resp.Response = "Error preparing query"
								resp.Status = 500
								json_resp, err := json.Marshal(resp)
								if err == nil {
									w.Write(json_resp)
								} else {
									resp.Response = "Error marshalling response"
									resp.Status = 500
								}
							}
						} else {
							resp.Response = "Error encrypting password"
							resp.Status = 500
							json_resp, err := json.Marshal(resp)
							if err == nil {
								w.Write(json_resp)
							} else {
								resp.Response = "Error marshalling response"
								resp.Status = 500
							}
						}
					}
				}
			} else {
				log.Println("Error encrypting old password: ", err)
				resp.Response = "Error encrypting ols password"
				resp.Status = 500
				json_resp, err := json.Marshal(resp)
				if err == nil {
					w.Write(json_resp)
				} else {
					resp.Response = "Error marshalling response"
					resp.Status = 500
				}
			}
		}
	} else {
		log.Println("Database connection Err: ", utils.ConnErr)
		resp.Response = "Database connection error"
		resp.Status = 500
		json_resp, err := json.Marshal(resp)
		if err == nil {
			w.Write(json_resp)
		} else {
			resp.Response = "Error marshalling response"
			resp.Status = 500
		}
	}
}
