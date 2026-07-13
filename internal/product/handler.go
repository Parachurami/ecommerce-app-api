package product

import (
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

func (h *handler) CreateProduct(res http.ResponseWriter, req *http.Request) {
	var params types.CreateProductParams
	if err := utils.Read(req, &params); err != nil {
		log.Print(err.Error())
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	if err := utils.Validate.Struct(params); err != nil {
		utils.Error(res, http.StatusBadRequest, err.Error())
		return
	}
	id := req.Context().Value("userId")
	role := req.Context().Value("userRole")
	if id == nil || role == nil {
		utils.Error(res, http.StatusForbidden, "Not Authenticated")
		return
	}
	userId, idOk := id.(string)
	userRole, roleOk := role.(string)
	if !idOk || !roleOk {
		utils.Error(res, http.StatusForbidden, "Invalid token")
		return
	}
	uRole := types.Role(userRole)
	uid, userErr := uuid.Parse(userId)
	if userErr != nil {
		log.Print("zError converting id: ", userErr)
		utils.Error(res, http.StatusForbidden, "Invalid token")
		return
	}
	if uRole == types.UserRole {
		utils.Error(res, http.StatusForbidden, "Only admins can perform this action")
		return
	}
	product, productErr := h.service.CreateProduct(req.Context(), uid, &params)
	if productErr != nil {
		log.Print("Error creating product: ", productErr.Error())
		utils.Error(res, http.StatusBadRequest, productErr.Error())
		return
	}
	utils.Write(res, product)

}

func (h *handler) GetProducts(res http.ResponseWriter, req *http.Request) {
	roleObject, roleObjectOk := req.Context().Value("userRole").(string)
	idObject, idObjectOk := req.Context().Value("userId").(string)
	if !roleObjectOk || !idObjectOk {
		utils.Error(res, http.StatusForbidden, "Invalid token")
		return
	}
	userId, userIdParseErr := uuid.Parse(idObject)
	if userIdParseErr != nil {

		utils.Error(res, http.StatusForbidden, "Invalid token")
		return
	}
	userRole := types.Role(roleObject)
	if userRole != types.AdminRole {
		utils.Error(res, http.StatusForbidden, "Only admins can perform this action")
		return
	}
	products, productsQueryErr := h.service.GetProducts(req.Context(), userId)
	if productsQueryErr != nil {
		statusCode := http.StatusForbidden
		if productsQueryErr == utils.InternalServerError {
			statusCode = http.StatusInternalServerError
		}
		utils.Error(res, statusCode, productsQueryErr.Error())
		return
	}
	utils.Write(res, products)
}
