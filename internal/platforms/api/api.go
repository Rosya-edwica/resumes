package api

import (
	"resumes/internal/models"
	"sync"
)

type API interface {
	// Метод создания ссылки, в которой будут лежать данные
	CreateQuery() (query string)

	// Подсчет резюме одной профессии по всей России для того, чтобы определить какой метод нам использовать FindResumesInRussia() или FindResumesInCurrentCity()
	CountResumesByQuery(url string) (count int)

	// Сбор всех вакансий с одного запроса по API (все вакансии профессии в городе)
	CollectAllResumesByQuery(position models.Position) (resumes []models.Resume)

	// Поиск резюме по всей России без привязки к городу
	FindResumesInRussia() (resumes []models.Resume)

	// Поиск резюме по конкретному городу для популярных профессий, которых больше 2000 на платформе
	FindResumesInCurrentCity(city models.City) (resumes []models.Resume)

	// Сбор всех резюме с одной страницы запроса
	CollectResumesFromPage(url string) (resumes []models.Resume)

	// Сбор одного конкретного резюме
	PutResumeToArrayByUrl(id string, wg *sync.WaitGroup, resumes *[]models.Resume)
}
