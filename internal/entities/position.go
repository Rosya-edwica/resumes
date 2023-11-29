package entities

import (
	"database/sql"
)

type Position struct {
	Id         int            `db:"id"`
	Name       string         `db:"name"`
	OtherNames sql.NullString `db:"other_names"`
}
