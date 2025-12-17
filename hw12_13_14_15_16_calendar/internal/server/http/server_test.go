package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

type Request struct {
	Date     string `json:"date"`
	ListType string `json:"listType"`
}

func (m *MockStorage) ClearEvents(ctx context.Context, _ string) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockStorage) EventList(ctx context.Context, date string, listType string) ([]domain.Event, error) {
	args := m.Called(ctx, date, listType)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockStorage) CreateEvent(ctx context.Context, event domain.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockStorage) UpdateEvent(ctx context.Context, event domain.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockStorage) DeleteEvent(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStorage) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockStorage) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHttpServer(t *testing.T) {
	mockStorage := new(MockStorage)
	server := NewServer(nil, nil, mockStorage, "")

	t.Run("TestGetEventList", func(t *testing.T) {
		events := []domain.Event{{ID: 1, Title: "Test Event", Date: time.Date(2025, time.February, 2, 6, 8, 23, 0, time.UTC)}}
		mockStorage.On("EventList", mock.Anything, "2025-02-21", "week").Return(events, nil)

		url := "/events?date=2025-02-21&listType=week"
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		server.GetEventList(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("TestCreateEvent", func(t *testing.T) {
		event := domain.Event{Title: "New Event", Date: time.Date(2025, time.February, 2, 6, 8, 23, 0, time.UTC)}
		mockStorage.On("CreateEvent", mock.Anything, event).Return(nil)

		body, _ := json.Marshal(event)
		r := httptest.NewRequest("POST", "/events/create", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.CreateEvent(w, r)

		resp := w.Result()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("TestDeleteEvent", func(t *testing.T) {
		mockStorage.On("DeleteEvent", mock.Anything, int64(1)).Return(nil)
		body := `{"id": 1}`

		r := httptest.NewRequest("POST", "/events/delete", bytes.NewReader([]byte(body)))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.DeleteEvent(w, r)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("TestCreateEventInvalidBody", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/events/create", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		server.CreateEvent(w, r)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
