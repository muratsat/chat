package handlers

import "net/http"

func SetupHandlers() {
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

}
