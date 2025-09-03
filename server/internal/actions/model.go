package actions

import "time"

type Action struct {
	ViewerId string `json:"viewer_id"`
	ViewedId string `json:"viewed_id"`
	Kind string `json:"action"`
	At time.Time `json:"at"`
}
