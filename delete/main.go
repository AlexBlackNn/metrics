package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func main() {
	db, err := sql.Open("pgx", "postgresql://app:app123@127.0.0.1:5432/metric_db?sslmode=disable")
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	stmt, err := db.Prepare(
		`INSERT INTO
    			app.counter_part (metric_id, name, value)
				VALUES ((SELECT uuid FROM app.types WHERE name = $1), $2, $3)`)
	if err != nil {
		fmt.Println("Unable to prepare statement:", err)
		return
	}
	for {
		_, err = stmt.Exec(
			"gauge",
			"test_gauge",
			42,
		)
		if err != nil {
			fmt.Println("Unable to insert to database:", err)
			return
		}

		rows, err := db.Query("SELECT name, statement, prepare_time FROM pg_prepared_statements")
		if err != nil {
			fmt.Println("Unable to query prepared statements:", err)
			return
		}
		defer rows.Close()

		var (
			name        string
			statement   string
			prepareTime time.Time
		)

		for rows.Next() {
			err = rows.Scan(&name, &statement, &prepareTime)
			if err != nil {
				fmt.Println("Unable to scan row:", err)
				continue
			}
			fmt.Printf("==========> name: %s \n, statement: %s \n, prepareTime: %s\n\n\n", name, statement, prepareTime)
		}

		time.Sleep(time.Second)
		fmt.Println("Inserted gauge value:", 42)
	}
}
