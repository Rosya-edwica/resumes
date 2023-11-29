package headhunter

import (
	"resumes/internal/models"

	"github.com/gocolly/colly"
)

var (
	universityTag = "div[data-qa='resume-block-education']"             // Блок Высшее образование
	examTag       = "div[data-qa='resume-block-attestation-education']" // Блок Экзамены, тесты
	coursesTag    = "div[data-qa='resume-block-additional-education']"  // Блок Повышение квалификации
	educationTags = []string{universityTag, examTag, coursesTag}
)

// В одном резюме может быть несколько блоков, связанных с образованием
func getEducation(html *colly.HTMLElement) (education []models.Education) {
	for _, item := range educationTags {
		ed := getEducationByTag(html, item)
		education = append(education, ed...)
	}
	return
}

func getEducationByTag(html *colly.HTMLElement, tag string) (education []models.Education) {
	edType := html.ChildText(tag + " span.resume-block__title-text.resume-block__title-text_sub") // Заголовок блока
	html.ForEach(tag+" div.resume-block-item-gap", func(i int, h *colly.HTMLElement) {
		if i > 0 {
			var edu models.Education
			edu.Type = edType
			edu.Title = h.ChildText("div[data-qa='resume-block-education-name']")
			edu.Direction = h.ChildText("div[data-qa='resume-block-education-organization']")
			edu.Year = h.ChildText("div.bloko-column.bloko-column_xs-4.bloko-column_s-2.bloko-column_m-2.bloko-column_l-2")
			education = append(education, edu)
		}
	})
	return
}
