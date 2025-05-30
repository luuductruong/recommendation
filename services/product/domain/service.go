package domain

import (
	"github.com/recommendation/services/core/domain/product"
)

type ProductDomainParam struct {
	ProductRepo product.ProductRepo
}

type domain struct {
	productRepo product.ProductRepo
}

func NewDomain(param *ProductDomainParam) product.Service {
	return &domain{
		productRepo: param.ProductRepo,
	}
}
