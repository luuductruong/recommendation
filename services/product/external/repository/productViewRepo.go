package repository

import (
	"github.com/recommendation/services/core/context"
	productDm "github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/helper/sql/query"
	"time"
)

type userViewHistory struct {
	ID        string
	UserID    string
	ProductID int64
	ViewAt    time.Time
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
	return &userViewHistoryRepo{}
}

type userViewHistoryRepo struct {
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

func (u *userViewHistoryQuery) Result() (*productDm.UserViewHistory, error) {
	return query.Result(u, mapUserViewHistoryToDm)
}

func (u *userViewHistoryQuery) ResultList() ([]*productDm.UserViewHistory, error) {
	return query.ResultList(u, mapUserViewHistoryToDm)
}
