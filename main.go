package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"strings"
	"web-scrapper/pkg/postgres"
)

const baseUrl string = "https://chelyabinsk.hh.ru"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	postgres.ConnectToDb()
}

func main() {
	c := colly.NewCollector()
	
	var page int
	var err error
	
	c.OnHTML("div.vacancy-serp-item__layout div.serp-item-controls a.bloko-button", getVacancy)
	
	c.OnHTML("div.pager span.pager-item-not-in-short-range span.pager-item-not-in-short-range a", func(element *colly.HTMLElement) {
		//fmt.Println("Visiting", element.Text)
		page, err = strconv.Atoi(element.Text)
		if err != nil {
			panic(err)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	
	for i := 0; i <= page; i++ {
		url := fmt.Sprintf("/search/vacancy?text=front-end&area=113&page=%d", i)
		c.Visit(baseUrl + url)
		
		fmt.Println(url)
	}
}

func getVacancy(element *colly.HTMLElement) {
	var qString []string = strings.Split(element.Attr("href"), "=")
	var vacancyId string = strings.Split(qString[1], "&")[0]
	
	isExist := postgres.IsVacancyExist(vacancyId)
	vacancyC := colly.NewCollector()
	url := baseUrl + "/vacancy/" + vacancyId
	
	if !isExist {
		//postgres.AddVacancy(vacancyId)
		vacancyC.OnHTML("span.bloko-header-section-2", func(element *colly.HTMLElement) {
			//fmt.Println(element.Text)
			fmt.Println(url)
		})
	}
	
	vacancyC.Visit(url)
}
