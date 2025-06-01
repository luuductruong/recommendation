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
	NotByProductIDs(productIDs ...int64) ProductQuery
	Limit(limit int) ProductQuery
	// result
	Result() (*Product, error)
	ResultList() ([]*Product, error)
}

type UserViewHistoryRepo interface {
	// input
	Upsert(ctx context.Context, product *UserViewHistory) error
	// query
	Query(ctx context.Context) UserViewHistoryQuery
	// get product_ids have the most viewed in time range. Return product_id and view_count
	MostViewedInTimeRange(ctx context.Context, viewFrom, viewTo time.Time, limit int32) ([]*SummaryProductView, error)
	// get product_ids have the most viewed. Return product_id and view_count
	MostView(ctx context.Context, limit int32) ([]*SummaryProductView, error)
	// get product_ids viewed by user. Return product_id and view_time
	RecentViewProductsByUser(ctx context.Context, userID string, limit int32) ([]*SummaryProductView, error)
	// get similar category with product_id. Return product_id and view_count
	GetMostViewedProductsInCategory(ctx context.Context, categoryID string, excludeProductID int64, limit int32) ([]*SummaryProductView, error)
	// get user also viewed product, exclude userID
	UsersWhoViewedProduct(ctx context.Context, productID int64, excludeUserID string) ([]string, error)
}

type UserViewHistoryQuery interface {
	// query
	ByProductID(productID int64) UserViewHistoryQuery
	ByUserID(userID string) UserViewHistoryQuery
	ViewedAfter(viewTime time.Time) UserViewHistoryQuery
	DistinctByProductID() UserViewHistoryQuery
	// ordering
	OrderByViewedTime(desc bool) UserViewHistoryQuery
	//
	Limit(limit int) UserViewHistoryQuery
	// result
	Result() (*UserViewHistory, error)
	ResultList() ([]*UserViewHistory, error)
}

type CategoryViewHistoryRepo interface {
	IncreaseViewCount(ctx context.Context, categoryID string) error
	Query(ctx context.Context) CategoryViewHistoryQuery
	Upsert(ctx context.Context, cateHis *CategoryViewHistory) error
}

type CategoryViewHistoryQuery interface {
	ByCategoryID(categoryID string) CategoryViewHistoryQuery
	Result() (*CategoryViewHistory, error)
}
