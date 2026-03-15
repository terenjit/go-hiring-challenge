package models

type PageFilter struct {
	Offset int
	Limit  int
}

type ProductFilter struct {
	Page          PageFilter
	Category      string
	PriceLessThan float64
}
