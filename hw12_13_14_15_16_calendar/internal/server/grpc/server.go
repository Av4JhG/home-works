package grpc

import (
	"context"
	"sync"
	"time"

	pb "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/calendar/pb"
	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/app"
	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/logger"
	domain "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedEventsServer
	logg    *logger.Logger
	mu      sync.Mutex
	app     Application
	storage app.Storage
}

type Application interface{}

func NewServer(logger *logger.Logger, app Application, storage app.Storage) *Server {
	return &Server{
		app:     app,
		storage: storage,
		logg:    logger,
	}
}

func (s *Server) EventList(ctx context.Context, req *pb.EventListRequest) (*pb.EventListResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var results []*pb.Event
	events, err := s.storage.EventList(ctx)
	if err != nil {
		s.logg.Error("cannot get events")
	}
	for _, event := range events {
		if event.Date.Format("2006-01-02") == req.Date.AsTime().Format("2006-01-02") {
			results = append(results, &pb.Event{
				Id:          event.ID,
				Title:       event.Title,
				Date:        timestamppb.New(event.Date),
				Duration:    timestamppb.New(time.Unix(0, int64(event.Duration))),
				Description: event.Description,
				UserId:      event.UserID,
				TimeShift:   event.TimeShift,
			})
		}
	}
	return &pb.EventListResponse{Results: results}, nil
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// на случай если захочу сериализаторы вынести
	event := &pb.Event{
		Title:       req.Title,
		Date:        req.Date,
		Duration:    req.Duration,
		Description: req.Description,
		UserId:      req.UserId,
		TimeShift:   req.TimeShift,
	}
	err := s.storage.CreateEvent(ctx, domain.Event{
		Title:       event.Title,
		Date:        event.Date.AsTime(),
		Duration:    event.Duration.AsTime().Sub(time.Unix(0, 0)),
		Description: event.Description,
		UserID:      event.UserId,
		TimeShift:   event.TimeShift,
	})
	if err != nil {
		s.logg.Error("cannot create event")
	}

	return &pb.CreateEventResponse{Event: event}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event := &pb.Event{
		Id:          req.Event.Id,
		Title:       req.Event.Title,
		Date:        req.Event.Date,
		Duration:    req.Event.Duration,
		Description: req.Event.Description,
		UserId:      req.Event.UserId,
		TimeShift:   req.Event.TimeShift,
	}
	err := s.storage.UpdateEvent(ctx, domain.Event{
		ID:          event.Id,
		Title:       event.Title,
		Date:        event.Date.AsTime(),
		Duration:    event.Duration.AsTime().Sub(time.Unix(0, 0)),
		Description: event.Description,
		UserID:      event.UserId,
		TimeShift:   event.TimeShift,
	})
	if err != nil {
		s.logg.Error("cannot update event")
	}

	return &pb.UpdateEventResponse{Event: req.Event}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage.DeleteEvent(ctx, req.Id)
	return &emptypb.Empty{}, nil
}
