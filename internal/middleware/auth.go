package middleware

import (
	"fmt"
	"github.com/ananaslegend/news-crud/internal/contexts"
	"github.com/ananaslegend/news-crud/pkg/jwt"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
)

var (
	ErrInvalidAccessToken = fmt.Errorf("invalid access token")
)

func Auth(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ParseToken(headerParts[1], []byte(secret))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(contexts.SetUserID(r.Context(), claims.UserID))

		next(w, r)
	}
}
