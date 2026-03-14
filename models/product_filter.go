package models

type ProductFilter struct {
	Offset        int     `json:"offset"`
	Limit         int     `json:"limit"`
	Category      string  `json:"category"`
	PriceLessThan float64 `json:"price_less_than"`
}
