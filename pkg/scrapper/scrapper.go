package scrapper

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
	"time"
	"web-scrapper/pkg/database"
	"web-scrapper/pkg/environment"
)

var C *colly.Collector
var page int
var message string

func InitScrapper() {
	C = colly.NewCollector()
}

func Run() {
	C.OnHTML("div.vacancy-serp-item__layout div.serp-item-controls a.bloko-button", getVacancy)

	C.OnHTML("div.pager span.pager-item-not-in-short-range span.pager-item-not-in-short-range a", getTotalPages)

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
		url := fmt.Sprintf("/search/vacancy?text=front-end&area=113&page=%d", i)
		C.Visit(baseUrl + url)
		fmt.Println(url)
	}
	fmt.Println("Scrapping completed")
}

func getVacancy(element *colly.HTMLElement) {
	message = ""
	var baseUrl string = environment.GetValue("BASE_URL")
	var qString []string = strings.Split(element.Attr("href"), "=")
	var vacancyId string = strings.Split(qString[1], "&")[0]
	isExist := database.IsVacancyExist(vacancyId)
	vacancyC := colly.NewCollector()
	url := baseUrl + "/vacancy/" + vacancyId

	if !isExist {
		database.AddVacancy(vacancyId)
		vacancyC.OnHTML("div.vacancy-title h1.bloko-header-section-1", func(element *colly.HTMLElement) {
			fmt.Println(element.Text)
			//telegram.Send(element.Text)
		})
	}
	vacancyC.Visit(url)
}

func getTotalPages(element *colly.HTMLElement) {
	var err error

	page, err = strconv.Atoi(element.Text)
	if err != nil {
		panic(err)
	}
}

func addToMessage(key, element *colly.HTMLElement) {
	message = fmt.Sprintf("%s: %s", key, message)
}
