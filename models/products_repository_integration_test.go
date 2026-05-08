package models_test

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/mytheresa/go-hiring-challenge/models"
)

func codesOf(products []models.Product) []string {
	codes := make([]string, len(products))
	for i, p := range products {
		codes[i] = p.Code
	}
	return codes
}

func TestProductsRepository_GetAllProducts(t *testing.T) {
	t.Run("returns products with no filter", func(t *testing.T) {
		withTx(t, func(db *gorm.DB) {
			category := models.Category{Code: "test-clothing", Name: "Clothing"}
			require.NoError(t, db.Create(&category).Error)

			require.NoError(t, db.Create(&models.Product{Code: "SHIRT-001", Price: decimal.NewFromFloat(49.99), CategoryID: category.ID}).Error)
			require.NoError(t, db.Create(&models.Product{Code: "JACKET-001", Price: decimal.NewFromFloat(149.99), CategoryID: category.ID}).Error)

			repo := models.NewProductsRepository(db)
			products, total, err := repo.GetAllProducts(models.ProductFilter{
				Page: models.PageFilter{Offset: 0, Limit: 100},
			})

			require.NoError(t, err)
			assert.GreaterOrEqual(t, total, int64(2))
			assert.Subset(t, codesOf(products), []string{"SHIRT-001", "JACKET-001"})
		})
	})

	t.Run("filters by category code", func(t *testing.T) {
		withTx(t, func(db *gorm.DB) {
			shoes := models.Category{Code: "test-shoes", Name: "Shoes"}
			clothing := models.Category{Code: "test-clothing", Name: "Clothing"}
			require.NoError(t, db.Create(&shoes).Error)
			require.NoError(t, db.Create(&clothing).Error)

			require.NoError(t, db.Create(&models.Product{Code: "SNEAKER-001", Price: decimal.NewFromFloat(89.99), CategoryID: shoes.ID}).Error)
			require.NoError(t, db.Create(&models.Product{Code: "TSHIRT-001", Price: decimal.NewFromFloat(29.99), CategoryID: clothing.ID}).Error)

			repo := models.NewProductsRepository(db)
			products, _, err := repo.GetAllProducts(models.ProductFilter{
				Page:     models.PageFilter{Offset: 0, Limit: 10},
				Category: "test-shoes",
			})

			require.NoError(t, err)
			codes := codesOf(products)
			assert.Contains(t, codes, "SNEAKER-001")
			assert.NotContains(t, codes, "TSHIRT-001")
		})
	})

	t.Run("filters by price less than", func(t *testing.T) {
		withTx(t, func(db *gorm.DB) {
			category := models.Category{Code: "test-accessories", Name: "Accessories"}
			require.NoError(t, db.Create(&category).Error)

			require.NoError(t, db.Create(&models.Product{Code: "BELT-001", Price: decimal.NewFromFloat(35.00), CategoryID: category.ID}).Error)
			require.NoError(t, db.Create(&models.Product{Code: "WATCH-001", Price: decimal.NewFromFloat(599.00), CategoryID: category.ID}).Error)

			repo := models.NewProductsRepository(db)
			products, _, err := repo.GetAllProducts(models.ProductFilter{
				Page:          models.PageFilter{Offset: 0, Limit: 100},
				PriceLessThan: 100.00,
			})

			require.NoError(t, err)
			codes := codesOf(products)
			assert.Contains(t, codes, "BELT-001")
			assert.NotContains(t, codes, "WATCH-001")
		})
	})

	t.Run("respects pagination limit", func(t *testing.T) {
		withTx(t, func(db *gorm.DB) {
			category := models.Category{Code: "test-shoes-page", Name: "Shoes"}
			require.NoError(t, db.Create(&category).Error)

			for i := 1; i <= 5; i++ {
				require.NoError(t, db.Create(&models.Product{
					Code:       fmt.Sprintf("BOOT-%03d", i),
					Price:      decimal.NewFromFloat(float64(i) * 20),
					CategoryID: category.ID,
				}).Error)
			}

			repo := models.NewProductsRepository(db)
			products, _, err := repo.GetAllProducts(models.ProductFilter{
				Page: models.PageFilter{Offset: 0, Limit: 2},
			})

			require.NoError(t, err)
			assert.Len(t, products, 2)
		})
	})
}

func TestProductsRepository_GetProductDetailByCode(t *testing.T) {
	t.Run("returns product with category and variants preloaded", func(t *testing.T) {
		withTx(t, func(db *gorm.DB) {
			category := models.Category{Code: "test-shoes", Name: "Shoes"}
			require.NoError(t, db.Create(&category).Error)

			product := models.Product{Code: "SNEAKER-002", Price: decimal.NewFromFloat(119.99), CategoryID: category.ID}
			require.NoError(t, db.Create(&product).Error)

			require.NoError(t, db.Create(&models.Variant{ProductID: product.ID, Name: "Size 40", SKU: "SNEAKER-002-40", Price: decimal.NewFromFloat(109.99)}).Error)
			require.NoError(t, db.Create(&models.Variant{ProductID: product.ID, Name: "Size 42", SKU: "SNEAKER-002-42"}).Error)

			repo := models.NewProductsRepository(db)
			result, err := repo.GetProductDetailByCode("SNEAKER-002")

			require.NoError(t, err)
			assert.Equal(t, "SNEAKER-002", result.Code)
			assert.Equal(t, "test-shoes", result.Category.Code)
			assert.Len(t, result.Variants, 2)
		})
	})

	t.Run("returns error when product code does not exist", func(t *testing.T) {
		withTx(t, func(db *gorm.DB) {
			repo := models.NewProductsRepository(db)
			_, err := repo.GetProductDetailByCode("DOES-NOT-EXIST")

			assert.Error(t, err)
		})
	})
}
