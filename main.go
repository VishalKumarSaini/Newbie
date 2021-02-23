package main

import (
	"fmt"
	"log"
	"net/http"

	controller "Newbie/controllers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	fmt.Println("Listning......")
	r.HandleFunc("/signUp", controller.SignUpHandler).Methods("POST")
	r.HandleFunc("/signUp/auth", controller.SignUpAuthHandler).Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	r.HandleFunc("/login/auth", controller.LoginAuthHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
