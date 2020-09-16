package model

type Artist struct {
	Name            string     `csv:"name"`
	OriginalName    string     `csv:"original_name"`
	Nationalities   StringList `csv:"nationalities"`
	PaintingSchools StringList `csv:"painting_schools"`
	ArtMovements    StringList `csv:"art_movements"`
	BirthDate       Date       `csv:"birth_date"`
	DeathDate       Date       `csv:"death_date"`
}
