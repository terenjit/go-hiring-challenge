package catalog

type ProductResponse struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category CategoryResponse `json:"category"`
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}

type ProductDetailResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category CategoryResponse  `json:"category"`
	Variants []VariantResponse `json:"variants"`
}

type VariantResponse struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}
