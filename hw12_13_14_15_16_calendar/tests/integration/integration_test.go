package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	domain "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/assert"
)

var baseURL = "http://localhost:8888"

// func TestMain(m *testing.M) {
// 	waitForServer()

// 	code := m.Run()
// 	os.Exit(code)
// }

// func waitForServer() {
// 	client := &http.Client{Timeout: 2 * time.Second}

// 	for i := 0; i < 10; i++ {
// 		fmt.Println("Try connect to server", i, " str: ", baseURL)
// 		resp, err := client.Get(baseURL + "/events")
// 		if err == nil && resp.StatusCode != http.StatusInternalServerError {
// 			return
// 		}
// 		time.Sleep(2 * time.Second)
// 	}
// 	panic("Server did not start in time")
// }

var event = domain.Event{
	ID:          1,
	Title:       "Тестовое событие",
	Date:        time.Now(),
	Duration:    2 * time.Hour,
	Description: "Это описание тестового события",
	UserID:      123,
	TimeShift:   0,
}

var event2 = domain.Event{
	ID:          1,
	Title:       "Обновленная встреча",
	Date:        time.Now(),
	Duration:    2 * time.Hour,
	Description: "Это описание тестового события",
	UserID:      123,
	TimeShift:   0,
}

func TestEventAPI(t *testing.T) {
	client := &http.Client{}
	ctx := context.Background()

	body, _ := json.Marshal(event)
	req, _ := http.NewRequestWithContext(ctx, "POST", baseURL+"/events/create", bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	var createdEvent map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&createdEvent)
	// fmt.Println(createdEvent)
	assert.Equal(t, event.Title, createdEvent["Title"])

	req, _ = http.NewRequestWithContext(ctx, "GET", baseURL+"/events", nil)
	resp, err = client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var events []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&events)
	assert.NotEmpty(t, events)

	body, _ = json.Marshal(event2)
	req, _ = http.NewRequestWithContext(ctx, "PUT", baseURL+"/events/update", bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	deleteReq := map[string]interface{}{
		"id": createdEvent["id"],
	}
	body, _ = json.Marshal(deleteReq)
	req, _ = http.NewRequestWithContext(ctx, "DELETE", baseURL+"/events/delete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}
