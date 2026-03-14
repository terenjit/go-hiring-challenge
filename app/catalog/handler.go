package catalog

import (
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductRepository interface {
	GetAllProducts(filter models.ProductFilter) ([]models.Product, int64, error)
}
type CatalogHandler struct {
	repo ProductRepository
}

func NewCatalogHandler(r ProductRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 10

	if o := r.URL.Query().Get("offset"); o != "" {
		val, err := strconv.Atoi(o)
		if err != nil || val < 0 {
			api.ErrorResponse(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = val
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		val, err := strconv.Atoi(l)
		if err != nil || val < 1 || val > 100 {
			api.ErrorResponse(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = val
	}

	filter := models.ProductFilter{
		Offset:        offset,
		Limit:         limit,
		Category:      r.URL.Query().Get("category"),
		PriceLessThan: 0,
	}

	if p := r.URL.Query().Get("price_less_than"); p != "" {
		val, err := strconv.ParseFloat(p, 64)
		if err != nil || val <= 0 {
			api.ErrorResponse(w, http.StatusBadRequest, "invalid price_less_than")
			return
		}
		filter.PriceLessThan = val
	}

	res, total, err := h.repo.GetAllProducts(filter)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map response
	products := make([]ProductResponse, len(res))
	for i, p := range res {
		products[i] = ProductResponse{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: CategoryResponse{
				Code: p.Category.Code,
				Name: p.Category.Name,
			},
		}
	}

	// Return the products as a JSON response
	api.OKResponse(w, ProductsResponse{
		Products: products,
		Total:    total,
		Offset:   offset,
		Limit:    limit,
	})
}
