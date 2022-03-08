package main

import (
	"fmt"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user/command/domain"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user/command/infra"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	service := domain.NewUserService(infra.NewUserRepositoryBoltDB(), infra.NewEventBusNats())

	handler := infra.HttpUserHandler{
		UserService: service,
	}

	router := mux.NewRouter()
	router.HandleFunc("/users/", handler.Create())
	router.HandleFunc("/users/{id}", handler.Update())

	fmt.Println("Started Server")
	panic(http.ListenAndServe(":8080", router))
}
