package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
)

func main() {
	c := colly.NewCollector()
	
	var baseUrl string = "https://chelyabinsk.hh.ru"
	
	c.OnHTML("div.vacancy-serp-item__layout div.serp-item-controls a.bloko-button", func(element *colly.HTMLElement) {
		var qString []string = strings.Split(element.Attr("href"), "=")
		var vacancyId string = strings.Split(qString[1], "&")[0]
		url := baseUrl + "/vacancy/" + vacancyId
		vacancyC := colly.NewCollector()
		
		vacancyC.OnRequest(func(request *colly.Request) {
			fmt.Println("Visiting", request.URL)
		})
		
		vacancyC.OnHTML("span.bloko-header-section-2", func(element *colly.HTMLElement) {
			fmt.Println(element.Text)
		})
		
		vacancyC.Visit(url)
		
	})
	
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	
	c.Visit(baseUrl + "/search/vacancy?text=junior+front-end&area=113&page=0")
}
