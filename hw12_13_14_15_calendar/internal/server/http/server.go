package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	*app.App
	server *http.Server
	router *mux.Router
}

type DateBody struct {
	Date string `json:"date"`
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
	s.router.HandleFunc("/event", s.createEvent).Methods("POST")
	s.router.HandleFunc("/event", s.updateEvent).Methods("PUT")
	s.router.HandleFunc("/event/{id}", s.deleteEvent).Methods("DELETE")
	s.router.HandleFunc("/events/day", s.listEventDay).Methods("POST")
	s.router.HandleFunc("/events/week", s.listEventWeek).Methods("POST")
	s.router.HandleFunc("/events/month", s.listEventMonth).Methods("POST")

	s.Logger.Info(fmt.Sprintf("rest server started on %s", s.server.Addr))
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Logger.Info("rest server is stoped")
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

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	event := storage.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		s.Logger.Error("body parser error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.Storage.CreateEvent(event)
	if err != nil {
		s.Logger.Error("event creating error: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "event created successfuly"})
	if err != nil {
		s.Logger.Error("error sending response")
	}
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	event := storage.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		s.Logger.Error("body parser error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := event.ID
	err = s.Storage.UpdateEvent(id, event)
	if err != nil {
		s.Logger.Error("event updating error: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "event updated successfuly"})
	if err != nil {
		s.Logger.Error("error sending response")
	}
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	if sId, ok := mux.Vars(r)["id"]; ok {
		id, err := strconv.Atoi(sId)
		if err != nil {
			s.Logger.Error("error converting id: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = s.Storage.DeleteEvent(id)
		if err != nil {
			s.Logger.Error("event deleting error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(map[string]string{"message": "event deleted successfuly"})
		if err != nil {
			s.Logger.Error("error sending response")
		}
		return
	}

	s.Logger.Error("body parser error")
	http.Error(w, "id required", http.StatusBadRequest)
}

func (s *Server) listEventDay(w http.ResponseWriter, r *http.Request) {
	body := DateBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		s.Logger.Error("body parser error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	day, err := time.Parse("2006-01-02T15:04:05-0700", body.Date)
	if err != nil {
		s.Logger.Error("error parsing date: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := s.Storage.ListEventDay(day)
	if err != nil {
		s.Logger.Error("event searching error: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.Logger.Error("error sending response")
	}
}

func (s *Server) listEventWeek(w http.ResponseWriter, r *http.Request) {
	body := DateBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		s.Logger.Error("body parser error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	day, err := time.Parse("2006-01-02T15:04:05-0700", body.Date)
	if err != nil {
		s.Logger.Error("error parsing date: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := s.Storage.ListEventWeek(day)
	if err != nil {
		s.Logger.Error("event searching error: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.Logger.Error("error sending response")
	}
}

func (s *Server) listEventMonth(w http.ResponseWriter, r *http.Request) {
	body := DateBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		s.Logger.Error("body parser error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	day, err := time.Parse("2006-01-02T15:04:05-0700", body.Date)
	if err != nil {
		s.Logger.Error("error parsing date: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := s.Storage.ListEventMonth(day)
	if err != nil {
		s.Logger.Error("event searching error: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.Logger.Error("error sending response")
	}
}
