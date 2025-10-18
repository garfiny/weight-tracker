package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"weight-tracker/src/store"
)

// setupRouter mirrors main.go route wiring used by tests
func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HandleRequest).Methods("GET")
	r.HandleFunc("/weights", ListWeights).Methods("GET")
	r.HandleFunc("/weights", CreateWeight).Methods("POST")
	r.HandleFunc("/weights/{id}", GetWeight).Methods("GET")
	r.HandleFunc("/weights/{id}", UpdateWeight).Methods("PUT")
	r.HandleFunc("/weights/{id}", DeleteWeight).Methods("DELETE")
	return r
}

// doRequest executes a request against the router and returns the recorder
func doRequest(r *mux.Router, method, path string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestWeightsRouter_TableDriven(t *testing.T) {
	store.Default.Reset()
	r := setupRouter()

	t.Run("happy-path CRUD", func(t *testing.T) {
		// create
		payload := map[string]interface{}{"date": "2025-10-18", "weight": 72.5, "notes": "morning"}
		b, _ := json.Marshal(payload)
		rec := doRequest(r, http.MethodPost, "/weights", b)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create expected 201, got %d", rec.Code)
		}
		var created map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
			t.Fatalf("decode create: %v", err)
		}
		id, _ := created["id"].(string)
		if id == "" {
			t.Fatalf("expected id in created response")
		}

		// list
		rec = doRequest(r, http.MethodGet, "/weights", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("list expected 200, got %d", rec.Code)
		}

		// get
		rec = doRequest(r, http.MethodGet, "/weights/"+id, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("get expected 200, got %d", rec.Code)
		}

		// update
		updated := map[string]interface{}{"date": "2025-10-18", "weight": 73.0, "notes": "updated"}
		ub, _ := json.Marshal(updated)
		rec = doRequest(r, http.MethodPut, "/weights/"+id, ub)
		if rec.Code != http.StatusOK {
			t.Fatalf("update expected 200, got %d", rec.Code)
		}

		// delete
		rec = doRequest(r, http.MethodDelete, "/weights/"+id, nil)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("delete expected 204, got %d", rec.Code)
		}

		// get after delete
		rec = doRequest(r, http.MethodGet, "/weights/"+id, nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("get after delete expected 404, got %d", rec.Code)
		}
	})

	t.Run("invalid json -> 400", func(t *testing.T) {
		store.Default.Reset()
		rec := doRequest(r, http.MethodPost, "/weights", []byte("not-json"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid json, got %d", rec.Code)
		}
	})

	t.Run("not found operations -> 404", func(t *testing.T) {
		store.Default.Reset()
		id := "does-not-exist"
		rec := doRequest(r, http.MethodGet, "/weights/"+id, nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("get expected 404, got %d", rec.Code)
		}
		rec = doRequest(r, http.MethodPut, "/weights/"+id, []byte(`{"date":"2025-10-18","weight":70}`))
		if rec.Code != http.StatusNotFound {
			t.Fatalf("put expected 404, got %d", rec.Code)
		}
		rec = doRequest(r, http.MethodDelete, "/weights/"+id, nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("delete expected 404, got %d", rec.Code)
		}
	})
}
