package cmd

import (
	"fmt"
	"github.com/recommendation/services/core/infra/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
		logger.Default.Panic("Error loading config: ", err)
	}
	log := logger.NewLogger(appConfig.Logger)
	logger.SetDefault(log)
	sql, err = db.NewSQL(appConfig.DataBase)
	if err != nil {
		logger.Default.Panic("Error connecting to database: ", err)
	}

	productDomain := domain.NewDomain(&domain.ProductDomainParam{
		ProductRepo:         repo.NewProductRepo(),
		UserViewHistoryRepo: repo.NewUserViewHistoryRepo(),
	})
	grpcHandler = handler.NewHandler(productDomain)

	grpcServe()
}

func grpcServe() {
	// Start gRPC server in goroutine
	lis, err := net.Listen("tcp", appConfig.Grpc.Address())
	if err != nil {
		logger.Default.Panic("failed to listen: %v", err)
	}
	fmt.Println("grpc server listening on :", lis.Addr().String())
	serve := middleware.NewServer(grpc.UnaryInterceptor(middleware.GrpcChainUnaryServer(sql)))
	defer serve.GracefulStop()

	appService.RegisterProductServiceServer(serve, grpcHandler)
	reflection.Register(serve)
	logger.Default.Fatal(serve.Serve(lis))
}
