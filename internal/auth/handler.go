package userAuth

import (
	"net/http"

	"github.com/Parachurami/ecommerce-app-api/types"
	"github.com/Parachurami/ecommerce-app-api/utils"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
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
	token, jwtErr := utils.GenerateJWT(user.Id.UUID, string(user.UserRole))
	if jwtErr != nil {
		utils.Error(res, http.StatusBadRequest, jwtErr.Error())
	}
	if err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	utils.Write(res, map[string]any{
		"user":  user,
		"token": token,
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
	token, jwtErr := utils.GenerateJWT(user.Id.UUID, string(user.UserRole))
	if jwtErr != nil {
		utils.Error(res, http.StatusBadRequest, jwtErr.Error())
		return
	}
	utils.Write(res, map[string]any{
		"user": map[string]any{
			"id":        user.Id,
			"email":     user.Email,
			"firstName": profile.FirstName,
			"lastName":  profile.LastName,
			"bio":       profile.Bio.String,
		},
		"token": token,
	})

}
