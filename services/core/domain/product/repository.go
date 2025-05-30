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

type ProductViewRepo interface {
	// input
	Upsert(ctx context.Context, product *ProductView) error
	// query
	Query(ctx context.Context) ProductViewQuery
}

type ProductViewQuery interface {
	// query
	ByProductID(productID int64) ProductViewQuery
	ByUserID(userID string) ProductViewQuery
	// result
	Result() (*ProductView, error)
	ResultList() ([]*ProductView, error)
}
