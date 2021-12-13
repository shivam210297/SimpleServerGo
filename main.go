package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"simpleServer/middleware"
	"simpleServer/services"
)

func main() {
	myroute := mux.NewRouter()
	myroute.HandleFunc("/signup", services.Signup).Methods("post")
	myroute.HandleFunc("/signin", services.Signin).Methods("post")
	myroute.Handle("/{userrole}/info", middleware.VerifyToken(http.HandlerFunc(services.Info))).Methods("Post")

	log.Fatal(http.ListenAndServe(":8080", myroute))
}
