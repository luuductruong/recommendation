package domain

import (
	"errors"
	"time"

	"github.com/recommendation/services/core/context"
	"github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/helper"
)

const deboundTime = time.Second * 1 // for example

func (d *domain) GetProductDetail(ctx context.Context, inp *product.GetProductDetailInp) (*product.Product, error) {
	d.logger.DebugCtx(ctx, "GetProductDetail with inp:", inp)
	if inp == nil || inp.ProductID <= 0 {
		d.logger.DebugCtx(ctx, "GetProductDetail with nil input")
		return nil, nil
	}
	// query product
	prod, err := d.productRepo.Query(ctx).ByProductID(inp.ProductID).Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "Error querying product")
		return nil, err
	}
	if prod == nil {
		d.logger.DebugCtx(ctx, "GetProductDetail: product is nil")
		return nil, errors.New("product not found")
	}
	//put back job for record user history
	go d.recordViewHistory(ctx.Clone(), inp.UserID, prod)
	return prod, nil
}

func (d *domain) recordViewHistory(ctx context.Context, userID string, prod *product.Product) {
	if prod == nil || userID == "" { // No need to record for guess user. Can be removed if we want tracking un authorize request
		return
	}
	now := time.Now()
	d.logger.DebugCtx(ctx, "recordViewHistory with prod:", prod)
	// query by user_id, product_id, newly
	existView, err := d.userViewHistory.Query(ctx).
		ByUserID(userID).
		ByProductID(prod.ProductID).
		OrderByViewedTime(true).
		Limit(1).
		Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "Error querying product view")
		return
	}
	if existView != nil && existView.ViewAt.UTC().After(now.Add(-deboundTime).UTC()) {
		d.logger.DebugCtx(ctx, "have view in past ", deboundTime)
		return
	}
	newView := &product.UserViewHistory{
		ID:        helper.NewStringUUID(),
		ViewAt:    now,
		ProductID: prod.ProductID,
		UserID:    userID,
	}
	err = d.userViewHistory.Upsert(ctx, newView)
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "Error upserting product view")
	}
	// setup category view history
	existsCate, err := d.cateViewHistory.Query(ctx).ByCategoryID(prod.CategoryID).Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "Error querying cate view")
	}
	if existsCate != nil {
		// increase
		err = d.cateViewHistory.IncreaseViewCount(ctx, prod.CategoryID)
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "Error increasing view count")
		}
	} else {
		newCateHistory := &product.CategoryViewHistory{
			ID:         helper.NewStringUUID(),
			CategoryID: prod.CategoryID,
			TotalView:  1,
			LastViewAt: now,
		}
		err = d.cateViewHistory.Upsert(ctx, newCateHistory)
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "Error upserting cate view")
		}
	}
	return
}
