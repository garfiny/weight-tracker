package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"weight-tracker/src/models"
	"weight-tracker/src/store"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ListWeights returns all weight entries.
func ListWeights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list := store.Default.List()
	json.NewEncoder(w).Encode(list)
}

// CreateWeight creates a new weight entry.
func CreateWeight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(r.Body)
	var entry models.WeightEntry
	if err := json.Unmarshal(body, &entry); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	entry.ID = uuid.New().String()

	store.Default.Create(entry)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

// GetWeight fetches a single weight by id.
func GetWeight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	entry, ok := store.Default.Get(id)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(entry)
}

// UpdateWeight updates an existing weight entry.
func UpdateWeight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	body, _ := io.ReadAll(r.Body)
	var updated models.WeightEntry
	if err := json.Unmarshal(body, &updated); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	up, ok := store.Default.Update(id, updated)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(up)
}

// DeleteWeight removes a weight entry.
func DeleteWeight(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if !store.Default.Delete(id) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
