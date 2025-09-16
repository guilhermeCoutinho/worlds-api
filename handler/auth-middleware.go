package handler

import (
	"context"
	"net/http"
	"strings"

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
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userId := strings.TrimPrefix(token, "Bearer ")
		r = r.WithContext(context.WithValue(r.Context(), UserIDCtxKey, userId))
		next.ServeHTTP(w, r)
	})
}
