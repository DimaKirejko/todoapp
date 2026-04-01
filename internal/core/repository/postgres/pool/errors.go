package core_postgres_pool

import "errors"

var (
	ErrNoRows = errors.New("no rows") // need replays all pgx.ErrNoRows
)
