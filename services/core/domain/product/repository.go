package product

import "github.com/recommendation/services/core/context"

type ProductRepo interface {
	Query(ctx context.Context) ProductQuery
}

type ProductQuery interface {
	// query
	ByProductID(productID int64) ProductQuery
	// result
	Result() (*Product, error)
	ResultList() ([]*Product, error)
}
