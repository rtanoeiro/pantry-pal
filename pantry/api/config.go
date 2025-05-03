package api

import "pantry-pal/pantry/database"

type Config struct {
	Db     *database.Queries
	Env    string
	Secret string
}
