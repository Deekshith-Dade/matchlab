package users

type User struct {
	ID		string 	`json:"id"`
	X			int 	`json:"x"`
	Y			int	`json:"y"`
	Active bool		`json:"active"`
	Distance int	`json:"distance"`

}
