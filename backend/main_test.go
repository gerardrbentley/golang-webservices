package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/gerardrbentley/places/service"
)

func TestLookupPlaceRoute(t *testing.T) {
	s := InMemoryPlaceService{
		LookupByNameFunc: func(searchName string) ([]PlaceRecord, error) {
			return []PlaceRecord{{Agency: searchName}}, nil
		},
	}
	h := setupHandlers(&s)

	w := httptest.NewRecorder()
	mockName := "seattle"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/place?name=%s", mockName), nil)
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Not OK %v", w.Code)
	}
	results := []PlaceRecord{}
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		log.Fatalln(err)
	}
	result := results[0]
	if result.Agency != mockName {
		t.Errorf("Not same Name: %v", result.Agency)
	}
}
