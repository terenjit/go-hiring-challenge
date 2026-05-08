package catalog

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAllProducts(filter models.ProductFilter) ([]models.Product, int64, error)
	GetProductDetailByCode(code string) (models.Product, error)
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
		Page: models.PageFilter{
			Offset: offset,
			Limit:  limit,
		},
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

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid code")
		return
	}

	res, err := h.repo.GetProductDetailByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.ErrorResponse(w, http.StatusNotFound, "product not found")
			return
		}
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map response
	variants := make([]VariantResponse, len(res.Variants))
	for i, variant := range res.Variants {
		price := variant.Price.InexactFloat64()
		if variant.Price.IsZero() {
			price = res.Price.InexactFloat64()
		}
		variants[i] = VariantResponse{
			Name:  variant.Name,
			SKU:   variant.SKU,
			Price: price,
		}
	}
	productDetail := ProductDetailResponse{
		Code:  res.Code,
		Price: res.Price.InexactFloat64(),
		Category: CategoryResponse{
			Code: res.Category.Code,
			Name: res.Category.Name,
		},
		Variants: variants,
	}

	api.OKResponse(w, productDetail)
}
