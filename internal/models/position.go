package models

import (
	"resumes/internal/entities"
	"strings"
)

type Position struct {
	Id         int
	Name       string
	OtherNames []string
}

func NewPositions(items []entities.Position) (positions []Position) {
	for _, i := range items {
		positions = append(positions, Position{
			Id:         i.Id,
			Name:       i.Name,
			OtherNames: strings.Split(i.OtherNames.String, "|"),
		})
	}
	return
}
