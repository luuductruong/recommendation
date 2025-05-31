package repository

import (
	"github.com/recommendation/services/core/context"
	productDm "github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/helper/sql/query"
)

type product struct {
	ProductID  int64 `gorm:"primary_key"`
	Name       string
	Price      float64
	CategoryID string
}

func mapProductToDm(source *product) *productDm.Product {
	if source == nil {
		return nil
	}
	return &productDm.Product{
		ProductID:  source.ProductID,
		Name:       source.Name,
		Price:      source.Price,
		CategoryID: source.CategoryID,
	}
}

func (p product) TableName() string {
	return "product"
}

func NewProductRepo() productDm.ProductRepo {
	return &productRepo{}
}

type productRepo struct {
}

type productQuery struct {
	query.BaseQuery
}

func (p *productRepo) Query(ctx context.Context) productDm.ProductQuery {
	return &productQuery{query.NewBQ(ctx.GetDbTx().Model(&product{}))}
}

func (p *productQuery) ByProductID(productID int64) productDm.ProductQuery {
	return query.Where(p, "product_id = ?", productID)
}

func (p *productQuery) Limit(limit int) productDm.ProductQuery {
	return query.Limit(p, limit)
}

func (p *productQuery) Result() (*productDm.Product, error) {
	return query.Result(p, mapProductToDm)
}

func (p *productQuery) ResultList() ([]*productDm.Product, error) {
	return query.ResultList(p, mapProductToDm)
}
