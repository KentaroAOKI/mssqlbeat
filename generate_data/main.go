package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB
var server = "<server_name>.database.windows.net"
var port = 1433
var user = "<user_id>"
var password = "<password>"
var database = "<database>"

// CREATE TABLE TimeCount (
//     ID INT PRIMARY KEY IDENTITY(1,1),
//     CurrentTime DATETIME,
//     CountValue INT
// );

func main() {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()

	// Insert data every second
	count := 0
	for {
		count++
		err = insertTimeCount(ctx, count)
		if err != nil {
			log.Fatal("Error inserting data: ", err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

func insertTimeCount(ctx context.Context, count int) error {
	tsql := `INSERT INTO TimeCount (CurrentTime, CountValue) VALUES (@CurrentTime, @CountValue);`
	_, err := db.ExecContext(
		ctx,
		tsql,
		sql.Named("CurrentTime", time.Now()),
		sql.Named("CountValue", count))
	return err
}
