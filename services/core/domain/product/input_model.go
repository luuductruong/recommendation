package product

type GetProductDetailInp struct {
	UserID    string
	ProductID int64
}

type GetRecommendationForUserInp struct {
	UserID string
	Limit  int32
}
