package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"hello-world/pkg/users"
	"net/http"
	"os"
)

type apiHandler func(http.ResponseWriter, *http.Request)

func main() {
	address := ":8080"
	errs := make(chan error)

	userStore, err := users.InitUserStore("test-user-table", "us-east-1")
	if err != nil {
		os.Exit(1)
	}

	svc, err := users.CreateUserService(userStore)
	if err != nil {
		os.Exit(2)
	}

	router := mux.NewRouter()
	router.Handle("/user", apiHandler(svc.AddRequest)).Methods("POST")
	go func() {
		fmt.Println("transport", "HTTP", "address", address)
		errs <- http.ListenAndServe(address, router)
	}()
	fmt.Println("fatal", "exit", "error", <-errs)
}

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn(w, r)
}
