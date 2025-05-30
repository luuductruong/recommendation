package cmd

import (
	cfn "github.com/recommendation/services/core/infra/config"
	"github.com/recommendation/services/core/infra/logger"
)

type Config struct {
	Http   *cfn.Client    `config:"http"`
	Logger *logger.Config `config:"logger"`
	// microservices
	Product *cfn.Client `config:"product"`
}
