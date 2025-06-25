package jwt

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/dorik33/DeNet/internal/config"
	"github.com/dorik33/DeNet/internal/utills"
	"github.com/go-chi/chi/v5"
)

func AuthMiddleware(logger *slog.Logger, cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing Authorization header")
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("Invalid authorization format", slog.String("header", authHeader))
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			userID, err := utills.ValidateToken(token, []byte(cfg.SecretKey))
			if err != nil {
				logger.Warn("Invalid token", slog.String("token", token), slog.String("error", err.Error()))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			urlUserID := chi.URLParam(r, "id")
			if urlUserID != "" {
				if strconv.Itoa(userID) != urlUserID {
					logger.Warn("Unauthorized access attempt", slog.Int("user_id", userID), slog.String("requested_id", urlUserID))
					http.Error(w, "You are not authorized to access this resource.", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
