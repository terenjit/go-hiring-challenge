package catalog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockProductRepo struct {
	products  []models.Product
	product   models.Product
	total     int64
	err       error
	detailErr error
}

func (m *mockProductRepo) GetAllProducts(filter models.ProductFilter) ([]models.Product, int64, error) {
	return m.products, m.total, m.err
}

func (m *mockProductRepo) GetProductDetailByCode(code string) (models.Product, error) {
	return m.product, m.detailErr
}

func TestHandleGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockProductRepo{
			products: []models.Product{
				{
					Code:  "PROD001",
					Price: decimal.NewFromFloat(10.99),
					Category: models.Category{
						Code: "clothing",
						Name: "Clothing",
					},
				},
				{
					Code:  "PROD002",
					Price: decimal.NewFromFloat(12.49),
					Category: models.Category{
						Code: "shoes",
						Name: "Shoes",
					},
				},
			},
			total: 2,
			err:   nil,
		}

		handler := NewCatalogHandler(repo)
		req := httptest.NewRequest("GET", "/catalog", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.JSONEq(t, `{
			"products": [
				{"code": "PROD001", "price": 10.99, "category": {"code": "clothing", "name": "Clothing"}},
				{"code": "PROD002", "price": 12.49, "category": {"code": "shoes", "name": "Shoes"}}
			],
			"total": 2,
			"offset": 0,
			"limit": 10
		}`, rec.Body.String())

	})

	t.Run("invalid limit", func(t *testing.T) {
		repo := &mockProductRepo{}
		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest("GET", "/catalog?limit=200", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "invalid limit"}`, rec.Body.String())
	})

	t.Run("invalid offset", func(t *testing.T) {
		repo := &mockProductRepo{}
		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest("GET", "/catalog?offset=-1", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "invalid offset"}`, rec.Body.String())
	})

	t.Run("error from db", func(t *testing.T) {
		repo := &mockProductRepo{
			err: errors.New("database error"),
		}
		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest("GET", "/catalog", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error": "database error"}`, rec.Body.String())
	})

	t.Run("filter by category returns filtered products", func(t *testing.T) {
		repo := &mockProductRepo{
			products: []models.Product{
				{
					Code:  "PROD002",
					Price: decimal.NewFromFloat(12.49),
					Category: models.Category{
						Code: "shoes",
						Name: "Shoes",
					},
				},
			},
			total: 1,
		}
		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest("GET", "/catalog?category=shoes", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
        "products": [
            {"code": "PROD002", "price": 12.49, "category": {"code": "shoes", "name": "Shoes"}}
        ],
        "total": 1,
        "offset": 0,
        "limit": 10
    }`, rec.Body.String())
	})

	t.Run("filter by price_less_than returns filtered products", func(t *testing.T) {
		repo := &mockProductRepo{
			products: []models.Product{
				{
					Code:  "PROD006",
					Price: decimal.NewFromFloat(5.50),
					Category: models.Category{
						Code: "shoes",
						Name: "Shoes",
					},
				},
			},
			total: 1,
		}
		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest("GET", "/catalog?price_less_than=10", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
        "products": [
            {"code": "PROD006", "price": 5.5, "category": {"code": "shoes", "name": "Shoes"}}
        ],
        "total": 1,
        "offset": 0,
        "limit": 10
    }`, rec.Body.String())
	})
}

func TestHandleByCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockProductRepo{
			product: models.Product{
				Code:  "PROD001",
				Price: decimal.NewFromFloat(10.99),
				Category: models.Category{
					Code: "clothing",
					Name: "Clothing",
				},
				Variants: []models.Variant{
					{Name: "Variant A", SKU: "SKU001A", Price: decimal.NewFromFloat(11.99)},
					{Name: "Variant B", SKU: "SKU001B", Price: decimal.NewFromFloat(0)},
				},
			},
		}
		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")
		rec := httptest.NewRecorder()
		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
			"code": "PROD001",
			"price": 10.99,
			"category": {"code": "clothing", "name": "Clothing"},
			"variants": [
				{"name": "Variant A", "sku": "SKU001A", "price": 11.99},
				{"name": "Variant B", "sku": "SKU001B", "price": 10.99}
			]
		}`, rec.Body.String())
	})

	t.Run("error from db", func(t *testing.T) {
		repo := &mockProductRepo{
			detailErr: errors.New("database error"),
		}
		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")
		rec := httptest.NewRecorder()
		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error": "database error"}`, rec.Body.String())
	})

	t.Run("product not found", func(t *testing.T) {
		repo := &mockProductRepo{
			detailErr: gorm.ErrRecordNotFound,
		}
		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")
		rec := httptest.NewRecorder()
		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"error": "product not found"}`, rec.Body.String())
	})

}
