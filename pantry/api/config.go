package api

import (
	"pantry-pal/pantry/database"
)

type Config struct {
	Db       *database.Queries
	Renderer *Templates
	Env      string
	Secret   string
}
