package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Payload struct {
	id int
	jwt.StandardClaims
}
type yty struct {
	Token string
}

func Signup(w http.ResponseWriter, r *http.Request) {

	err := godotenv.Load()
	if err != nil {
		http.Error(w, "internal server error", 500)
	}
	psqlconn := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))
	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		http.Error(w, "internal servor error", 500)
	}
	var cred credentials

	err = json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		http.Error(w, "Bad Request Error", 400)
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Srever Error", 500)
	}

	sql := `insert into user_info (email,password) values($1,$2)`
	_, err = db.Exec(sql, cred.Email, string(hashedpassword))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	//fmt.Println(result.RowsAffected())
	fmt.Println("inserted successfully")
	defer db.Close()

}

func Signin(w http.ResponseWriter, r *http.Request) {
	var cred credentials
	err := godotenv.Load()
	psqlconn := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))

	err = json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		http.Error(w, "Bad Request Error", 400)
	}
	// open database
	var temppass string
	var id int
	db, _ := sql.Open("postgres", psqlconn)
	sql := `select password,user_id from user_info where email=$1`
	result := db.QueryRow(sql, cred.Email)
	err = result.Scan(&temppass, &id)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	err = bcrypt.CompareHashAndPassword([]byte(temppass), []byte(cred.Password))
	if err != nil {
		http.Error(w, "Not a valid user", 401)
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	payload := &Payload{id, jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte("hello"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	//u := &yty{tokenString}
	//w.Write([]byte(u.token))
	err = json.NewEncoder(w).Encode(&yty{Token: tokenString})
	if err != nil {
		http.Error(w, "internal server error", 500)
	}
}

func Info(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct{ Text string }{"Its working"})
}
