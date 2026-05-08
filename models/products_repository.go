package models

import (
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(filter ProductFilter) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{})

	if filter.Category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", filter.Category)
	}

	if filter.PriceLessThan > 0 {
		query = query.Where("products.price < ?", filter.PriceLessThan)
	}
	query.Count(&total)

	if err := query.Debug().Preload("Category").
		Offset(filter.Page.Offset).
		Limit(filter.Page.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductsRepository) GetProductDetailByCode(code string) (Product, error) {
	var product Product
	if err := r.db.Debug().Preload("Category").Preload("Variants").Where("code = ?", code).First(&product).Error; err != nil {
		return Product{}, err
	}
	return product, nil
}
