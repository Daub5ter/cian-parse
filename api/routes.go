package api

import (
	"cian-parse/internals/app/handlers"
	"github.com/gorilla/mux"
)

func CreateRoutes(immovableHandler *handlers.ImmovablesHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/immovables/create", immovableHandler.Create).Methods("POST")
	r.HandleFunc("/immovables/find/{id}", immovableHandler.FindOne).Methods("GET")
	r.HandleFunc("/immovables/find_all", immovableHandler.FindAll).Methods("GET")
	r.HandleFunc("/immovables/update/{id}", immovableHandler.Update).Methods("PUT")
	r.HandleFunc("/immovables/delete/{id}", immovableHandler.Delete).Methods("DELETE")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(handlers.NotFound).GetHandler()
	return r
}
