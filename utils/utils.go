package utils

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Parachurami/ecommerce-app-api/config"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

var Psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func Write(res http.ResponseWriter, data any) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	return json.NewEncoder(res).Encode(map[string]any{
		"message": "success",
		"data":    data,
	})
}

func Read(req *http.Request, data any) error {
	if req.Body == nil {
		return errors.New("Body cannot be empty")
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func Error(res http.ResponseWriter, errCode int, message any) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(errCode)
	json.NewEncoder(res).Encode(map[string]any{
		"message": message,
	})
}

func GenerateJWT(id uuid.UUID, role string) (string, error) {
	expiration := time.Hour * 24 * 7
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":         id,
		"userRole":       role,
		"expirationTime": time.Now().Add(expiration),
	})

	return token.SignedString([]byte(config.JWT_SECRET))

}

func WithJWT(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		authToken := getTokenFromHeaders(req)
		token, tokenErr := verifyToken(authToken)
		if tokenErr != nil {
			Error(res, http.StatusForbidden, "Not Authenticated")
			return
		}

		if !token.Valid {
			Error(res, http.StatusForbidden, "Invalid token")
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		id, idOk := claims["userId"].(string)
		userRole := claims["userRole"].(string)
		expirationTime := claims["expirationTime"].(string)
		actualTime, timeErr := time.Parse(time.RFC3339Nano, expirationTime)
		if timeErr != nil {
			log.Print("Error converting time")
			return
		}
		if actualTime.Compare(time.Now()) < 1 {
			Error(res, http.StatusForbidden, "Token Has Expired")
			return
		}
		if !idOk {
			log.Print("Error getting user id")
			return
		}
		ctx := req.Context()
		ctx = context.WithValue(ctx, "userId", id)
		ctx = context.WithValue(ctx, "userRole", userRole)
		req = req.WithContext(ctx)
		handlerFunc(res, req)
	}
}

func getTokenFromHeaders(req *http.Request) string {
	token := req.Header.Get("Authorization")
	if token == "" {
		return ""
	}
	return strings.Split(token, " ")[1]
}

func verifyToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(config.JWT_SECRET), nil
	})
}

func HashPassword(password string) (string, error) {
	token, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func CompareHash(password string, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		log.Print("Error in hash comparison", err)
		return errors.New("Invalid credentials")
	}
	return nil
}
