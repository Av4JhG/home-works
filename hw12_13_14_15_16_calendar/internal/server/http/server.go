package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/app"
	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/logger"
	domain "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/model"
)

type Server struct {
	logg    *logger.Logger
	app     Application
	Addr    string
	Handler http.Handler
	storage app.Storage
}

type Logger interface {
	*logger.Logger
}

type Application interface{}

func NewServer(logger *logger.Logger, app Application, storage app.Storage, dsn string) *Server {
	mux := http.NewServeMux()
	server := &Server{
		app:     app,
		Addr:    dsn,
		logg:    logger,
		storage: storage,
	}
	mux.HandleFunc("/events", server.GetEventList)
	mux.HandleFunc("/events/create", server.CreateEvent)
	mux.HandleFunc("/events/update", server.UpdateEvent)
	mux.HandleFunc("/events/delete", server.DeleteEvent)

	server.Handler = loggingMiddleware(mux, logger)
	return server
}

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:         s.Addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      s.Handler,
	}

	s.logg.Infoln("start server on ", s.Addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		s.logg.Errorln(err)
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	os.Exit(1)
	return nil
}

func (s *Server) GetEventList(w http.ResponseWriter, r *http.Request) {
	// dateQuery := r.URL.Query().Get("date")
	// if dateQuery == "" {
	// 	http.Error(w, "Missing date parameter", http.StatusBadRequest)
	// 	return
	// }

	// date, err := time.Parse("2006-01-02", dateQuery)
	// if err != nil {
	// 	http.Error(w, "Invalid date format", http.StatusBadRequest)
	// 	return
	// }

	events, err := s.storage.EventList(r.Context())
	if err != nil {
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		return
	}

	results := make([]domain.Event, 0)
	results = append(results, events...)
	// for _, event := range events {
	// if event.Date.Format("2006-01-02") == date.Format("2006-01-02") {
	// results = append(results, event)
	// }
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req domain.Event
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.storage.CreateEvent(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var req domain.Event
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.storage.UpdateEvent(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.storage.DeleteEvent(r.Context(), req.ID)
	if err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Event deleted"}`))
}
