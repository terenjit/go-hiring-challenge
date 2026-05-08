package categories

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
}
