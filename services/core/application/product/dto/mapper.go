package dto

import (
	"github.com/recommendation/services/core/application/product/model"
	"github.com/recommendation/services/core/domain/product"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapProductFromDm(source *product.Product) *model.Product {
	if source == nil {
		return nil
	}
	return &model.Product{
		Id:         source.ProductID,
		Name:       source.Name,
		Price:      source.Price,
		CategoryId: source.CategoryID,
	}
}

func MapViewStatusFromDm(source *product.SummaryProductView) *model.ViewStatus {
	if source == nil {
		return nil
	}
	totalView := int64(0)
	if source.ViewCount != nil {
		totalView = *source.ViewCount
	}
	var ts *timestamppb.Timestamp
	if source.ViewAt != nil {
		ts = timestamppb.New(*source.ViewAt)
	}
	return &model.ViewStatus{
		ProductId: source.ProductID,
		TotalView: totalView,
		LastView:  ts,
	}
}
