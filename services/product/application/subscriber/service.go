package subscriber

import (
	"context"
	appContext "github.com/recommendation/services/core/context"
	"github.com/recommendation/services/core/domain/product"
	"github.com/recommendation/services/core/infra/pubsub"
)

type Handler interface {
	pubsub.BaseSubsHandler
}

func NewHandler(domain product.Service) Handler {
	return &handler{
		domain: domain,
	}
}

type handler struct {
	domain product.Service
}

func (h *handler) RouteSetup() map[pubsub.MessageName]pubsub.SubscriptionHandler {
	return map[pubsub.MessageName]pubsub.SubscriptionHandler{}
}

func (h *handler) Test(ctx context.Context, msg *pubsub.Message) error {
	appCtx := appContext.FromContext(ctx)
	type Test struct {
		ProductID int64
	}
	var test Test
	err := msg.ScanPayload(&test)
	if err != nil {
		return err
	}
	if _, err = h.domain.GetProductDetail(appCtx, &product.GetProductDetailInp{
		UserID:    "123",
		ProductID: test.ProductID,
	}); err != nil {
		return err
	}
	return nil
}
