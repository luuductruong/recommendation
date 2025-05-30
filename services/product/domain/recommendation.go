package domain

import (
	"time"

	"github.com/recommendation/services/core/context"
	"github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/helper"
)

const MinRecommendationProduct = 10

func (d *domain) GetRecommendationForUser(ctx context.Context, inp *product.GetRecommendationForUserInp) ([]int64, error) {
	d.logger.DebugCtx(ctx, "inp:", inp)
	if inp == nil || inp.UserID == "" {
		d.logger.DebugCtx(ctx, "inp with nil userID")
		return nil, nil
	}
	productIds := make([]int64, 0)
	recentView, err := d.userViewHistory.RecentViewProductsByUser(ctx, inp.UserID, inp.Limit)
	if err != nil {
		d.logger.DebugCtx(ctx, "RecentViewProductsByUser query error:", err)
		return nil, err
	}
	for _, v := range recentView {
		productIds = append(productIds, v.ProductID)
	}
	if int32(len(productIds)) == inp.Limit { // never len(productIds) > inp.Limit
		return productIds, nil
	}
	popular, _ := d.GetPopularProducts(ctx, MinRecommendationProduct) // ignore error
	if len(popular) == 0 {
		d.logger.DebugCtx(ctx, "GetPopularProducts query error:", err)
		return productIds, nil
	}
	// merge response
	productIds = append(productIds, popular...)
	resp := helper.Unique(productIds)
	if int32(len(resp)) > inp.Limit {
		resp = resp[:inp.Limit]
	}
	return resp, nil
}

// GetPopularProducts get product with most view
func (d *domain) GetPopularProducts(ctx context.Context, limit int32) ([]int64, error) {
	d.logger.DebugCtx(ctx, "GetPopularProducts with limit:", limit)
	to := time.Now()
	from := to.AddDate(0, 0, -2) // default in 2 days, maybe change later
	res, err := d.userViewHistory.MostPopularProductsInTimeRange(ctx, from, to, limit)
	if err != nil {
		d.logger.ErrorCtx(ctx, err)
		return nil, err
	}
	popular := []int64{}
	for _, v := range res {
		popular = append(popular, v.ProductID)
	}
	return popular, nil
}
