package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/worlds-api/utils"
	"github.com/sirupsen/logrus"
)

type UserIDCtxKeyType string

const UserIDCtxKey = UserIDCtxKeyType("userID")

type AuthMiddleware struct {
	logger logrus.FieldLogger
}

func NewAuthMiddleware(logger logrus.FieldLogger) *AuthMiddleware {
	return &AuthMiddleware{logger: logger}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := utils.LoggerFromCtx(r.Context())
		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Error("No token provided")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userId := strings.TrimPrefix(token, "Bearer ")
		userIdUUID, err := uuid.Parse(userId)
		if err != nil {
			logger.WithError(err).Error("Invalid user ID")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), UserIDCtxKey, userIdUUID))
		next.ServeHTTP(w, r)
	})
}
