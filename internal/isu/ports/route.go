package ports

import (
	"github.com/gorilla/mux"
)

func Init(server *HttpServer) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", server.CallbackHandler)
	r.HandleFunc("/auth", server.HandleRedirect)
	userApi := r.PathPrefix("/user").Subrouter()
	userApi.HandleFunc("/users/add", server.AddUser).Methods("POST")
	userApi.HandleFunc("/users/find", server.GetPublicInfoByPhoneNumber).Methods("GET")
	userApi.HandleFunc("/users/update", server.UpdatePublicInfo).Methods("PUT")
	userApi.Use(server.Auth)
	adminApi := r.PathPrefix("/admin").Subrouter()
	adminApi.HandleFunc("/users", server.AllUserData).Methods("GET")
	adminApi.HandleFunc("/users", server.UpdateFullInfo).Methods("PUT")
	adminApi.Use(server.Auth)
	return r
}
