package userAuth

import (
	"log"
	"net/http"
	"time"

	"github.com/Parachurami/ecommerce-app-api/types"
	"github.com/Parachurami/ecommerce-app-api/utils"
	"github.com/redis/go-redis/v9"
)

type handler struct {
	service Service
	db      *redis.Client
}

func NewHandler(service Service, db *redis.Client) *handler {
	return &handler{
		service: service,
		db:      db,
	}
}

func (h *handler) LoginUser(res http.ResponseWriter, req *http.Request) {
	var params types.LoginUserParams
	if err := utils.Read(req, &params); err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	if err := utils.Validate.Struct(params); err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.Login(&params)
	if err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	td, jwtErr := utils.GenerateJWT(user.Id.UUID, string(user.UserRole))
	if jwtErr != nil {
		utils.Error(res, http.StatusBadRequest, jwtErr.Error())
		return
	}
	if err := h.db.Set(req.Context(), td.RefreshUUID, td.RefreshToken, td.RtExpires).Err(); err != nil {
		log.Print("Error storing token in redis: ", err)
		utils.Error(res, http.StatusInternalServerError, utils.InternalServerError.Error())
		return
	}
	log.Print("Token saved to redis")
	http.SetCookie(res, &http.Cookie{
		Name:     "refresh-token",
		Value:    td.RefreshToken,
		Expires:  time.Now().Add(td.RtExpires),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
		MaxAge:   int(td.RtExpires),
	})
	http.SetCookie(res, &http.Cookie{
		Name:     "access-token",
		Value:    td.AccessToken,
		Expires:  time.Now().Add(td.AtExpires),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(td.AtExpires),
	})
	http.SetCookie(res, &http.Cookie{
		Name:     "refresh-token-id",
		Value:    td.RefreshUUID,
		Expires:  time.Now().Add(td.AtExpires),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(td.AtExpires),
	})
	utils.Write(res, map[string]any{
		"user": user,
	})

}

func (h *handler) RegisterUser(res http.ResponseWriter, req *http.Request) {
	var params types.RegisterUserParams
	if err := utils.Read(req, &params); err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	if err := utils.Validate.Struct(params); err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	password, err := utils.HashPassword(params.Password)
	if err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	user, profile, registrationErr := h.service.Register(&types.RegisterUserParams{FirstName: params.FirstName,
		LastName: params.LastName,
		Email:    params.Email,
		RoleType: params.RoleType,
		Password: password})
	if registrationErr != nil {
		utils.Error(res, http.StatusBadRequest, registrationErr.Error())
		return
	}
	td, jwtErr := utils.GenerateJWT(user.Id.UUID, string(user.UserRole))
	if jwtErr != nil {
		utils.Error(res, http.StatusBadRequest, jwtErr.Error())
		return
	}
	if err := h.db.Set(req.Context(), td.RefreshUUID, td.RefreshToken, td.RtExpires).Err(); err != nil {
		log.Print("Error storing token in redis: ", err)
		utils.Error(res, http.StatusInternalServerError, utils.InternalServerError.Error())
		return
	}
	log.Print("Token saved to redis")
	http.SetCookie(res, &http.Cookie{
		Name:     "refresh-token",
		Value:    td.RefreshToken,
		Expires:  time.Now().Add(td.RtExpires),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
		MaxAge:   int(td.RtExpires),
	})
	http.SetCookie(res, &http.Cookie{
		Name:     "access-token",
		Value:    td.AccessToken,
		Expires:  time.Now().Add(td.AtExpires),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(td.AtExpires),
	})
	http.SetCookie(res, &http.Cookie{
		Name:     "refresh-token-id",
		Value:    td.RefreshUUID,
		Expires:  time.Now().Add(td.AtExpires),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(td.AtExpires),
	})
	utils.Write(res, map[string]any{
		"user": map[string]any{
			"id":        user.Id,
			"email":     user.Email,
			"firstName": profile.FirstName,
			"lastName":  profile.LastName,
			"bio":       profile.Bio.String,
		},
	})
}

func (h *handler) Logout(res http.ResponseWriter, req *http.Request) {
	cookie, cookieErr := req.Cookie("refresh-token-id")
	if cookieErr != nil || cookie == nil {
		log.Print("Error retreiving cookie: ", cookieErr)
		utils.Error(res, http.StatusUnauthorized, "Unauthorized or Invalid Token")
		return
	}
	if err := h.db.Del(req.Context(), cookie.Value).Err(); err != nil {
		log.Print("Error deleting token in redis: ", err)
		utils.Error(res, http.StatusInternalServerError, utils.InternalServerError.Error())
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:     "refresh-token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
		MaxAge:   -1,
	})
	http.SetCookie(res, &http.Cookie{
		Name:     "access-token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.SetCookie(res, &http.Cookie{
		Name:     "refresh-token-id",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	utils.Write(res, "User Logged Out")
}
