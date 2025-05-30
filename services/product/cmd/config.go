package cmd

import (
	"github.com/recommendation/services/core/infra/config"
	"github.com/recommendation/services/core/infra/db"
	"github.com/recommendation/services/core/infra/logger"
)

type Config struct {
	Grpc     *config.Client `config:"grpc"`
	Logger   *logger.Config `config:"logger"`
	DataBase *db.Config     `config:"database"`
}
