package main

import (
	"fmt"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user-query/domain"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user-query/infra"
	"net/http"
)

func main() {
	_, err := domain.NewUserService(infra.NewUserRepositoryBoltDB(), infra.NewEventBusNats())
	if err != nil {
		panic(err)
	}

	//handler := infra.HttpUserHandler{
	//	UserService: service,
	//}
	//
	//router := mux.NewRouter()
	//router.HandleFunc("/users/", handler.Create())
	//router.HandleFunc("/users/{id}", handler.Update())

	fmt.Println("Started Server")
	panic(http.ListenAndServe(":8080", nil))
}
