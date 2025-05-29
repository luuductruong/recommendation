package cmd

import (
	cfn "github.com/recommentation/service/core/infra/config"
)

type Config struct {
	Http *cfn.Client `config:"http"`
	// microservices
	Product *cfn.Client `config:"product"`
	Order   *cfn.Client `config:"order"`
}
