package api

import "net/http"

func (api API) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{uuid}", api.UserController.NewUser)
	mux.HandleFunc("GET /user/{uuid}", api.UserController.GetUser)
	mux.HandleFunc("PATCH /user/{uuid}", api.UserController.UpdateUser)
	mux.HandleFunc("DELETE /user/{uuid}", api.UserController.DeleteUser)
	return mux
}
