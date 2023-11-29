package models

type Language struct {
	Name  string
	Level string
}

type Education struct {
	ResumeId  string
	Title     string
	Direction string
	Year      string
	Type      string
}

type Experience struct {
	ResumeId         string
	Post             string
	Period           string
	DurationInMonths int
	Branch           string
	Subbranch        string
}

type Resume struct {
	Id               string
	Platform         string
	Title            string
	Category         string
	City             string
	CityId           int
	PositionId       int
	Url              string
	Salary           int
	Currency         string
	GlobalExperience string
	SkillsList       []string
	SpecsList        []string
	LanguagesList    []Language
	ExperienceList   []Experience
	EducationList    []Education
}
