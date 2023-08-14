package model

import "time"

type Sentinel struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Hosts     string `json:"hosts"`
	CreatedAt time.Time
}
