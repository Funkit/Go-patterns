package API

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type Server struct {
	port   int
	router chi.Router
}

type Health struct {
	System1 string `json:"system_1"`
	System2 string `json:"system_2"`
}

func (h *Health) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewServer(port int) *Server {
	return &Server{
		router: chi.NewRouter(),
		port:   port,
	}
}

func (s *Server) AddMiddlewares(middlewares ...func(handler http.Handler) http.Handler) {
	s.router.Use(middlewares...)
}

func (s *Server) SubRoutes(baseURL string, r chi.Router) {
	s.router.Mount(baseURL, r)
}

func (s *Server) Run() error {
	log.Printf("Listening on port %v\n", s.port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", s.port), s.router); err != nil {
		return err
	}
	return nil
}

func (s *Server) InitializeRoutes() {
	s.router.Get("/health", s.getSystemHealth())
}

func (s *Server) getSystemHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := Health{
			System1: "OK",
			System2: "ERROR",
		}

		if err := render.Render(w, r, &health); err != nil {
			Handle(w, r, Wrap(err, ErrRender))
		}
	}
}
