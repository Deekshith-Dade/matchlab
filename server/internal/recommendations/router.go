package recommendations

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)


func Routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	repo := NewRepository(db)


	r.Route("/{user_id}", func(api chi.Router){
		api.Get("/", func(w http.ResponseWriter, req *http.Request){
			userId := chi.URLParam(req, "user_id") 
			topKStr := req.URL.Query().Get("topk")
			topK, err := strconv.Atoi(topKStr)
			if err != nil {
				http.Error(w, "invalid topk parameter", http.StatusBadRequest)
			}

			data, err := repo.getRecommendations(req.Context(), userId, topK)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return

			}

			_ = json.NewEncoder(w).Encode(data)
		
		})
	})

	return r
}
