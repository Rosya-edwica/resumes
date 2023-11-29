package main

import (
	"fmt"
	"log"
	config "resumes/configs"
	"resumes/internal/database"
	"resumes/internal/database/city"
	"resumes/internal/database/position"
	"resumes/internal/platforms/headhunter"
)

func main() {
	db, err := database.New(config.GetDBConfig())
	if err != nil {
		log.Fatal(err)
	}

	cities, err := city.NewRep(db).GetCities()
	if err != nil || len(cities) == 0 {
		log.Fatal(fmt.Sprintf("Error: %s\tLen(cities)=%d", err, len(cities)))
	}
	positions, err := position.NewRep(db).GetPositions()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: %s\tLen(positions)=%d", err, len(positions)))
	}

	hh := headhunter.HeadHunter{
		Positions: positions,
		Cities:    cities,
		DB:        db,
	}
	hh.ParseAll()
	db.Close()
}
