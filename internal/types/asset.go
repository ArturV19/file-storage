package types

import "time"

type Asset struct {
	Name         string    `json:"name"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	CreatedAt    time.Time `json:"created_at"`
}
