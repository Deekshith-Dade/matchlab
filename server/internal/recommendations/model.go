package recommendations

type Recommendation struct {
	UserID 	string 	`json:"user_id"`
	Rank		int 	`json:"rank"`
}
