package product

type GetProductDetailInp struct {
	UserID    string
	ProductID int64
}

type GetRecommendationForUserInp struct {
	UserID    string
	ProductID int64
	Limit     int32
}

type CreateProductInp struct {
	Name       string
	Price      float64
	CategoryID string
}
