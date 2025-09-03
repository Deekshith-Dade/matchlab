package actions

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(db *sql.DB) http.Handler {
	r := ch.NewRouter()
	repo := NewRepository(db)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		viewerId := req.URL.Query().Get("viewer_id")
		data, err := repo.List(req.Context(), viewerId)
		if err != nil { http.Error(w, err.Error(), 500); return}
		_ = json.NewEncoder(w).Encode(data)
	})

	r.Post("/", func(w http.ResponseWriter, req *http.Request){
		var act Action
		if err := json.NewDecoder(req.Body).Decode(&act); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := repo.Create(req.Context(), act); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})



	return r
}
