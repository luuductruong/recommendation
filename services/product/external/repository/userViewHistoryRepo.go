package repository

import (
	"fmt"
	"github.com/recommendation/services/core/context"
	productDm "github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/helper"
	"github.com/recommendation/services/core/helper/sql/query"
	"time"
)

type userViewHistory struct {
	ID        string
	UserID    string
	ProductID int64
	ViewAt    time.Time
}

type summaryProductView struct {
	ProductID int64
	ViewCount *int64
	ViewAt    *time.Time
}

func mapSummaryProductViewToDm(source *summaryProductView) *productDm.SummaryProductView {
	if source == nil {
		return nil
	}
	return &productDm.SummaryProductView{
		ProductID: source.ProductID,
		ViewCount: source.ViewCount,
		ViewAt:    source.ViewAt,
	}
}

func mapUserViewHistoryToDm(source *userViewHistory) *productDm.UserViewHistory {
	if source == nil {
		return nil
	}
	return &productDm.UserViewHistory{
		ID:        source.ID,
		UserID:    source.UserID,
		ProductID: source.ProductID,
		ViewAt:    source.ViewAt,
	}
}
func mapUserViewHistoryFromDm(source *productDm.UserViewHistory) *userViewHistory {
	if source == nil {
		return nil
	}
	return &userViewHistory{
		ID:        source.ID,
		UserID:    source.UserID,
		ProductID: source.ProductID,
		ViewAt:    source.ViewAt,
	}
}

func (u userViewHistory) TableName() string {
	return "user_view_history"
}

func NewUserViewHistoryRepo() productDm.UserViewHistoryRepo {
	return &userViewHistoryRepo{
		TableName: userViewHistory{}.TableName(),
	}
}

type userViewHistoryRepo struct {
	TableName string
}

func (u *userViewHistoryRepo) MostViewedInTimeRange(ctx context.Context, viewFrom, viewTo time.Time, limit int32) ([]*productDm.SummaryProductView, error) {
	var result []*summaryProductView
	sql := fmt.Sprintf(`SELECT *
			FROM (
    			SELECT product_id,
			    COUNT(*) AS view_count
				FROM %s
			    WHERE view_at >= $1 AND view_at <= $2
			    GROUP BY product_id
			    ORDER BY view_count DESC
			) AS pivc
			LIMIT $3`, u.TableName)

	q := ctx.GetDbTx().Raw(sql, viewFrom, viewTo, limit)
	err := q.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return helper.MapList(result, mapSummaryProductViewToDm), nil
}

func (u *userViewHistoryRepo) RecentViewProductsByUser(ctx context.Context, userID string, limit int32) ([]*productDm.SummaryProductView, error) {
	var result []*summaryProductView
	sql := fmt.Sprintf(`SELECT product_id, view_at
FROM (
  SELECT product_id, view_at,
         ROW_NUMBER() OVER (PARTITION BY product_id ORDER BY view_at DESC) as rn
  FROM %s
  WHERE user_id = $1
) AS ranked
WHERE rn = 1
ORDER BY view_at DESC
LIMIT $2`, u.TableName)

	q := ctx.GetDbTx().Raw(sql, userID, limit)
	err := q.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return helper.MapList(result, mapSummaryProductViewToDm), nil
}

func (u *userViewHistoryRepo) MostView(ctx context.Context, limit int32) ([]*productDm.SummaryProductView, error) {
	var result []*summaryProductView
	sql := fmt.Sprintf(`SELECT DISTINCT product_id, count(*) as view_count
								FROM %s
								GROUP BY product_id
								ORDER BY view_count desc
								LIMIT $1;`, u.TableName)
	q := ctx.GetDbTx().Raw(sql, limit)
	err := q.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return helper.MapList(result, mapSummaryProductViewToDm), nil
}

func (u *userViewHistoryRepo) GetMostViewedProductsInCategory(ctx context.Context, categoryID string, excludeProductID int64, limit int32) ([]*productDm.SummaryProductView, error) {
	var result []*summaryProductView

	sql := fmt.Sprintf(`
		SELECT product_id, COUNT(*) AS view_count
		FROM %s uvh
		JOIN products p ON uvh.product_id = p.id
		WHERE p.category = $1 AND uvh.product_id != $2
		GROUP BY uvh.product_id
		ORDER BY view_count DESC
		LIMIT $3;`, u.TableName)

	q := ctx.GetDbTx().Raw(sql, categoryID, excludeProductID, limit)
	err := q.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return helper.MapList(result, mapSummaryProductViewToDm), nil
}

func (u *userViewHistoryRepo) UsersWhoViewedProduct(ctx context.Context, productID int64, excludeUserID string) ([]string, error) {
	var userIDs []string
	sql := fmt.Sprintf(`SELECT DISTINCT user_id FROM %s WHERE product_id = $1 AND user_id != $2`, u.TableName)
	q := ctx.GetDbTx().Raw(sql, productID, excludeUserID)
	err := q.Find(&userIDs).Error
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (u *userViewHistoryRepo) Upsert(ctx context.Context, view *productDm.UserViewHistory) error {
	return query.Upsert(ctx.GetDbTx(), view, mapUserViewHistoryFromDm)
}

type userViewHistoryQuery struct {
	query.BaseQuery
}

func (u *userViewHistoryRepo) Query(ctx context.Context) productDm.UserViewHistoryQuery {
	return &userViewHistoryQuery{query.NewBQ(ctx.GetDbTx().Model(&userViewHistory{}))}
}

func (u *userViewHistoryQuery) ByProductID(productID int64) productDm.UserViewHistoryQuery {
	return query.Where(u, "product_id = ?", productID)
}

func (u *userViewHistoryQuery) ByUserID(userID string) productDm.UserViewHistoryQuery {
	return query.Where(u, "user_id = ?", userID)
}

func (u *userViewHistoryQuery) ViewedAfter(viewTime time.Time) productDm.UserViewHistoryQuery {
	return query.Where(u, "view_at > ?", viewTime)
}

func (u *userViewHistoryQuery) DistinctByProductID() productDm.UserViewHistoryQuery {
	u.SetDB(u.GetDB().Distinct("product_id"))
	return u
}

func (u *userViewHistoryQuery) Result() (*productDm.UserViewHistory, error) {
	return query.Result(u, mapUserViewHistoryToDm)
}

func (u *userViewHistoryQuery) ResultList() ([]*productDm.UserViewHistory, error) {
	return query.ResultList(u, mapUserViewHistoryToDm)
}

// ordering
func (u *userViewHistoryQuery) OrderByViewedTime(desc bool) productDm.UserViewHistoryQuery {
	return query.OrderBy(u, "view_at", desc)
}

// ordering
func (u *userViewHistoryQuery) Limit(limit int) productDm.UserViewHistoryQuery {
	return query.Limit(u, limit)
}
