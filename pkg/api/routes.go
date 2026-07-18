package api

import "net/http"

func (api API) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
