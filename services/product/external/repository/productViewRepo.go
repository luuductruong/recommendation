package repository

import (
	"github.com/recommendation/services/core/context"
	productDm "github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/helper/sql/query"
	"time"
)

type productView struct {
	ID        string
	UserID    string
	ProductID int64
	ViewAt    time.Time
}

func mapProductViewToDm(source *productView) *productDm.ProductView {
	if source == nil {
		return nil
	}
	return &productDm.ProductView{
		ID:        source.ID,
		UserID:    source.UserID,
		ProductID: source.ProductID,
		ViewAt:    source.ViewAt,
	}
}
func mapProductViewFromDm(source *productDm.ProductView) *productView {
	if source == nil {
		return nil
	}
	return &productView{
		ID:        source.ID,
		UserID:    source.UserID,
		ProductID: source.ProductID,
		ViewAt:    source.ViewAt,
	}
}

func (v productView) TableName() string {
	return "product_view"
}

func NewProductViewRepo() productDm.ProductViewRepo {
	return &productViewRepo{}
}

type productViewRepo struct {
}

func (v *productViewRepo) Upsert(ctx context.Context, view *productDm.ProductView) error {
	return query.Upsert(ctx.GetDbTx(), view, mapProductViewFromDm)
}

type productViewQuery struct {
	query.BaseQuery
}

func (v *productViewRepo) Query(ctx context.Context) productDm.ProductViewQuery {
	return &productViewQuery{query.NewBQ(ctx.GetDbTx().Model(&productView{}))}
}

func (v *productViewQuery) ByProductID(productID int64) productDm.ProductViewQuery {
	return query.Where(v, "product_id = ?", productID)
}

func (v *productViewQuery) ByUserID(userID string) productDm.ProductViewQuery {
	return query.Where(v, "user_id = ?", userID)
}

func (v *productViewQuery) Result() (*productDm.ProductView, error) {
	return query.Result(v, mapProductViewToDm)
}

func (v *productViewQuery) ResultList() ([]*productDm.ProductView, error) {
	return query.ResultList(v, mapProductViewToDm)
}
