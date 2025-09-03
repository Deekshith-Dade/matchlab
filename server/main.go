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

	port := os.Getenv("SERVER_PORT")
	fmt.Println("PORT:",port)

	serverEnv := os.Getenv("SERVER_ENV")
	fmt.Println("ENV:", serverEnv)
	
	r := api.NewRouter(db.Connection)
	log.Printf("listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

