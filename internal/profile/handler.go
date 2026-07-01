package profile

import (
	"errors"
	"log"
	"net/http"

	"github.com/Parachurami/ecommerce-app-api/types"
	"github.com/Parachurami/ecommerce-app-api/utils"
	"github.com/google/uuid"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) UpdateProfile(res http.ResponseWriter, req *http.Request) {
	var params types.UpdateProfileParams
	if err := utils.Read(req, &params); err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	id, ok := req.Context().Value("userId").(string)

	if !ok {
		utils.Error(res, http.StatusBadRequest, "Not authenticated")
		return
	}
	userId, parsingErr := uuid.Parse(id)
	if parsingErr != nil {
		utils.Error(res, http.StatusBadRequest, parsingErr.Error())
		return
	}
	profile, updateErr := h.service.UpdateProfile(userId, &params, req.Context())
	if updateErr != nil {
		log.Print("Error updating user: ", updateErr.Error())
		utils.Error(res, http.StatusBadRequest, errors.New("Error updating user"))
		return
	}

	if err := utils.Write(res, map[string]any{
		"email":     profile.Email,
		"bio":       profile.Bio.String,
		"image":     profile.ProfileImage.String,
		"firstName": profile.FirstName,
		"lastName":  profile.LastName,
		"updatedAt": profile.UpdatedAt,
		"createdAt": profile.CreatedAt,
	}); err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}

}

func (h *handler) GetProfile(res http.ResponseWriter, req *http.Request) {
	id, ok := req.Context().Value("userId").(string)
	if !ok {
		utils.Error(res, http.StatusForbidden, "Not Authenticated")
		return
	}
	userId, idErr := uuid.Parse(id)
	if idErr != nil {
		utils.Error(res, http.StatusForbidden, "Invalid token")
		return
	}
	profile, err := h.service.GetProfile(userId, req.Context())
	if err != nil {
		log.Print("Error fetching users: ", err)
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	utils.Write(res, map[string]any{
		"email":     profile.Email,
		"bio":       profile.Bio.String,
		"image":     profile.ProfileImage.String,
		"firstName": profile.FirstName,
		"lastName":  profile.LastName,
		"updatedAt": profile.UpdatedAt,
		"createdAt": profile.CreatedAt,
	})
}
