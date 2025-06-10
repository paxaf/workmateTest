package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/paxaf/BrandScoutTest/internal/controller"
	"github.com/paxaf/BrandScoutTest/internal/entity"
)

type MockUsecase struct {
	tasks      map[string]entity.Task
	keyCounter atomic.Uint64
	returnErr  bool
}

func (m *MockUsecase) Get(key string) (entity.Task, bool) {
	if m.returnErr {
		return entity.Task{}, false
	}
	return m.tasks[key], true
}

func (m *MockUsecase) Set(Task entity.Task) {
	if m.returnErr {
		return
	}
	key := strconv.FormatUint(m.keyCounter.Add(1), 10)
	Task.Id = key
	m.tasks[key] = Task
}

func (m *MockUsecase) GetAll() []entity.Task {
	if m.returnErr {
		return nil
	}
	tasks := make([]entity.Task, 0, len(m.tasks))
	for _, q := range m.tasks {
		tasks = append(tasks, q)
	}
	return tasks
}

func (m *MockUsecase) Delete(id string) error {
	if m.returnErr {
		return errors.New("mock error")
	}
	if _, exists := m.tasks[id]; !exists {
		return errors.New("not found")
	}
	delete(m.tasks, id)
	return nil
}

func TestAddHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{tasks: make(map[string]entity.Task)}
		h := controller.New(mockUsecase)

		Task := entity.Task{Title: "Do sometning", Content: "Okay"}
		body, _ := json.Marshal(Task)
		req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Add(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}
		if len(mockUsecase.tasks) != 1 {
			t.Fatalf("Task not added to service")
		}

		for _, q := range mockUsecase.tasks {
			if q.Title != "Do sometning" || q.Content != "Okay" {
				t.Errorf("Unexpected Task content: %+v", q)
			}
			if q.Id == "" {
				t.Error("Task ID not generated")
			}
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		t.Parallel()
		h := controller.UsecaseHandler{}
		req := httptest.NewRequest(http.MethodGet, "/add", nil)
		w := httptest.NewRecorder()

		h.Add(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		t.Parallel()
		h := controller.UsecaseHandler{}
		req := httptest.NewRequest(http.MethodPost, "/add", nil)
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		h.Add(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		t.Parallel()
		h := controller.UsecaseHandler{}
		req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Add(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestGetAllHandler(t *testing.T) {
	t.Parallel()

	t.Run("success with tasks", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{
			tasks: map[string]entity.Task{
				"1": {Id: "1", Title: "Do test for workmate", Content: "Well well well"},
			},
		}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected JSON content, got %s", ct)
		}

		var resp entity.TaskResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if len(resp.Tasks) != 1 {
			t.Fatalf("Expected 1 Task, got %d", len(resp.Tasks))
		}
		if resp.Tasks[0].Id != "1" || resp.Tasks[0].Content != "Well well well" {
			t.Errorf("Unexpected Task data: %+v", resp.Tasks[0])
		}
	})

	t.Run("empty storage", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{tasks: make(map[string]entity.Task)}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp entity.TaskResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if len(resp.Tasks) != 0 {
			t.Errorf("Expected 0 tasks, got %d", len(resp.Tasks))
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		t.Parallel()
		h := controller.UsecaseHandler{}
		req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})
}
func TestDeleteHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{
			tasks: map[string]entity.Task{
				"test-id": {Id: "test-id", Title: "Do test for workmate", Content: "Well well well"},
			},
		}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodDelete, "/delete/test-id", nil)
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		if _, exists := mockUsecase.tasks["test-id"]; exists {
			t.Errorf("Task was not deleted")
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{tasks: make(map[string]entity.Task)}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodDelete, "/delete/missing-id", nil)
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{returnErr: true}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodDelete, "/delete/id", nil)
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("missing id in path", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{tasks: make(map[string]entity.Task)}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodDelete, "/delete/", nil)
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing id, got %d", w.Code)
		}
	})
}

func TestParseQuoteFromReq(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		Task := entity.Task{Title: "Do test for workmate", Content: "Well well well"}
		body, _ := json.Marshal(Task)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		result, err := controller.ParseTaskFromReq(req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.Title != Task.Title || result.Content != Task.Content {
			t.Errorf("Parsed Task doesn't match original")
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "text/plain")

		_, err := controller.ParseTaskFromReq(req)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("POST", "/", strings.NewReader("{invalid}"))
		req.Header.Set("Content-Type", "application/json")

		_, err := controller.ParseTaskFromReq(req)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})

	t.Run("empty body", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set("Content-Type", "application/json")

		_, err := controller.ParseTaskFromReq(req)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}
