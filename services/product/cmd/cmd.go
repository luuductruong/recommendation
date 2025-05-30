package cmd

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/recommendation/services/core/application/middleware"
	appService "github.com/recommendation/services/core/application/product/service"
	"github.com/recommendation/services/core/infra/config"
	"github.com/recommendation/services/core/infra/db"
	handler "github.com/recommendation/services/product/application/grpchandler"
	"github.com/recommendation/services/product/domain"
	repo "github.com/recommendation/services/product/external/repository"
)

var (
	sql         db.SQL
	appConfig   *Config
	grpcHandler appService.ProductServiceServer
)

func Run() {
	var err error
	//appConfig = &Config{}
	err = config.LoadConfig(&appConfig)
	if err != nil {
		log.Fatal(err)
	}
	sql, err = db.NewSQL(appConfig.DataBase)
	if err != nil {
		fmt.Println("ERRRRRRRRRRRRR ", err)
		//log.Fatal(err)
	}

	productDomain := domain.NewDomain(&domain.ProductDomainParam{
		ProductRepo: repo.NewProductRepo(),
	})
	grpcHandler = handler.NewHandler(productDomain)

	grpcServe()
}

func grpcServe() {
	// Start gRPC server in goroutine
	lis, err := net.Listen("tcp", ":8003")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("grpc server listening on :", lis.Addr().String())
	serve := middleware.NewServer(grpc.UnaryInterceptor(middleware.GrpcChainUnaryServer(sql)))
	defer serve.GracefulStop()

	appService.RegisterProductServiceServer(serve, grpcHandler)
	log.Fatal(serve.Serve(lis))
}
