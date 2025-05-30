package domain

import (
	"errors"

	"github.com/recommendation/services/core/context"
	"github.com/recommendation/services/core/domain/product"
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
	go d.updateUserHistory(ctx, inp.UserID, prod)
	return prod, nil
}

func (d *domain) updateUserHistory(ctx context.Context, userID string, prod *product.Product) {

}
