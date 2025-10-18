package models

// WeightEntry represents a single weight record.
type WeightEntry struct {
	ID     string  `json:"id"`
	Date   string  `json:"date"` // ISO-8601 date e.g. 2025-10-18
	Weight float64 `json:"weight"`
	Notes  string  `json:"notes,omitempty"`
}
