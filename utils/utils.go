package utils

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Parachurami/ecommerce-app-api/config"
	"github.com/Parachurami/ecommerce-app-api/types"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

var (
	InternalServerError = errors.New("Internal Server Error")
	UserNotFound        = errors.New("User Not Found")
	InvalidTokenError   = errors.New("Invalid Token")
)

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
	if req.Body == nil || req.ContentLength < 1 {
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

func GenerateJWT(id uuid.UUID, role string) (*types.TokenDetails, error) {
	accessExpiration := time.Hour * 24 * 7
	refreshExpiration := time.Hour * 24 * 30

	accessTokenId := uuid.New().String()
	refreshTokenId := uuid.New().String()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tokenId":        accessTokenId,
		"userId":         id,
		"userRole":       role,
		"expirationTime": time.Now().Add(accessExpiration),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tokenId":        refreshTokenId,
		"userId":         id,
		"userRole":       role,
		"expirationTime": time.Now().Add(refreshExpiration),
	})
	var err error = nil
	at, err := accessToken.SignedString([]byte(config.JWT_SECRET))
	rt, err := refreshToken.SignedString([]byte(config.JWT_SECRET))

	return &types.TokenDetails{
		AccessToken:  at,
		RefreshToken: rt,
		AccessUUID:   accessTokenId,
		RefreshUUID:  refreshTokenId,
		AtExpires:    accessExpiration,
		RtExpires:    refreshExpiration,
	}, err

}

func WithJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authToken := getTokenFromHeaders(req)
		token, tokenErr := verifyToken(authToken)
		if tokenErr != nil {
			log.Print("Token error: ", tokenErr)
			Error(res, http.StatusUnauthorized, "Unauthorized or Invalid Token")
			return
		}

		if !token.Valid {
			Error(res, http.StatusUnauthorized, "Invalid token")
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
		next.ServeHTTP(res, req)
	})
}

func getTokenFromHeaders(req *http.Request) string {
	cookie, err := req.Cookie("access-token")
	if cookie == nil || err != nil {
		return ""
	}
	return cookie.Value
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

func ScanRow(row pgx.Rows, product *types.Product) error {
	/*
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		    userId UUID NOT NULL,
		    name TEXT NOT NULL,
		    description TEXT NOT NULL,
		    budget FLOAT NOT NULL,
		    skills TEXT[],
		    duration INTERVAL NOT NULL,
		    expiration INTERVAL NOT NULL,
		    image_url TEXT,
		    deliverables TEXT[],
		    createdAt TIMESTAMPTZ DEFAULT NOW(),
		    updatedAt TIMESTAMPTZ DEFAULT NOW(),
	*/
	return row.Scan(&product.Id, &product.UserId, &product.Name, &product.Description,
		&product.Budget, &product.Skills, &product.Duration, &product.Expiration, &product.ImageUrl,
		&product.Deliverables, &product.CreatedAt, &product.UpdatedAt)
}
