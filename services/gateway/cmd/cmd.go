package cmd

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"

	coreMdw "github.com/recommendation/services/core/application/middleware"
	productGw "github.com/recommendation/services/core/application/product/service"
	cfn "github.com/recommendation/services/core/infra/config"
	grpcHelper "github.com/recommendation/services/core/infra/grpc"
	gatewayMdw "github.com/recommendation/services/gateway/middleware"
)

var (
	appConfig *Config

	// client connect
	productConn *grpc.ClientConn
)

func Run() {
	var err error
	// HARD code app config. Can read from file .yml or env in deployment
	appConfig = &Config{
		Product: &cfn.Client{
			Host: "localhost",
			Port: "8003", // connect to product
			SSL:  false,
		},
		Http: &cfn.Client{
			Host: "localhost",
			Port: "8080", // expose api
			SSL:  false,
		},
	}

	// client connect
	productConn, err = grpcHelper.NewClient(appConfig.Product)
	if err != nil {
		log.Fatal(err)
	}

	// make context with cancel
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// setup timeout for request
	runtime.DefaultContextTimeout = 30 * time.Second // for example

	gMux := runtime.NewServeMux()

	// register client connect
	err = productGw.RegisterProductServiceHandler(ctx, gMux, productConn)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", gatewayMdw.WithCORS(gMux))
	combineMd := gatewayMdw.ChainCombine(coreMdw.WithResponseType())
	server := http.Server{
		Handler: combineMd(mux),
	}
	listen, err := net.Listen("tcp", appConfig.Http.Address())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Gateway start listening on " + appConfig.Http.Address())
	log.Fatal(server.Serve(listen))
}
