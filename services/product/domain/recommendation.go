package domain

import (
	"errors"
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
	if inp.ProductID > 0 {
		return d.GetRelatedProducts(ctx, inp.ProductID, queryLimit)
	}
	return d.GetCollaborativeRecommendation(ctx, inp.UserID, queryLimit)
}

/*
GetPopularProducts get products with most view
  - Get most viewed product in past 2 day (Each product must view as least 5 times)
  - If not enough limit, merge with list most view product
*/
func (d *domain) GetPopularProducts(ctx context.Context, limit int32) ([]*product.SummaryProductView, error) {
	d.logger.DebugCtx(ctx, "GetPopularProducts with limit:", limit)
	to := time.Now().UTC()
	from := to.AddDate(0, -1, 0).UTC() // default in 1 month, maybe change later
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
func (d *domain) GetRelatedProducts(ctx context.Context, productID int64, limit int32) ([]*product.SummaryProductView, error) {
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
	if len(viewHistory) == 0 {
		viewHistory, err = d.GetPopularProducts(ctx, limit)
		if err != nil {
			d.logger.ErrorCtx(ctx, err)
			return nil, err
		}
	}
	return viewHistory, nil
}

/*
GetCollaborativeRecommendation recommend products based on other users' view history who viewed the same products as current user
  - Get products user viewed recently
  - For each product, find other users who viewed it (excluding current user)
  - Aggregate products viewed by these other users
  - Sort aggregated products by count desc
*/
func (d *domain) GetCollaborativeRecommendation(ctx context.Context, userID string, limit int32) ([]*product.SummaryProductView, error) {
	d.logger.DebugCtx(ctx, "GetCollaborativeRecommendation with userID:", userID)
	if userID == "" {
		return nil, errors.New("userID is empty")
	}

	// Step 1: Get products user viewed recently
	userViewed, err := d.userViewHistory.RecentViewProductsByUser(ctx, userID, limit)
	if err != nil {
		d.logger.ErrorCtx(ctx, err)
		return nil, err
	}
	if len(userViewed) == 0 {
		// fallback popular products
		d.logger.DebugCtx(ctx, "recent view products is empty, returning popular recommendations")
		return d.GetPopularProducts(ctx, limit)
	}

	// Step 2: For each product, find other users who viewed it (excluding current user)
	otherUsers := make(map[string]struct{})
	for _, prod := range userViewed {
		users, err := d.userViewHistory.UsersWhoViewedProduct(ctx, prod.ProductID, userID)
		if err != nil {
			d.logger.ErrorCtx(ctx, err)
			continue
		}
		for _, u := range users {
			otherUsers[u] = struct{}{}
		}
	}

	if len(otherUsers) == 0 {
		// no similar users, fallback popular
		d.logger.DebugCtx(ctx, "no other users found, returning popular recommendations")
		return d.GetPopularProducts(ctx, limit)
	}

	// Step 3: Aggregate products viewed by these other users
	aggProductViews := make(map[int64]int64) // productID -> count
	for u := range otherUsers {
		views, err := d.userViewHistory.RecentViewProductsByUser(ctx, u, 20)
		if err != nil {
			d.logger.ErrorCtx(ctx, err)
			continue
		}
		for _, v := range views {
			aggProductViews[v.ProductID] += 1
		}
	}

	resp := make([]*product.SummaryProductView, 0)
	for prodID, count := range aggProductViews {
		resp = append(resp, &product.SummaryProductView{
			ProductID: prodID,
			ViewCount: helper.AnyToPointer(count),
		})
	}
	if len(resp) == 0 {
		// fallback popular products
		return d.GetPopularProducts(ctx, limit)
	}

	// Step 5: Sort aggregated products by count desc
	resp = helper.Sort(resp, func(i, j int) bool {
		return *resp[i].ViewCount > *resp[j].ViewCount // decreasing
	})
	return resp, nil
}
