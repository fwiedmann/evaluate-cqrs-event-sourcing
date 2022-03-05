package infra

import (
	"encoding/json"
	"fmt"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user-command/domain"
	"github.com/gorilla/mux"
	"net/http"
)

type UserCreateJson struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type UserCreateResponseJson struct {
	Id string `json:"id"`
}

type HttpUserHandler struct {
	UserService *domain.UserService
}

func (uh *HttpUserHandler) Create() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Got User Create Request")
		defer request.Body.Close()
		user := &UserCreateJson{}

		if err := json.NewDecoder(request.Body).Decode(user); err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		u, err := uh.UserService.Create(domain.CreateUser{
			Name:    user.Name,
			Surname: user.Surname,
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(201)
		json.NewEncoder(writer).Encode(&UserCreateResponseJson{Id: u.Id})
	}
}

type UserUpdateJson struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
type UserUpdateResponseJson struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func (uh *HttpUserHandler) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)

		fmt.Println("Got User Create Request")
		defer request.Body.Close()
		user := &UserUpdateJson{}

		if err := json.NewDecoder(request.Body).Decode(user); err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		u, err := uh.UserService.Update(domain.UpdateUser{
			Id:      vars["id"],
			Name:    user.Name,
			Surname: user.Surname,
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(200)
		json.NewEncoder(writer).Encode(&UserUpdateResponseJson{
			Id:      u.Id,
			Name:    user.Name,
			Surname: user.Surname,
		})
	}
}
