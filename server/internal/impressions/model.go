package impressions

import "time"

type Impression struct {
	ViewerID string		`json:"viewer_id"`	
	ViewedID string 	`json:"viewed_id"`
	Rank 		 int 			`json:"rank"`
	At 			 time.Time `json:"at"`
}
