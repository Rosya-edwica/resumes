package headhunter

import (
	"regexp"
	"resumes/internal/models"
	"strings"

	"github.com/gocolly/colly"
)

func getExperience(html *colly.HTMLElement) (experience []models.Experience) {
	html.ForEach("div[data-qa='resume-block-experience'] div.resume-block-item-gap", func(i int, h *colly.HTMLElement) {
		if i > 0 {
			var exp models.Experience
			exp.Post = h.ChildText("div[data-qa='resume-block-experience-position']")
			exp.Period = getExperiencePeriod(h)
			exp.DurationInMonths = experienceToMonths(getExperienceDuration(h))
			experience = append(experience, exp)
		}
	})
	return
}

// Здесь используются регулярки для того, чтобы успешно отделить нужные данные от мусора
// Первая регулярка: `\W+ \d{4} — \W+ \d{4}` - ищет такие строки "Июнь 2014 — Август 2014"
// Вторая регулярка: `\W+ \d{4} —.*?\d` - ищет такие строки "Июнь 2014 — по настояющее время"
func getExperiencePeriod(html *colly.HTMLElement) string {
	periodText := html.ChildText("div.bloko-column.bloko-column_xs-4.bloko-column_s-2.bloko-column_m-2.bloko-column_l-2")
	periodText = strings.ReplaceAll(periodText, "\u00a0", " ")
	re := regexp.MustCompile(`\W+ \d{4} — \W+ \d{4}|\w+ \d{4} — \w+ \d{4}|\w+ \d{2} — \w+ \d{4}|\w+ \d{2} — \w+ \d{2}|\W+ \d{2} — \W+ \d{4}|\W+ \d{2} — \W+ \d{2}`) // Июнь 2014 — Август 20143 месяца
	items := re.FindAllString(periodText, -1)
	if len(items) > 0 {
		return items[0]
	} else {
		re = regexp.MustCompile(`\W+ \d{4} —.*?\d|\w+ \d{4} —.*?\d|\W+ \d{2} —.*?\d|\w+ \d{2} —.*?\d`) // Июнь 2014 — по настоящее время3 месяца
		item := re.FindString(periodText)
		return removeLasChar(item)
	}
}

func getExperienceDuration(html *colly.HTMLElement) string {
	durationText := html.ChildText("div.bloko-column.bloko-column_xs-4.bloko-column_s-2.bloko-column_m-2.bloko-column_l-2")
	if durationText == "" {
		return ""
	}
	durationText = strings.ReplaceAll(durationText, "\u00a0", " ")
	re := regexp.MustCompile(`\W+ \d{4} — \W+ \d{4}|\w+ \d{4} — \w+ \d{4}|\w+ \d{2} — \w+ \d{4}|\w+ \d{2} — \w+ \d{2}|\W+ \d{2} — \W+ \d{4}|\W+ \d{2} — \W+ \d{2}`)
	items := re.Split(durationText, -1)
	if len(items) > 1 {
		durationStr := items[len(items)-1]
		return durationStr
	} else {
		re = regexp.MustCompile(`\W+ \d{4} —.*?\d|\w+ \d{4} —.*?\d|\W+ \d{4} —.*?\d|\w+ \d{4} —.*?\d|\W+ \d{2} —.*?\d|\w+ \d{2} —.*?\d`)
		items := strings.Split(re.FindString(durationText), " ")
		item := items[len(items)-1]
		re = regexp.MustCompile(`\d+`)
		items = re.FindAllString(item, -1)
		item = items[len(items)-1]

		re = regexp.MustCompile(`\W+ \d{4} — .*?\d+|\w+ \d{4} — .*?\d+|\W+ \d{2} — .*?\d+|\w+ \d{2} — .*?\d+`)
		newItems := re.Split(durationText, -1)
		durationStr := item + strings.Join(newItems, " ")
		return durationStr
	}
}

func getGlobalExperience(html *colly.HTMLElement) string {
	var (
		exp   string
		reSub = regexp.MustCompile(`Опыт работы|Work experience|`)
	)
	html.ForEach("span.resume-block__title-text.resume-block__title-text_sub", func(i int, h *colly.HTMLElement) {
		if strings.Contains(h.Text, "Опыт работы") || strings.Contains(h.Text, "Work experience") {
			exp = reSub.ReplaceAllString(h.Text, "")
			exp = strings.ReplaceAll(exp, "\u00a0", " ")
		}

	})
	return exp
}
