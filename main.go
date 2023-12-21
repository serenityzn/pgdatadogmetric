package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type pgConnect struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func dbConnect(connect pgConnect) (*sql.DB, error) {

	psqlConnect := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connect.host, connect.port, connect.user, connect.password, connect.dbname)

	db, err := sql.Open("postgres", psqlConnect)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Unable to Ping Database! [ERROR]")
		return nil, err
	}

	fmt.Println("Connected! [OK]")
	return db, nil
}

func dbQuery(db *sql.DB) (*sql.Row, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM products")

	return row, nil
}

func main() {

	mydb, err := configInit()
	if err != nil {
		panic(err)
	}

	db, err := dbConnect(mydb)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	row, err := dbQuery(db)
	if err != nil {
		panic(err)
	}

	var count int
	err = row.Scan(&count)
	if err != nil {
		panic(err)
	}

	fmt.Println(count)

	/*err = datadogCreateMetricTag()
	if err != nil {
		panic(err)
	}*/

	err = datadogSubmitMetric(float64(count))
	if err != nil {
		panic(err)
	}

}
