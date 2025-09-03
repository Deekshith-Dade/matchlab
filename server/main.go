package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	db "github.com/deekshith-dade/matchlab/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10*time.Second))

	// Routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/", listUsers)  // GET /users
		r.Post("/", createUser) // POST /users
		r.Get("/{id}", getUser) // GET /user/{id}
		r.Put("/{id}", updateUser) // PUT /users/{id}
	})

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}


func listUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := db.Connection.QueryContext(ctx, `SELECT id, x, y, active, distance FROM users`)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.X, &u.Y, &u.Active, &u.Distance); err != nil {
			httpError(w, err, http.StatusInternalServerError)
			return
		}
		out = append(out, u)

	}

	if err:=rows.Err(); err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	jsonOK(w, out)
}


func getUser(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r, "id")
	var u User
	err := db.Connection.QueryRowContext(r.Context(), `SELECT id, x, y, active, distance FROM users WHERE id = $1`, id).Scan(&u.Id, &u.X, &u.Y, &u.Active, &u.Distance)
	if err == sql.ErrNoRows {
		httpError(w, err, http.StatusNotFound)
		return
	}

	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	jsonOK(w, u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}
	_, err := db.Connection.ExecContext(
		r.Context(),
		`INSERT INTO users (id, x, y, active, distance) VALUES ($1, $2, $3, $4, $5)`,
		u.Id, u.X, u.Y, u.Active, u.Distance,
	)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"status": "created"})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}
	// optional: ensure path id wins
	u.Id = id

	res, err := db.Connection.ExecContext(
		r.Context(),
		`UPDATE users SET x=$2, y=$3, active=$4, distance=$5 WHERE id=$1`,
		u.Id, u.X, u.Y, u.Active, u.Distance,
	)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		httpError(w, context.Canceled, http.StatusNotFound)
		return
	}
	jsonOK(w, map[string]string{"status": "updated"})
}


func jsonOK(w http.ResponseWriter, v any){
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, e error, code int){
	http.Error(w, http.StatusText(code)+":"+e.Error(), code)
}
