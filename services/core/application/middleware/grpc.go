package middleware

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/recommendation/services/core/infra/db"
	"github.com/recommendation/services/core/middleware"
)

func GrpcChainUnaryServer(sql db.SQL) grpc.UnaryServerInterceptor {
	middleWareChain := []grpc.UnaryServerInterceptor{
		middleware.GrpcDatabaseTx(sql),
	}
	return grpc_middleware.ChainUnaryServer(middleWareChain...)
}

func NewServer(option grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(option, grpc.StatsHandler(otelgrpc.NewServerHandler()))
}
