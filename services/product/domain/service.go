package domain

import (
	"github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/infra/logger"
)

type ProductDomainParam struct {
	ProductRepo         product.ProductRepo
	UserViewHistoryRepo product.UserViewHistoryRepo
}

type domain struct {
	logger          logger.Logger
	productRepo     product.ProductRepo
	userViewHistory product.UserViewHistoryRepo
}

func NewDomain(param *ProductDomainParam) product.Service {
	return &domain{
		logger:          logger.Default,
		productRepo:     param.ProductRepo,
		userViewHistory: param.UserViewHistoryRepo,
	}
}
