package context

import (
	"context"
	"github.com/recommendation/services/core/middleware"
	"gorm.io/gorm"
)

type ctxInternal struct {
	context.Context
}

func FromContext(ctx context.Context) Context {
	return &ctxInternal{ctx}
}

type Context interface {
	context.Context
	GetDbTx() *gorm.DB
}

func (c *ctxInternal) GetDbTx() *gorm.DB {
	return middleware.DbTxFromContext(c)
}
