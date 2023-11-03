package scrapper

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
	"time"
	"web-scrapper/pkg/database"
	"web-scrapper/pkg/environment"
	"web-scrapper/pkg/telegram"
)

var C *colly.Collector
var page int

func InitScrapper() {
}

func Run() {
	C = colly.NewCollector()

	C.OnHTML("div.vacancy-serp-item__layout div.serp-item-controls a.bloko-button", getVacancy)

	C.OnHTML("div.pager span.pager-item-not-in-short-range span.pager-item-not-in-short-range a", getTotalPages)

	C.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	executeVacancy()

	ticker := time.NewTicker(2 * time.Hour)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				executeVacancy()
			}
		}
	}()
	select {}

}

func executeVacancy() {
	var baseUrl string = environment.GetValue("BASE_URL")

	for i := 0; i <= page; i++ {
		url := fmt.Sprintf("/search/vacancy?text=Frontend&area=113&page=%d", i)
		C.Visit(baseUrl + url)
		fmt.Println(url)
		fmt.Println("Scrapping page " + strconv.Itoa(i) + " completed")

	}
}

func getVacancy(element *colly.HTMLElement) {
	var message string
	var baseUrl string = environment.GetValue("BASE_URL")
	var qString []string = strings.Split(element.Attr("href"), "=")
	var vacancyId string = strings.Split(qString[1], "&")[0]
	isExist := database.IsVacancyExist(vacancyId)
	vacancyC := colly.NewCollector()
	url := baseUrl + "/vacancy/" + vacancyId

	if !isExist {
		database.AddVacancy(vacancyId)
		vacancyC.OnHTML("div.vacancy-title h1.bloko-header-section-1", func(element *colly.HTMLElement) {
			//addToMessage("Название", element, &message)
			message = message + "Название: " + element.Text + "\n"
		})
		vacancyC.OnHTML("div[data-qa=vacancy-salary] span.bloko-header-section-2", func(element *colly.HTMLElement) {
			//addToMessage("Зарплата", element, &message)
			message = message + "Зарплата: " + element.Text + "\n"
		})
		vacancyC.OnHTML("span[data-qa=vacancy-experience]", func(element *colly.HTMLElement) {
			//addToMessage("Опыт", element, &message)
			message = message + "Опыт: " + element.Text + "\n"
		})
		//vacancyC.OnHTML("div.bloko-tag-list", func(element *colly.HTMLElement) {
		//	addToMessage("Навыки", element)
		//})
		vacancyC.OnHTML("p.vacancy-creation-time-redesigned", func(element *colly.HTMLElement) {
			//addToMessage("Дата создания", element, &message)
			message = message + "Дата создания: " + element.Text + "\n"
		})
		vacancyC.OnHTML("div.g-user-content", func(element *colly.HTMLElement) {
			//addToMessage("Описание", element, &message)
			message = message + "Описание: " + element.Text + "\n"
		})

		vacancyC.OnScraped(func(r *colly.Response) {
			message = message + url
			//fmt.Print(message)
			telegram.Send(message)
			//fmt.Println("Scraped", r.Request.URL)
		})
	}
	vacancyC.Visit(url)
}

func addToMessage(key string, element *colly.HTMLElement, message *string) {
	*message = *message + fmt.Sprintf("%s: %s", key, element.Text)
	*message = *message + "\n"
}

func getTotalPages(element *colly.HTMLElement) {
	var err error

	page, err = strconv.Atoi(element.Text)
	if err != nil {
		panic(err)
	}

	fmt.Println("Total pages: " + strconv.Itoa(page))
}
