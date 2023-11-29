package resume

import (
	"fmt"
	"resumes/internal/models"
	"resumes/pkg/logger"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Reposititory struct {
	db *sqlx.DB
}

func NewRep(db *sqlx.DB) *Reposititory {
	return &Reposititory{db: db}
}

func (r *Reposititory) SaveResumes(resumes []models.Resume) error {
	tx, err := r.db.Beginx()
	if err != nil {
		logger.Log.Printf("ОШИБКА при создании транзакции:%s\n", err)
		return err
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				logger.Log.Printf("ОШИБКА при откате транзакции:%s\n", err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			logger.Log.Printf("ОШИБКА при подтверждении транзакции:%s\n", err)
			return
		}
	}()

	err = r.saveResumes(tx, resumes)
	if err != nil {
		logger.Log.Printf("ОШИБКА при сохранении резюме:%s\n", err)
		return err
	}
	return nil
}

func (r *Reposititory) saveResumes(tx *sqlx.Tx, resumes []models.Resume) error {
	if len(resumes) == 0 {
		return nil
	}
	var (
		educationList  []models.Education
		experienceList []models.Experience
		valuesQuery    = make([]string, 0, len(resumes))
		valuesArgs     = make([]interface{}, 0, len(resumes))
	)

	for _, resume := range resumes {
		var ed []models.Education
		var ex []models.Experience
		for _, i := range resume.EducationList {
			i.ResumeId = resume.Id
			ed = append(ed, i)
		}
		for _, i := range resume.ExperienceList {
			i.ResumeId = resume.Id
			ex = append(ex, i)
		}
		educationList = append(educationList, ed...)
		experienceList = append(experienceList, ex...)

		languages := prepareLanguagesToSave(resume.LanguagesList)
		valuesQuery = append(valuesQuery, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valuesArgs = append(valuesArgs, resume.Id, resume.PositionId, resume.GlobalExperience, resume.CityId, strings.Join(resume.SpecsList, "|"), resume.Salary, resume.Currency, resume.Title, resume.Platform, resume.City, languages, strings.Join(resume.SkillsList, "|"), resume.Url)
	}
	query := fmt.Sprintf(`INSERT IGNORE INTO resume(id, position_id, general_experience, city_id, specs, salary, currency, title, platform, city, languages, skills, url)	VALUES %s`, strings.Join(valuesQuery, ","))
	_, err := tx.Exec(query, valuesArgs...)
	if err != nil {
		return err
	}

	err = r.saveEducation(tx, educationList)
	if err != nil {
		logger.Log.Printf("ОШИБКА при сохранении опыта:%s\n", err)
	}
	err = r.saveExperince(tx, experienceList)
	if err != nil {
		logger.Log.Printf("ОШИБКА при сохранении образования:%s\n", err)
	}
	return nil
}

func (r *Reposititory) saveEducation(tx *sqlx.Tx, educationList []models.Education) error {
	var (
		valuesQuery = make([]string, 0, len(educationList))
		valuesArgs  = make([]interface{}, 0, len(educationList))
	)

	for _, ed := range educationList {
		valuesQuery = append(valuesQuery, "(?, ?, ?, ?, ?)")
		valuesArgs = append(valuesArgs, ed.ResumeId, ed.Title, ed.Direction, ed.Year, ed.Type)

	}
	query := fmt.Sprintf(`INSERT IGNORE INTO resume_education(resume_id, name, direction, year, type) VALUES %s`, strings.Join(valuesQuery, ","))
	_, err := tx.Exec(query, valuesArgs...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Reposititory) saveExperince(tx *sqlx.Tx, experienceList []models.Experience) error {
	var (
		valuesQuery = make([]string, 0, len(experienceList))
		valuesArgs  = make([]interface{}, 0, len(experienceList))
	)

	for _, exp := range experienceList {
		valuesQuery = append(valuesQuery, "(?, ?, ?, ?, ?, ?)")
		valuesArgs = append(valuesArgs, exp.ResumeId, exp.Post, exp.Period, exp.DurationInMonths, exp.Branch, exp.Subbranch)

	}
	query := fmt.Sprintf(`INSERT IGNORE INTO resume_experience(resume_id, post, period, duration_in_months, branch, subbranch) VALUES %s`, strings.Join(valuesQuery, ","))
	_, err := tx.Exec(query, valuesArgs...)
	if err != nil {
		return err
	}
	return nil
}

func prepareLanguagesToSave(languages []models.Language) string {
	var langs []string

	for _, item := range languages {
		lng := fmt.Sprintf("%s: %s", item.Name, item.Level)
		langs = append(langs, lng)
	}

	return strings.Join(langs, "|")
}
