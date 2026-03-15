package categories

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

type mockCategoryRepo struct {
	categories []models.Category
	category   models.Category
	total      int64
	err        error
	detailErr  error
}

func (m *mockCategoryRepo) GetAllCategories() ([]models.Category, error) {
	return m.categories, m.err
}

func (m *mockCategoryRepo) CreateCategory(category models.Category) error {
	return m.err
}

func TestHandleGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockCategoryRepo{
			categories: []models.Category{
				{
					Code: "clothing",
					Name: "Clothing",
				},
				{
					Code: "shoes",
					Name: "Shoes",
				},
			},
		}
		handler := NewCategoryHandler(repo)

		req := httptest.NewRequest("GET", "/categories", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[
			{"code": "clothing", "name": "Clothing"},
			{"code": "shoes", "name": "Shoes"}
		]`, rec.Body.String())
	})

	t.Run("error from db", func(t *testing.T) {
		repo := &mockCategoryRepo{
			err: errors.New("database error"),
		}
		handler := NewCategoryHandler(repo)

		req := httptest.NewRequest("GET", "/categories", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error": "database error"}`, rec.Body.String())
	})

	t.Run("empty response", func(t *testing.T) {
		repo := &mockCategoryRepo{
			categories: []models.Category{},
		}
		handler := NewCategoryHandler(repo)

		req := httptest.NewRequest("GET", "/categories", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[]`, rec.Body.String())
	})
}

func TestHandleCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockCategoryRepo{
			err: nil,
		}

		body := CreateCategoryRequest{
			Code: "clothing",
			Name: "Clothing",
		}

		jsonBody, _ := json.Marshal(body)
		handler := NewCategoryHandler(repo)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(jsonBody))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"code": "clothing", "name": "Clothing"}`, rec.Body.String())
	})

	t.Run("empty request body", func(t *testing.T) {
		repo := &mockCategoryRepo{
			err: nil,
		}
		handler := NewCategoryHandler(repo)

		req := httptest.NewRequest("POST", "/categories", nil)
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "invalid request body"}`, rec.Body.String())
	})

	t.Run("empty request body", func(t *testing.T) {
		repo := &mockCategoryRepo{
			err: nil,
		}

		body := CreateCategoryRequest{
			Code: "",
			Name: "",
		}

		jsonBody, _ := json.Marshal(body)
		handler := NewCategoryHandler(repo)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(jsonBody))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "code and name are required"}`, rec.Body.String())
	})

	t.Run("duplicate code", func(t *testing.T) {
		repo := &mockCategoryRepo{
			err: errors.New("duplicate key"),
		}

		body := CreateCategoryRequest{
			Code: "clothing",
			Name: "Clothing",
		}

		jsonBody, _ := json.Marshal(body)
		handler := NewCategoryHandler(repo)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(jsonBody))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.JSONEq(t, `{"error": "category code already exists"}`, rec.Body.String())
	})

	t.Run("error from db", func(t *testing.T) {
		repo := &mockCategoryRepo{
			err: errors.New("database error"),
		}

		body := CreateCategoryRequest{
			Code: "clothing",
			Name: "Clothing",
		}

		jsonBody, _ := json.Marshal(body)
		handler := NewCategoryHandler(repo)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(jsonBody))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error": "database error"}`, rec.Body.String())
	})

}
