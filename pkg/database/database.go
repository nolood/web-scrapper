package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"web-scrapper/pkg/environment"
)

var Db *sql.DB

func ConnectToDb() {
	user := environment.GetValue("PG_USERNAME")
	host := environment.GetValue("PG_HOST")
	password := environment.GetValue("PG_PASSWORD")
	dbname := environment.GetValue("PG_DB")
	port := environment.GetValue("PG_PORT")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = Db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully connected to %s database!\n", dbname)
}

func InitDatabase() {
	ConnectToDb()

	_, err := Db.Exec("CREATE TABLE IF NOT EXISTS vacancies (id integer PRIMARY KEY)")
	if err != nil {
		panic(err)
	}
}

func AddVacancy(vacancyId string) {
	_, err := Db.Exec("INSERT INTO vacancies (id) VALUES ($1)", vacancyId)
	if err != nil {
		return
	}
}

func IsVacancyExist(vacancyId string) bool {
	rows, err := Db.Query("SELECT id FROM vacancies WHERE id = $1", vacancyId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}
