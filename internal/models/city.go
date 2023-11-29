package models

type City struct {
	HeadhunterID int    `db:"id_hh"`
	SuperjobID   int    `db:"id_superjob"`
	EdwicaID     int    `db:"id_edwica"`
	Name         string `db:"name"`
}
