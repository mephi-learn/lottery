package auth

import (
	"encoding/json"
	"homework/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecret   = "MySecret"
	jwtDuration = 6 * time.Hour
)

func Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			http.Error(w, "invalid auth header", http.StatusBadRequest)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		user := models.User{}
		if err = json.Unmarshal([]byte(claims.Subject), &user); err != nil {
			http.Error(w, "Invalid user info", http.StatusUnauthorized)
			return
		}

		ctx := user.ToContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GenerateJWTToken(userJson string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userJson,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
