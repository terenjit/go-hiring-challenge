package categories

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoryRepository interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category models.Category) error
}
type CategoryHandler struct {
	repo CategoryRepository
}

func NewCategoryHandler(r CategoryRepository) *CategoryHandler {
	return &CategoryHandler{
		repo: r,
	}
}

func (h *CategoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	res, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map response
	categories := make([]CategoryResponse, len(res))
	for i, p := range res {
		categories[i] = CategoryResponse{
			Code: p.Code,
			Name: p.Name,
		}
	}

	// Return the products as a JSON response
	api.OKResponse(w, categories)
}
func (h *CategoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Code == "" || req.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "code and name are required")
		return
	}

	category := models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.CreateCategory(category); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			api.ErrorResponse(w, http.StatusConflict, "category code already exists")
			return
		}
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.CreatedResponse(w, CategoryResponse{
		Code: category.Code,
		Name: category.Name,
	})
}
