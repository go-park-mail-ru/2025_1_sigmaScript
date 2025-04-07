package router

import "net/http"

type AuthHandlerInterface interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Session(w http.ResponseWriter, r *http.Request)
}

type StaffPersonHandlerInterface interface {
	GetPerson(w http.ResponseWriter, r *http.Request)
}
