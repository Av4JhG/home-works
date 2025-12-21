package memorystorage

import (
	"context"
	"strconv"
	"sync"

	domain "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/model"
)

type Storage struct {
	events map[string]domain.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: map[string]domain.Event{},
		mu:     sync.RWMutex{},
	}
}

func (s *Storage) Connect(_ context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	// TODO
	return nil
}

func (s *Storage) CreateEvent(_ context.Context, event domain.Event) error {
	s.mu.Lock()
	event.ID = int64(len(s.events) + 1)
	s.events[strconv.Itoa(int(event.ID))] = event
	s.mu.Unlock()
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, event domain.Event) error {
	s.mu.Lock()
	s.events[strconv.Itoa(int(event.ID))] = event
	s.mu.Unlock()
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id int64) error {
	delete(s.events, strconv.Itoa(int(id)))
	return nil
}

func (s *Storage) EventList(_ context.Context, _ string, _ string) ([]domain.Event, error) {
	list := make([]domain.Event, 0, len(s.events))
	s.mu.Lock()
	for _, v := range s.events {
		// fmt.Println("list ", v)
		list = append(list, v)
	}
	s.mu.Unlock()
	return list, nil
}

func (s *Storage) ClearEvents(_ context.Context, _ string) error {
	// TODO
	return nil
}
