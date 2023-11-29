package headhunter

import (
	"fmt"
	"regexp"
	"resumes/internal/database/resume"
	"resumes/internal/models"
	"resumes/pkg/logger"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

// Отправная точка для запуска парсинга по профессиям
func (hh *HeadHunter) ParseAll() {
	for _, pos := range hh.Positions {
		hh.ScrapePosition(pos)
	}
}

// Для того, чтобы спарсить профессию, нужно объединить название профессии с другими наименованиями и убрать из этого списка дубликаты
// После чего итерируемся по всем наименованиям ОДНОЙ профессии и возвращаем список резюме
func (hh *HeadHunter) ScrapePosition(pos models.Position) {
	hh.CurrentPositionId = pos.Id
	hh.PositionName = pos.Name
	pos.OtherNames = append(pos.OtherNames, pos.Name)
	r := resume.NewRep(hh.DB)
	for _, name := range getUniqueNames(pos.OtherNames) {
		hh.CurrentPositionName = name
		resumes := hh.CollectAllResumesByQuery()
		if len(resumes) == 0 {
			logger.Log.Printf("Не нашлось резюме для профессии: %s", name)
			continue
		}
		err := r.SaveResumes(resumes)
		if err != nil {
			logger.Log.Printf("ОШИБКА при сохранении %d резюме: %s", len(resumes), err)
		} else {
			logger.Log.Printf("УСПЕШНО сохранили %d резюме", len(resumes))
		}
	}
	return
}

// Здесь нужно реализовать определение количество доступных резюме
// Если количество меньше 2000, ищем резюме без привязки к городу
// Если количество превышает 2000, то ищем тотечно по городам все резюме (чтобы охватить как можно больше)
func (hh *HeadHunter) CollectAllResumesByQuery() (resumes []models.Resume) {
	resumes = hh.FindResumesInRussia()

	// for _, city := range hh.Cities {
	// 	items := hh.FindResumesInCurrentCity(city)
	// 	resumes = append(resumes, items...)
	// }

	return
}

// Для поиска без привязки к городу достаточно передать методу FindResumesInCurrentCity пустую структуру города
func (hh *HeadHunter) FindResumesInRussia() (resumes []models.Resume) {
	return hh.FindResumesInCurrentCity(models.City{})
}

// Здесь мы составляем поисковый запрос под конкретный город, итерируемся по всем страницам выдачи и собираем резюме
func (hh *HeadHunter) FindResumesInCurrentCity(city models.City) (resumes []models.Resume) {
	hh.CurrentCityId = city.HeadhunterID
	hh.CurrentCityEdwicaId = city.EdwicaID
	var pageNum int
	for pageNum < 50 {
		url := fmt.Sprintf("%s&page=%d", hh.CreateQuery(), pageNum)
		pageResumes := hh.CollectResumesFromPage(url)
		if len(pageResumes) == 0 {
			break
		}
		pageNum++
		logger.Log.Printf("Количество вакансий на странице №:%d\tКоличество резюме:%d\tПрофессия:%s:%s\tГород:%d\n", pageNum, len(pageResumes), hh.PositionName, hh.CurrentPositionName, city.HeadhunterID)
		resumes = append(resumes, pageResumes...)
	}
	return
}

// Собираем из выдачи ссылки на резюме и парсим их
func (hh *HeadHunter) CollectResumesFromPage(url string) (resumes []models.Resume) {
	var (
		resumesUrl = []string{}
		wg         sync.WaitGroup
	)
	html, err := getHTMLBody(url)
	if err != nil {
		logger.Log.Printf("Ошибка при получении HTML:%s\tURL:%s\n", err, url)
		return
	}
	// Находим ссылки на резюме
	html.ForEach("div[data-qa='resume-serp__resume']", func(i int, h *colly.HTMLElement) {
		itemUrl := "https://hh.ru" + h.ChildAttr("a", "href")
		resumesUrl = append(resumesUrl, itemUrl)
	})

	// Параллельно парсим резюме
	wg.Add(len(resumesUrl))
	for _, i := range resumesUrl {
		go hh.PutResumeToArrayByUrl(i, &wg, &resumes)
	}
	wg.Wait()
	return
}

// Не знаю почему не использую каналы для синхронизации данных, но здесь сохраняю данные в переданный слайс
func (hh *HeadHunter) PutResumeToArrayByUrl(url string, wg *sync.WaitGroup, resumes *[]models.Resume) {
	var resume models.Resume
	defer wg.Done()
	html, err := getHTMLBody(url)
	if err != nil {
		logger.Log.Printf("Ошибка при получении HTML:%s\tURL:%s\n", err, url)
		return
	}
	salary, currency := getSalary(html)

	resume.Url = url
	resume.Platform = "hh"
	resume.Salary = salary
	resume.Currency = currency
	resume.Id = getId(url)
	resume.City = getCity(html)
	resume.Title = getTitle(html)
	resume.SpecsList = getSpecs(html)
	resume.SkillsList = getSkills(html)
	resume.LanguagesList = getLanguages(html)
	resume.EducationList = getEducation(html)
	resume.ExperienceList = getExperience(html)
	resume.GlobalExperience = getGlobalExperience(html)
	resume.PositionId = hh.CurrentPositionId
	resume.CityId = hh.getEdwicaIdByCity(getCity(html))

	*resumes = append(*resumes, resume)
	return
}

// У нас не будет id правильного города, если парсим города без привязки к городу
// Чтобы данные нормально хранились в БД, нам нужно по названию городу найти его edwica_id
// Если edwica_id не найден, то в БД будет храниться название города и id = 0
func (hh *HeadHunter) getEdwicaIdByCity(name string) int {
	for _, i := range hh.Cities {
		if strings.ToLower(name) == strings.ToLower(i.Name) {
			return i.EdwicaID
		}
	}
	return 0
}

// Вырезаем id из ссылки резюме
func getId(url string) string {
	re := regexp.MustCompile(`.*?resume\/|\?.*`)
	id := re.ReplaceAllString(url, "")
	return id
}

func getCity(html *colly.HTMLElement) (city string) {
	html.ForEach("span[data-qa='resume-personal-address']", func(i int, h *colly.HTMLElement) {
		if i == 0 {
			city = h.Text
		}
	})
	return
}

func getTitle(html *colly.HTMLElement) (title string) {
	html.ForEach("div.resume-block__title-text-wrapper", func(i int, h *colly.HTMLElement) {
		if i == 0 {
			title = h.ChildText("h2")
		}
	})
	return
}

func getSkills(html *colly.HTMLElement) (skills []string) {
	html.ForEach("div.bloko-tag-list div.bloko-tag.bloko-tag_inline.bloko-tag_countable", func(i int, h *colly.HTMLElement) {
		skills = append(skills, h.ChildText("span"))
	})
	return
}

func getLanguages(html *colly.HTMLElement) (languages []models.Language) {
	html.ForEach("div[data-qa='resume-block-languages'] div.bloko-tag.bloko-tag_inline", func(i int, h *colly.HTMLElement) {
		text := strings.Split(h.Text, " — ")
		languages = append(languages, models.Language{
			Name:  text[0],
			Level: text[1],
		})
	})
	return
}

func getSpecs(html *colly.HTMLElement) (specs []string) {
	html.ForEach("li.resume-block__specialization", func(i int, h *colly.HTMLElement) {
		specs = append(specs, h.Text)
	})
	return
}

func getSalary(html *colly.HTMLElement) (salary int, currency string) {
	var reDigit = regexp.MustCompile(`\d+`)
	html.ForEach("span.resume-block__salary", func(i int, h *colly.HTMLElement) {
		text := strings.ReplaceAll(h.Text, "\u2009", "")
		text = strings.ReplaceAll(text, "\u00a0", " ")
		salary, _ = strconv.Atoi(reDigit.FindString(text))
		currency = reDigit.ReplaceAllString(text, "")
		if strings.Contains(currency, "руб") || strings.Contains(currency, "₽") {
			currency = "RUB"
		}
	})
	return
}
