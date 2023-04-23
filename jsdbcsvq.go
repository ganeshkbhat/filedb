package jsdbcsvq

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mithrandie/csvq/lib/query"

	_ "github.com/mithrandie/csvq-driver"
)

func FileSQLQuery(directorypath string, querystring string, datastructure interface{}) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	directorypath = directorypath // "/path/to/data/directory"

	db, err := sql.Open("csvq", directorypath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	// queryString := "SELECT id, first_name, country_code FROM `users.csv` WHERE id = '12'"
	// r := db.QueryRowContext(ctx, queryString)
	r := db.QueryRowContext(ctx, querystring)

	var (
		id          int
		firstName   string
		countryCode string
	)

	if err := r.Scan(&id, &firstName, &countryCode); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No Rows.")
		} else if csvqerr, ok := err.(query.Error); ok {
			panic(fmt.Sprintf("Unexpected error: Number: %d  Message: %s", csvqerr.Number(), csvqerr.Message()))
		} else {
			panic("Unexpected error: " + err.Error())
		}
	} else {
		fmt.Printf("Result: [id]%3d  [first_name]%10s  [country_code]%3s\n", id, firstName, countryCode)
	}
	return r
}
