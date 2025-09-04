package http

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/deekshith-dade/matchlab/internal/actions"
	"github.com/deekshith-dade/matchlab/internal/impressions"
	"github.com/deekshith-dade/matchlab/internal/matches"
	"github.com/deekshith-dade/matchlab/internal/recommendations"
	"github.com/deekshith-dade/matchlab/internal/users"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10*time.Second))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request){w.Write([]byte("ok"))})

	r.Route("/", func(api chi.Router){
		api.Mount("/users", users.Routes(db))
		api.Mount("/impressions", impressions.Routes(db))
		api.Mount("/actions", actions.Routes(db))
		api.Mount("/matches", matches.Routes(db))
		api.Mount("/recommendations", recommendations.Routes(db))
	})

	return r
}
