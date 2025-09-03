package matches

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)


func Routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	repo := NewRepository(db)

	r.Get("/", func(w http.ResponseWriter, req *http.Request){
		userId := req.URL.Query().Get("user_id")
		data, err := repo.ListByUser(req.Context(), userId)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_ = json.NewEncoder(w).Encode(data)
	})

	r.Post("/", func(w http.ResponseWriter, req *http.Request){
		var m Match
		if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), 400);
			return
		}
		if err := repo.Create(req.Context(), m); err != nil {
			http.Error(w, err.Error(), 500);
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	return r
}
