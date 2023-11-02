package main

import (
	"web-scrapper/pkg/database"
	"web-scrapper/pkg/environment"
	"web-scrapper/pkg/scrapper"
	"web-scrapper/pkg/telegram"
)

func init() {
	environment.InitEnvironment()
	database.InitDatabase()
	telegram.InitTelegramBot()
	scrapper.InitScrapper()
}

func main() {
	scrapper.Run()
}
