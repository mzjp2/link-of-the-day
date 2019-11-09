package storage

import "time"

// Record is the representation of a storage row
type Record struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Count     int       `json:"count"`
	DateAdded time.Time `json:"date_added"`
	Scheduled time.Time `json:"scheduled"`
}

// Service defines a storage service for URLs
type Service interface {
	Save(string, time.Time, time.Time) (int64, error)
	Load(int) (*Record, error)
	LoadScheduled(time.Time) (*Record, error)
	LoadLast() (*Record, error)
	UpdateCount(int) error
	Close() error
}
