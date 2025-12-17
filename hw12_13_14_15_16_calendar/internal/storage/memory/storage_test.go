package memorystorage

import (
	"context"
	"os/signal"
	"syscall"
	"testing"
	"time"

	domain "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require" //nolint: depguard
)

func TestStorage(t *testing.T) {
	testEvent := domain.Event{
		ID:          1,
		Title:       "Morning Jog",
		Date:        time.Date(2025, time.January, 8, 6, 0, 0, 0, time.UTC), // 8 Jan 2025, 06:00 UTC
		Duration:    time.Hour * 1,                                          // 1 час
		Description: "A refreshing morning jog through the park.",
		UserID:      12345,
		TimeShift:   15, // Уведомление за 15 минут до события
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	t.Run("empty list", func(t *testing.T) {
		storage := New()
		list, err := storage.EventList(ctx, "", "")
		require.ErrorIs(t, err, nil)

		require.Equal(t, 0, len(list))
	})

	t.Run("not empty list", func(t *testing.T) {
		storage := New()
		storage.CreateEvent(ctx, testEvent)
		list, err := storage.EventList(ctx, "", "")
		require.ErrorIs(t, err, nil)

		require.Equal(t, 1, len(list))
	})
}
