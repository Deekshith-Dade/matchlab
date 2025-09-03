package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var Connection *sql.DB

type Config struct {
	User 			string 	`json:"user"`
	Password 	string 	`json:"password"`
	Host 			string	`json:"host"`	 
	Port 			int			`json:"port"`
	DBname 		string 	`json:"dbname"`
}

func InitDB() {
	configFile, err := os.Open("./db/dbConfig.json")
	if err != nil {
		log.Fatal("Error opening dbConfig.json:", err)
	}
	defer configFile.Close()

	var config Config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Fatal("Error decoding dbconfig.json", err)
	}
 
	connStr := fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			config.User, config.Password, config.Host, config.Port, config.DBname,
		)

	log.Printf("Connecting to %s\n", connStr)

	var connErr error
	Connection, connErr = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", connErr)	
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := Connection.PingContext(ctx); err != nil {
		Connection.Close()
		log.Fatalf("Error pinging: %v", fmt.Errorf("db.Ping: %w", err))
	}


	log.Println("database connected")
}
