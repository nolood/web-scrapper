package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var Db *sql.DB

func ConnectToDb() {
	user, exist := os.LookupEnv("PG_USERNAME")
	if !exist {
		fmt.Println("PG_USERNAME does not exist")
	}

	host, exist := os.LookupEnv("PG_HOST")
	if !exist {
		fmt.Println("PG_HOST does not exist")
	}

	password, exist := os.LookupEnv("PG_PASSWORD")
	if !exist {
		fmt.Println("PG_PASSWORD does not exist")
	}

	dbname, exist := os.LookupEnv("PG_DB")
	if !exist {
		fmt.Println("PG_DB does not exist")
	}

	port, exist := os.LookupEnv("PG_PORT")
	if !exist {
		fmt.Println("PG_PORT does not exist")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	//defer Db.Close()

	err = Db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("CREATE TABLE IF NOT EXISTS vacancies (id integer PRIMARY KEY)")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully connected to %s database!\n", dbname)
}

func AddVacancy(vacancyId string) {
	Db.Exec("INSERT INTO vacancies (id) VALUES ($1)", vacancyId)
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
