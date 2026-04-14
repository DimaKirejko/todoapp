package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("no rows") // need replays all pgx.ErrNoRows
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrUnknown            = errors.New("unknown") //13;27
)
