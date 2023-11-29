package headhunter

import (
	"net/url"
	"regexp"
	"resumes/internal/models"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-faster/errors"
	"github.com/gocolly/colly"
	"github.com/jmoiron/sqlx"
)

var HEADERS = map[string]string{
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
}

type HeadHunter struct {
	Positions           []models.Position
	Cities              []models.City
	DB                  *sqlx.DB
	CurrentPositionId   int
	CurrentPositionName string
	PositionName        string
	CurrentCityId       int
	CurrentCityEdwicaId int
}

// FIXME: RussiaCityId ONLY FOR HH
const RussiaCityId = 113

func (hh *HeadHunter) CreateQuery() (query string) {
	var params = make(url.Values)
	var strCityid string
	if hh.CurrentCityId != 0 {
		strCityid = strconv.Itoa(hh.CurrentCityId)
	} else {
		strCityid = strconv.Itoa(RussiaCityId)
	}

	params.Add("text", hh.CurrentPositionName) // Какую профессию ищем?
	params.Add("area", strCityid)              // В каком городе ищем?
	params.Add("no_magic", "true")
	params.Add("ored_clusters", "true")
	params.Add("order_by", "relevance")
	params.Add("items_on_page", "100")
	params.Add("search_period", "0")
	params.Add("logic", "normal")
	params.Add("pos", "position")
	params.Add("exp_period", "all_time")
	params.Add("exp_company_size", "any")
	params.Add("hhtmFrom", "resume_search_result")
	params.Add("job_search_status", "unknown")
	params.Add("job_search_status", "active_search")
	params.Add("job_search_status", "looking_for_offers")
	params.Add("job_search_status", "has_job_offer")
	params.Add("job_search_status", "accepted_job_offer")
	params.Add("job_search_status", "not_looking_for_job")

	return "https://hh.ru/search/resume?" + params.Encode()
}

func getUniqueNames(names []string) (unique []string) {
	allKeys := make(map[string]bool)
	for _, item := range names {
		lower := strings.TrimSpace(strings.ToLower(item))
		if _, value := allKeys[lower]; !value {
			if len(item) < 2 {
				continue
			}
			allKeys[lower] = true
			unique = append(unique, strings.TrimSpace(item))
		}
	}
	return
}

func getHTMLBody(url string) (body *colly.HTMLElement, err error) {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		for key, value := range HEADERS {
			r.Headers.Set(key, value)
		}
	})
	c.OnHTML("body", func(h *colly.HTMLElement) {
		body = h
	})

	err = c.Visit(url)
	if err != nil {
		return nil, errors.Wrap(err, "get html")
	}
	return body, nil
}

func removeLasChar(str string) string {
	for len(str) > 0 {
		_, size := utf8.DecodeLastRuneInString(str)
		return str[:len(str)-size]
	}
	return str
}

// Приводим такую строку "3 года 2 месяца" в количество месяцев: 38
func experienceToMonths(text string) int {
	var (
		reSpace       = regexp.MustCompile(` +`)
		reMonth       = regexp.MustCompile(`\d{2} м|\d{2} m|\d м|\d m`)
		reYear        = regexp.MustCompile(`\d{2} г|\d{2} y|\d г|\d y|\d{2} л|\d л`)
		months, years int
	)

	text = reSpace.ReplaceAllString(text, " ")
	monthsText := reMonth.FindString(text)
	items := strings.Split(monthsText, " ")
	if len(items) > 0 {
		monthsText = items[0]
		months, _ = strconv.Atoi(monthsText)
	} else {
		months = 0
	}

	yearsText := reYear.FindString(text)
	items = strings.Split(yearsText, " ")
	if len(items) > 0 {
		yearsText = items[0]
		years, _ = strconv.Atoi(yearsText)
	}

	if years > 0 {
		return years*12 + months
	} else {
		return months
	}
}
