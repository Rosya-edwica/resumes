package position

import (
	"resumes/internal/entities"
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

func (r *Reposititory) GetPositions() ([]models.Position, error) {
	query := `SELECT id, name, other_names FROM position WHERE id NOT IN (select distinct position_id from resume) ORDER BY id ASC`
	positions := make([]entities.Position, 0)
	err := r.db.Select(&positions, query)
	if err != nil {
		return nil, errors.Wrap(err, "select position")
	}

	return models.NewPositions(positions), nil
}
