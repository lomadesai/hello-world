package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"hello-world/pkg/users"
	"net/http"
)

type apiHandler func(http.ResponseWriter, *http.Request)

func main() {
	address := ":8080"
	errs := make(chan error)

	router := mux.NewRouter()
	router.Handle("/user", apiHandler(users.AddRequest)).Methods("POST")
	go func() {
		fmt.Println("transport", "HTTP", "address", address)
		errs <- http.ListenAndServe(address, router)
	}()
	fmt.Println("fatal", "exit", "error", <- errs)
}

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn(w, r)
}