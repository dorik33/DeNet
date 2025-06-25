package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dorik33/DeNet/internal/models"
	"github.com/dorik33/DeNet/internal/service"
	"github.com/dorik33/DeNet/internal/service/serviceerrors"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	RegisterHandler() http.HandlerFunc
	LoginHandler() http.HandlerFunc
	LeaderboardHandler() http.HandlerFunc
	SetReferrerHandler() http.HandlerFunc
	StatusHandler() http.HandlerFunc
	CompleteTaskHandler() http.HandlerFunc
}

type handler struct {
	userService service.UserService
	logger      *slog.Logger
}

func NewHandlers(userService service.UserService, logger *slog.Logger) Handlers {
	return &handler{
		userService: userService,
		logger:      logger,
	}
}

func (h *handler) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			h.logger.Info("Invalid method")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Error("Failed to decode request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		if req.Password != req.ConfirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		err := h.userService.Register(r.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, serviceerrors.ErrUserAlreadyExists) {
				http.Error(w, "User already exists", http.StatusConflict)
				return
			}
			http.Error(w, "Failed to register", http.StatusBadRequest)
			return
		}

		h.logger.Info("User successfully created")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User successfully created"})

	}
}

func (h *handler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			h.logger.Info("Invalid method")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		token, err := h.userService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, serviceerrors.ErrUserNotFound) {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			if errors.Is(err, serviceerrors.ErrInvalidPassword) {
				http.Error(w, "Invalid password", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func (h *handler) LeaderboardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			h.logger.Info("Invalid method")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		limit := 10
		if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
			if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
				limit = l
			}
		}

		users, err := h.userService.GetLeaderboard(r.Context(), limit)
		if err != nil {
			h.logger.Error("Failed to get leaderboard", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			h.logger.Error("Failed to encode leaderboard response", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func (h *handler) SetReferrerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			h.logger.Info("Invalid method")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userIDStr := chi.URLParam(r, "id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var req models.SetReferrerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Error("Failed to decode request body", slog.String("error", err.Error()))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = h.userService.SetReferrer(r.Context(), userID, req.ReferrerID)
		if err != nil {
			if errors.Is(err, serviceerrors.ErrUserNotFound) {
				http.Error(w, "Referrer not found", http.StatusNotFound)
				return
			}
			h.logger.Error("Failed to set referrer", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Referrer set successfully"})
	}
}

func (h *handler) StatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			h.logger.Info("Invalid method")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userIDStr := chi.URLParam(r, "id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		status, err := h.userService.Status(r.Context(), userID)
		if err != nil {
			if errors.Is(err, serviceerrors.ErrUserNotFound) {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			h.logger.Error("Failed to get status", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			h.logger.Error("Failed to encode status response", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func (h *handler) CompleteTaskHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			h.logger.Info("Invalid method")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		userIDstr := chi.URLParam(r, "id")
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		var req models.CompleteTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Error("Failed to decode request body", slog.String("error", err.Error()))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = h.userService.CompleteTask(r.Context(), userID, req.TaskID)
		if err != nil {
			if errors.Is(err, serviceerrors.ErrTaskNotFound) {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}
			if errors.Is(err, serviceerrors.ErrTaskAlreadyDone) {
				http.Error(w, "Task already completed", http.StatusConflict)
				return
			}
			h.logger.Error("Failed to complete task", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Task completed successfully"})
	}
}
