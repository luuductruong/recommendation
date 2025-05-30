package domain

import (
	"github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/infra/logger"
)

type ProductDomainParam struct {
	ProductRepo     product.ProductRepo
	ProductViewRepo product.ProductViewRepo
}

type domain struct {
	logger          logger.Logger
	productRepo     product.ProductRepo
	productViewRepo product.ProductViewRepo
}

func NewDomain(param *ProductDomainParam) product.Service {
	return &domain{
		logger:          logger.Default,
		productRepo:     param.ProductRepo,
		productViewRepo: param.ProductViewRepo,
	}
}
