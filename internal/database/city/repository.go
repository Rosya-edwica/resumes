package city

import (
	"resumes/internal/models"

	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
)

type Reposititory struct {
	db *sqlx.DB
}

func NewRep(db *sqlx.DB) *Reposititory {
	return &Reposititory{db: db}
}

func (r *Reposititory) GetCities() ([]models.City, error) {
	query := `SELECT id_edwica, id_superjob, id_hh, name FROM h_city ORDER BY id_hh ASC`
	cities := make([]models.City, 0)
	err := r.db.Select(&cities, query)
	if err != nil {
		return nil, errors.Wrap(err, "select city")
	}

	return cities, nil
}
