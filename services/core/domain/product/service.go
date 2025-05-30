package product

import "github.com/recommendation/services/core/context"

type Service interface {
	GetProductDetail(ctx context.Context, inp *GetProductDetailInp) (*Product, error)
	GetRecommendationForUser(ctx context.Context, inp *GetRecommendationForUserInp) ([]int64, error)
}
