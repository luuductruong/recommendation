package domain

import (
	"time"

	"github.com/recommendation/services/core/context"
	"github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/helper"
)

const DefaultRecommendationNumber = 10

/*
GetRecommendationForUser return list recomment product_ids for user.
  - If user didn't get product, return most popular product.
  - If user have viewed, but not enough limit, append list popular product til fulfill limit.
  - Edge case: If not found any record in view history, just random products
*/
func (d *domain) GetRecommendationForUser(ctx context.Context, inp *product.GetRecommendationForUserInp) ([]*product.SummaryProductView, error) {
	d.logger.DebugCtx(ctx, "inp:", inp)
	if inp == nil || inp.UserID == "" {
		return d.GetPopularProducts(ctx, DefaultRecommendationNumber)
	}
	queryLimit := inp.Limit
	if queryLimit <= 0 {
		queryLimit = DefaultRecommendationNumber
	}
	recentView, err := d.userViewHistory.RecentViewProductsByUser(ctx, inp.UserID, queryLimit)
	if err != nil {
		d.logger.DebugCtx(ctx, "RecentViewProductsByUser query error:", err)
		return nil, err
	}
	if int32(len(recentView)) == queryLimit { // never len(productIds) > inp.Limit
		return recentView, nil
	}
	popular, _ := d.GetPopularProducts(ctx, queryLimit) // ignore error
	if len(popular) == 0 {
		d.logger.DebugCtx(ctx, "GetPopularProducts query error:", err)
		return recentView, nil
	}
	// merge response
	recentView = append(recentView, popular...)
	resp := helper.UniqBy(recentView, func(s *product.SummaryProductView) int64 {
		return s.ProductID
	})
	if len(resp) == 0 { // Edge case not found any history
		prods, _ := d.productRepo.Query(ctx).Limit(int(inp.Limit)).ResultList()
		if len(prods) == 0 {
			return nil, nil
		}
		for _, prod := range prods {
			resp = append(resp, &product.SummaryProductView{
				ProductID: prod.ProductID,
				ViewCount: helper.AnyToPointer(int64(0)),
				ViewAt:    helper.AnyToPointer(time.Now()),
			})
		}
	}
	if int32(len(resp)) > inp.Limit {
		resp = resp[:inp.Limit]
	}
	return resp, nil
}

/*
GetPopularProducts get products with most view
  - Get most viewed product in past 2 day (Each product must view as least 5 times)
  - If not enough limit, merge with list most view product
*/
func (d *domain) GetPopularProducts(ctx context.Context, limit int32) ([]*product.SummaryProductView, error) {
	d.logger.DebugCtx(ctx, "GetPopularProducts with limit:", limit)
	to := time.Now()
	from := to.AddDate(0, 0, -2) // default in 2 days, maybe change later
	res, err := d.userViewHistory.MostViewedInTimeRange(ctx, from, to, limit)
	if err != nil {
		d.logger.ErrorCtx(ctx, err)
		return nil, err
	}
	res = helper.SelectMap(res, func(summaryProductView *product.SummaryProductView) bool {
		return summaryProductView.ViewCount != nil && *summaryProductView.ViewCount > 5 // Example
	})
	if int32(len(res)) < limit {
		// find most viewed products
		mostView, _ := d.userViewHistory.MostView(ctx, limit)
		if len(mostView) == 0 {
			return res, nil
		}
		for _, v := range mostView {
			res = append(res, v)
		}
		res = helper.UniqBy(res, func(s *product.SummaryProductView) int64 {
			return s.ProductID
		})
		if int32(len(res)) > limit {
			res = res[:limit]
		}
	}
	return res, nil
}

/*
GetRelatedProducts
  - get
*/
func (d *domain) GetRelatedProducts(ctx context.Context, productID int64, limit int32) ([]int64, error) {
	d.logger.DebugCtx(ctx, "GetRelatedProducts with productID:", productID)
	if productID <= 0 {
		return nil, nil
	}
	prod, err := d.productRepo.Query(ctx).ByProductID(productID).Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err)
		return nil, err
	}
	if prod == nil {
		return nil, nil
	}
	viewHistory, err := d.userViewHistory.GetMostViewedProductsInCategory(ctx, prod.CategoryID, productID, limit)
	if err != nil {
		d.logger.ErrorCtx(ctx, err)
		return nil, err
	}
	relatedProducts := []int64{}
	for _, v := range viewHistory {
		relatedProducts = append(relatedProducts, v.ProductID)
	}
	return relatedProducts, nil
}
