package grpchandler

import (
	"context"

	"github.com/recommendation/services/core/application/product/dto"
	"github.com/recommendation/services/core/application/product/service"
	appContext "github.com/recommendation/services/core/context"
	"github.com/recommendation/services/core/domain/product"
)

type handler struct {
	service.UnimplementedProductServiceServer
	productDomain product.Service
}

func NewHandler(productDomain product.Service) service.ProductServiceServer {
	return &handler{
		productDomain: productDomain,
	}
}

func (h *handler) GetProductDetail(ctx context.Context, req *dto.GetProductDetailReq) (*dto.GetProductDetailResp, error) {
	appCtx := appContext.FromContext(ctx)
	prod, err := h.productDomain.GetProductDetail(appCtx, &product.GetProductDetailInp{
		UserID:    req.UserId,
		ProductID: req.ProductId,
	})
	if err != nil {
		return nil, err
	}
	return &dto.GetProductDetailResp{
		Product: dto.MapProductFromDm(prod),
	}, nil
}

func (h *handler) GetRecommendationForUser(ctx context.Context, req *dto.GetRecommendationForUserReq) (*dto.GetRecommendationForUserResp, error) {
	appCtx := appContext.FromContext(ctx)
	productIds, err := h.productDomain.GetRecommendationForUser(appCtx, &product.GetRecommendationForUserInp{
		UserID: req.UserId,
		Limit:  req.Limit,
	})
	if err != nil {
		return nil, err
	}
	return &dto.GetRecommendationForUserResp{
		ProductIds: productIds,
	}, nil
}
