package users

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type setActiveReq struct {
	Active bool `json:"active"`
}

func Routes(db *sql.DB) http.Handler {
		r := chi.NewRouter()
		repo := NewRepository(db)

		r.Post("/", func(w http.ResponseWriter, req *http.Request){
			var u User
			if err := json.NewDecoder(req.Body).Decode(&u); err != nil {
				http.Error(w, err.Error(), 400);
			return
		}
			if err := repo.Create(req.Context(), u); err != nil {
				http.Error(w, err.Error(), 500);
			return
		}

		w.WriteHeader(http.StatusCreated)

	})

	r.Get("/", func(w http.ResponseWriter, req *http.Request){
		data, err := repo.List(req.Context())
		if err != nil {
			http.Error(w, err.Error(), 500);	
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(data)
	})

	r.Route("/{user_id}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, req *http.Request) { 
			userId := chi.URLParam(req, "user_id")

			data, err := repo.ListByID(req.Context(), userId)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if len(data) == 0 {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(data[0])

		})

		r.Patch("/active", func(w http.ResponseWriter, req *http.Request) {
			userId := chi.URLParam(req, "user_id")

			var body setActiveReq
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			affected, err := repo.SetActive(req.Context(), userId, body.Active)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if affected == 0 {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			data, err := repo.ListByID(req.Context(), userId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)	
			}
			if len(data) == 0 {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(data[0])
		})
	})

	return r
}

