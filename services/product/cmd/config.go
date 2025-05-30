package cmd

import "github.com/recommendation/services/core/infra/db"

type Config struct {
	DataBase *db.Config `config:"database"`
}
