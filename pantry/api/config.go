package api

import (
	"pantry-pal/pantry/database"
)

type Config struct {
	Db       *database.Queries
	Port     string
	DBUrl    string
	Renderer Renderer
	Env      string
	Secret   string
}
