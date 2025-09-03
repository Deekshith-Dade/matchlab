package matches

import "time"

type Match struct {
	U1 string `json:"u1"`
	U2 string `json:"u2"`
	At time.Time `json:"at"`
}
