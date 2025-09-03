package impressions

import (
	"database/sql"
	"encoding/json"
	"net/http"
)


func Routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	repo := NewRepository(db)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		viewerId := req.URL.Query().Get("viewer_id")
		data, err := repo.List(req.Context(), viewerId)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_ = json.NewEncoder(w).Encode(data)
	})

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
		var imp Impression
		if err := json.NewDecoder(req.Body).Decode(&imp) ; err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := repo.Create(req.Context(), imp); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	return r
}
