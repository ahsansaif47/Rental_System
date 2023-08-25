package Auth

import (
	utils "Rental_System/Utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	emailUname := r.FormValue("email_uname")
	password := r.FormValue("password")

	isEmail := false
	if strings.Contains(emailUname, "@") {
		isEmail = true
	}

	fmt.Println("email_uname is: ", emailUname)
	fmt.Println("Password is: ", password)
	type loginResp struct {
		Response string `json:"resp"`
		Status   int
	}

	resp := loginResp{
		Response: "",
		Status:   200,
	}

	if utils.ConnErr != nil {
		log.Fatalln("Database connection err: ", utils.ConnErr)
	}
	encryptedPass, err := utils.Encrypt(password)
	if err != nil {
		log.Fatalln("Error encrypting password: ", err)
	}

	if hashStatus := utils.Compare_Encryption(password, encryptedPass); hashStatus {
		type User struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		user := User{
			Email:    "",
			Name:     "",
			Password: "",
		}

		var users *sql.Rows
		if isEmail {
			users, err = utils.ConnStr.Query(`SELECT email, name, password FROM users WHERE email = $1 and password = $2`,
				emailUname, encryptedPass)
			if err != nil {
				log.Fatalln("Error querying user: ", err)
			}
		} else {
			users, err = utils.ConnStr.Query(`SELECT email, name, password FROM users WHERE name = $1 and password = $2`,
				emailUname, encryptedPass)
			if err != nil {
				log.Fatalln("Error querying user: ", err)
			}
		}

		defer users.Close()

		if rows_err := utils.Rows_iteration_error_check(users); rows_err != nil {
			log.Println("Error iterating rows: ", rows_err)
		}
		for users.Next() {
			err := users.Scan(&user.Email, &user.Name, &user.Password)
			if err != nil {
				log.Println("Error setting user value: ", err)
			}
		}
		if (emailUname == user.Email || emailUname == user.Name) && encryptedPass == user.Password {
			resp.Response = "User Found"
			resp.Status = 200
			json.NewEncoder(w).Encode(resp)
		} else {
			resp.Response = "User Not Found"
			resp.Status = 200
			json.NewEncoder(w).Encode(resp)
		}
	}
}
