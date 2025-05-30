package domain

import (
	"errors"
	"fmt"

	"github.com/recommendation/services/core/context"
	"github.com/recommendation/services/core/domain/product"
)

func (d *domain) GetProductDetail(ctx context.Context, inp *product.GetProductDetailInp) (*product.Product, error) {
	fmt.Println("GetProductDetail", inp)
	if inp == nil || inp.ProductID <= 0 {
		fmt.Println("GetProductDetail: product is nil")
		return nil, nil
	}
	// query product
	prod, err := d.productRepo.Query(ctx).ByProductID(inp.ProductID).Result()
	if err != nil {
		fmt.Println("GetProductDetail:", err)
		return nil, err
	}
	if prod == nil {
		fmt.Println("GetProductDetail: product is nil")
		return nil, errors.New("product not found")
	}
	//put back job for record user history
	go d.updateUserHistory(ctx, inp.UserID, prod)
	return prod, nil
}

func (d *domain) updateUserHistory(ctx context.Context, userID string, prod *product.Product) {

}
