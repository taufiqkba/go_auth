package main

import (
	"fmt"
	"github.com/taufiqkba/go_auth/controllers"
	"net/http"
)

func main() {
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/login", controllers.Login)

	fmt.Println("Server running on port: 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
