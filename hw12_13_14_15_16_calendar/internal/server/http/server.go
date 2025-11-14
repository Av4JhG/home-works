package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	logg    *logger.Logger
	app     Application
	Addr    string
	Handler http.Handler
}

type Logger interface {
	*logger.Logger
}

type Application interface{}

type Handler struct{}

func (h *Handler) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/hello":
		fmt.Println("hello endpoint do some work")
		time.Sleep(2 * time.Second)
	case "/test":
		fmt.Println("test")
	}
}

func NewServer(logger *logger.Logger, app Application, dsn string) *Server {
	myHandler := &Handler{}
	return &Server{
		app:     app,
		Addr:    dsn,
		Handler: loggingMiddleware(myHandler, logger),
		logg:    logger,
	}
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
