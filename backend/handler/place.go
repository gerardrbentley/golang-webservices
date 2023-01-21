package handler

import (
	"encoding/json"
	"log"
	"net/http"

	. "github.com/gerardrbentley/places/service"
)

type PlaceError struct {
	Error string `json:"error"`
}

func PlaceHandler(service PlaceService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		searchName := r.URL.Query().Get("name")
		if searchName == "" {
			w.WriteHeader(http.StatusBadRequest)
			p := PlaceError{Error: "no name query"}
			if err := json.NewEncoder(w).Encode(p); err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		records, err := service.LookupByName(searchName)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(records); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
