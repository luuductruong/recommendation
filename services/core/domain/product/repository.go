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

type UserViewHistoryRepo interface {
	// input
	Upsert(ctx context.Context, product *UserViewHistory) error
	// query
	Query(ctx context.Context) UserViewHistoryQuery
}

type UserViewHistoryQuery interface {
	// query
	ByProductID(productID int64) UserViewHistoryQuery
	ByUserID(userID string) UserViewHistoryQuery
	// result
	Result() (*UserViewHistory, error)
	ResultList() ([]*UserViewHistory, error)
}
