package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	db "github.com/deekshith-dade/matchlab/db"
	api	 "github.com/deekshith-dade/matchlab/http"
)


type User struct {
	Id string `json:"id"` 
	X int			`json:"x"`
	Y int     `json:"y"`
	Active bool `json:"active"`
	Distance int `json:"distance"`
}


func main() {

	 db.InitDB()
	 defer db.Connection.Close()
	
	if err := db.ClearAllTables(db.Connection); err != nil {
		log.Fatalf("failed to clear tables: %v", err)
	}
	fmt.Println("Cleared Tables to Start Simulation Fresh")

	port := os.Getenv("SERVER_PORT")
	fmt.Println("PORT:",port)

	serverEnv := os.Getenv("SERVER_ENV")
	fmt.Println("ENV:", serverEnv)
	
	r := api.NewRouter(db.Connection)
	log.Printf("listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

