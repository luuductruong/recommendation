package domain

import (
	"errors"
	"time"

	"github.com/recommendation/services/core/context"
	"github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/helper"
)

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
	go d.updateUserViewHistory(ctx.Clone(), inp.UserID, prod)
	return prod, nil
}

func (d *domain) updateUserViewHistory(ctx context.Context, userID string, prod *product.Product) {
	if prod == nil {
		return
	}
	d.logger.DebugCtx(ctx, "updateUserViewHistory with prod:", prod)
	// query by user_id, product_id
	existView, err := d.userViewHistory.Query(ctx).
		ByUserID(userID).
		ByProductID(prod.ProductID).
		Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "Error querying product view")
		return
	}
	if existView != nil {
		d.logger.DebugCtx(ctx, "updateUserViewHistory with exist view:", existView)
		//update
		existView.ViewAt = time.Now()
		err = d.userViewHistory.Upsert(ctx, existView)
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "Error upserting product view")
		}
		return
	}
	d.logger.DebugCtx(ctx, "updateUserViewHistory with nil view")
	newView := &product.UserViewHistory{
		ID:        helper.NewStringUUID(),
		ViewAt:    time.Now(),
		ProductID: prod.ProductID,
		UserID:    userID,
	}
	err = d.userViewHistory.Upsert(ctx, newView)
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "Error upserting product view")
	}
	return
}
