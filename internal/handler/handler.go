package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Oxygenss/yandex_final_project/internal/config"
	"github.com/Oxygenss/yandex_final_project/internal/models"
	"github.com/Oxygenss/yandex_final_project/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service service.Service
	cfg     *config.Config
}

func New(service service.Service, cfg config.Config) *Handler {
	return &Handler{
		service: service,
		cfg:     &cfg}
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var AuthRequest models.SignInRequest
	err = json.Unmarshal(body, &AuthRequest)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	password := h.cfg.Auth.Password
	if password == "" {
		writeJSONError(w, "Authentication is not configured", http.StatusInternalServerError)
		return
	}

	if password != AuthRequest.Password {
		writeJSONError(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"expires": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(h.cfg.Auth.Secret))
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(models.SignInResponse{Token: signedToken})
}

func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, "Identifier not specified", http.StatusBadRequest)
		return
	}

	err := h.service.DoneTask(idStr)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, "Identifier not specified", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteTask(idStr)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})

}

func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var task models.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.EditTask(task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, "Identifier not specified", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTaskByID(idStr)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var tasks []models.Task
	var err error

	if search != "" {
		tasks, err = h.service.SearchTasks(search)
	} else {
		tasks, err = h.service.GetTasks()
	}

	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetTasksResponse{Tasks: tasks})
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var task models.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.AddTask(task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AddTaskResponse{ID: id})
}

func (h *Handler) NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "неверный формат даты для параметра 'now'", http.StatusBadRequest)
		return
	}

	next, err := h.service.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
