package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/paxaf/workmateTest/internal/entity"
)

func (h *UsecaseHandler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	Task, err := ParseTaskFromReq(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return

	}
	h.service.Set(*Task)
	w.WriteHeader(http.StatusCreated)
}

func (h *UsecaseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	Tasks := h.service.GetAll()
	resp := entity.TaskResponse{Tasks: Tasks}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *UsecaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	key := path[2]
	err := h.service.Delete(key)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UsecaseHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	values := r.URL.Query()
	key := values.Get("id")

	Task, ok := h.service.Get(key)
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	resp := Task
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func ParseTaskFromReq(r *http.Request) (*entity.Task, error) {
	var Task entity.Task
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("invalid content type: %s", contentType)
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&Task)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	return &Task, nil
}
