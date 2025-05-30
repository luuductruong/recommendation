package product

import (
	"github.com/recommendation/services/core/context"
	"time"
)

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
	MostPopularProductsInTimeRange(ctx context.Context, viewFrom, viewTo time.Time, limit int32) ([]*SummaryProductView, error)
	RecentViewProductsByUser(ctx context.Context, userID string, limit int32) ([]*SummaryProductView, error)
}

type UserViewHistoryQuery interface {
	// query
	ByProductID(productID int64) UserViewHistoryQuery
	ByUserID(userID string) UserViewHistoryQuery
	ViewedAfter(viewTime time.Time) UserViewHistoryQuery
	DistinctByProductID() UserViewHistoryQuery
	// ordering
	OrderByViewedTime(desc bool) UserViewHistoryQuery
	// result
	Result() (*UserViewHistory, error)
	ResultList() ([]*UserViewHistory, error)
}
