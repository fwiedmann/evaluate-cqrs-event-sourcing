package main

import (
	"fmt"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user/query/domain"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user/query/infra"
	"net/http"
)

func main() {
	service, err := domain.NewUserService(infra.NewUserRepositoryBoltDB())
	if err != nil {
		panic(err)
	}

	infra.NewEventBusNats(service).Subscribe()
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
