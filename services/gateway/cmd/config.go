package cmd

import (
	cfn "github.com/recommendation/services/core/infra/config"
)

type Config struct {
	Http *cfn.Client `config:"http"`
	// microservices
	Product *cfn.Client `config:"product"`
	Order   *cfn.Client `config:"order"`
}
