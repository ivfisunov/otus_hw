package internalhttp

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	*app.App
	server *http.Server
	router *mux.Router
}

func NewServer(app *app.App, host string, port string) *Server {
	router := mux.NewRouter()

	addr := net.JoinHostPort(host, port)
	server := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return &Server{app, server, router}
}

func (s *Server) Start(ctx context.Context) error {
	s.router.Use(s.loggingMiddleware)
	s.router.HandleFunc("/hello", s.helloHandler).Methods("GET")

	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"message": "hello"})
	if err != nil {
		s.Logger.Error("error sending response")
	}
}
