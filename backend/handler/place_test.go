package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/gerardrbentley/places/service"
)

func TestPlaceHandler(t *testing.T) {
	ctx := context.Background()
	mockName := "seattle"

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("/place?name=%s", mockName),
		nil,
	)
	w := httptest.NewRecorder()

	s := InMemoryPlaceService{
		LookupByNameFunc: func(searchName string) ([]PlaceRecord, error) {
			return []PlaceRecord{{Agency: searchName}}, nil
		},
	}
	r := http.NewServeMux()
	r.Handle("/place", PlaceHandler(&s))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Not OK %v", w.Code)
	}
	var results []PlaceRecord
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		log.Fatalln(err)
	}
	result := results[0]
	if result.Agency != mockName {
		t.Errorf("Not same name: %v", result.Agency)
	}
}
