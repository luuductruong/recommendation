package dto

import (
	"github.com/recommendation/services/core/application/product/model"
	"github.com/recommendation/services/core/domain/product"
)

func MapProductFromDm(source *product.Product) *model.Product {
	if source == nil {
		return nil
	}
	return &model.Product{
		Id:    source.ProductID,
		Name:  source.Name,
		Price: source.Price,
	}
}
