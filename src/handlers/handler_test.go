package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"weight-tracker/src/store"
)

func TestHandleRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	HandleRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["message"] != "Hello, World!" {
		t.Fatalf("unexpected message: %v", resp["message"])
	}
}

func TestStoreReset(t *testing.T) {
	// ensure Reset exists and clears store
	store.Default.Reset()
	list := store.Default.List()
	if len(list) != 0 {
		t.Fatalf("expected empty store after reset, got %d items", len(list))
	}
}
