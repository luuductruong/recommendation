package repository

import (
	"github.com/recommendation/services/core/context"
	productDm "github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/helper/sql/query"
	"time"
)

type categoryViewHistory struct {
	ID         string
	CategoryID string
	TotalView  int64
	LastViewAt time.Time
}

func mapCategoryViewHistoryToDm(source *categoryViewHistory) *productDm.CategoryViewHistory {
	if source == nil {
		return nil
	}
	return &productDm.CategoryViewHistory{
		ID:         source.ID,
		CategoryID: source.CategoryID,
		TotalView:  source.TotalView,
		LastViewAt: source.LastViewAt,
	}
}

func mapCategoryViewHistoryFromDm(source *productDm.CategoryViewHistory) *categoryViewHistory {
	if source == nil {
		return nil
	}
	return &categoryViewHistory{
		ID:         source.ID,
		CategoryID: source.CategoryID,
		TotalView:  source.TotalView,
		LastViewAt: source.LastViewAt,
	}
}

func (c categoryViewHistory) TableName() string {
	return "category_view_history"
}

func NewCategoryViewHistoryRepo() productDm.CategoryViewHistoryRepo {
	return &categoryViewHistoryRepo{}
}

type categoryViewHistoryRepo struct {
}

func (c *categoryViewHistoryRepo) IncreaseViewCount(ctx context.Context, categoryID string) error {
	sql := `UPDATE category_view_history
			SET total_view = total_view + 1,
			    last_view_at = NOW()
			WHERE category_id = ?`
	return ctx.GetDbTx().Exec(sql, categoryID).Error
}

func (c *categoryViewHistoryRepo) Upsert(ctx context.Context, cateHis *productDm.CategoryViewHistory) error {
	return query.Upsert(ctx.GetDbTx(), cateHis, mapCategoryViewHistoryFromDm)
}

type categoryViewHistoryQuery struct {
	query.BaseQuery
}

func (c *categoryViewHistoryQuery) ByCategoryID(categoryID string) productDm.CategoryViewHistoryQuery {
	return query.Where(c, "category_id = ?", categoryID)
}

func (c *categoryViewHistoryQuery) Result() (*productDm.CategoryViewHistory, error) {
	return query.Result(c, mapCategoryViewHistoryToDm)
}

func (c *categoryViewHistoryRepo) Query(ctx context.Context) productDm.CategoryViewHistoryQuery {
	return &categoryViewHistoryQuery{query.NewBQ(ctx.GetDbTx().Model(&categoryViewHistory{}))}
}
