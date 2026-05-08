package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/mytheresa/go-hiring-challenge/models"
)

func TestCategoriesRepository_GetAllCategories(t *testing.T) {
	t.Run("returns all categories", func(t *testing.T) {
		withTx(t, func(db *gorm.DB) {
			require.NoError(t, db.Create(&models.Category{Code: "test-bags", Name: "Bags"}).Error)
			require.NoError(t, db.Create(&models.Category{Code: "test-sunglasses", Name: "Sunglasses"}).Error)

			repo := models.NewCategoriesRepository(db)
			categories, err := repo.GetAllCategories()

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(categories), 2)

			codes := make([]string, len(categories))
			for i, c := range categories {
				codes[i] = c.Code
			}
			assert.Contains(t, codes, "test-bags")
			assert.Contains(t, codes, "test-sunglasses")
		})
	})
}

func TestCategoriesRepository_CreateCategory(t *testing.T) {
	t.Run("creates a category successfully", func(t *testing.T) {
		withTx(t, func(db *gorm.DB) {
			repo := models.NewCategoriesRepository(db)
			err := repo.CreateCategory(models.Category{Code: "test-jewellery", Name: "Jewellery"})

			require.NoError(t, err)

			var category models.Category
			require.NoError(t, db.Where("code = ?", "test-jewellery").First(&category).Error)
			assert.Equal(t, "Jewellery", category.Name)
		})
	})

	t.Run("returns error on duplicate category code", func(t *testing.T) {
		withTx(t, func(db *gorm.DB) {
			require.NoError(t, db.Create(&models.Category{Code: "test-hats", Name: "Hats"}).Error)

			repo := models.NewCategoriesRepository(db)
			err := repo.CreateCategory(models.Category{Code: "test-hats", Name: "Hats"})

			assert.Error(t, err)
		})
	})
}
