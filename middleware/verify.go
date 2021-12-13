package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

type payload struct {
	Id int
	jwt.StandardClaims
}
type ErrorMsg struct {
	Err string
}

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		//if header != " " {
		//	bearer := strings.Split(header, " ")
		//	if len(bearer) == 2 {
		payLoad := &payload{}
		token, error := jwt.ParseWithClaims(header, payLoad, func(token *jwt.Token) (interface{}, error) {
			return []byte("hello"), nil
		})

		if error != nil {
			json.NewEncoder(w).Encode(ErrorMsg{error.Error()})
			return
		}
		
		if token.Valid {
			err := godotenv.Load()
			if err != nil {
				http.Error(w, "internal server error", 500)
			}
			psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))
			// open database
			db, err := sql.Open("postgres", psqlconn)
			if err != nil {
				http.Error(w, "internal servor error", 500)
			}
			sql := `select role from user_role where user_id=$1`
			result, err := db.Query(sql, payLoad.Id)
			if err != nil {
				http.Error(w, "internal server error", 500)
			}
			var role string
			err = result.Scan(&role)
			if err != nil {
				http.Error(w, "internal server error", 500)
			}
			vars := mux.Vars(r)
			Userrole := vars["userrole"]
			if Userrole != role {
				json.NewEncoder(w).Encode(ErrorMsg{"Unauthorised User"})
				next.ServeHTTP(w, r)
			}
		} else {
			json.NewEncoder(w).Encode(ErrorMsg{"Unauthorised user"})
		}

		//	} else {
		//		json.NewEncoder(w).Encode(ErrorMsg{"invalid token"})
		//	}
		//} else {
		//	json.NewEncoder(w).Encode(ErrorMsg{"header cannot be empty"})

	})
}
