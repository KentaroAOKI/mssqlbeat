package beater

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"

	"github.com/KentaroAOKI/mssqlbeat/config"
	_ "github.com/microsoft/go-mssqldb"
)

// PublishMssqlData publishes SQL data to the Beat client.
func (bt *mssqlbeat) PublishMssqlData(b *beat.Beat, input *config.Input, thread_no int) error {

	var db *sql.DB
	var err error

	// Create the connection string for the SQL server.
	connString := fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s;TLSMin=%s",
		input.MssqlserverHost,
		input.MssqlserverPort,
		input.MssqlserverUserId,
		input.MssqlserverPassword,
		input.MssqlserverDatabase,
		input.MssqlserverTlsmin)

	// Open a connection to the SQL server.
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
		return err
	}
	defer db.Close()

	// Create a context for the database operations.
	ctx := context.Background()

	// Read the last timestamp from the file.
	lastTime, err := readLastTime(input.Field, input.SqlTimeInitializeWithCurrentTime)
	if err != nil {
		log.Fatal("Error read last time: ", err.Error())
		return err
	}

	// Execute the SQL query with the last timestamp as a parameter.
	tsql := input.SqlQuery
	rows, err := db.QueryContext(ctx, tsql, sql.Named("LastTime", lastTime))
	if err != nil {
		log.Fatal("Error getting data: ", err.Error())
		return err
	}
	defer rows.Close()

	// Get the column names from the result set.
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("Error getting columns: ", err.Error())
		return err
	}

	// Prepare slices to hold the column values.
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Iterate over the rows in the result set.
	var maxTime time.Time
	for rows.Next() {
		var currentTime time.Time
		var hasCurrentTime bool

		// Scan the row values into the value pointers.
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal("Error scanning result: ", err.Error())
			return err
		}

		// Create a map to hold the column values.
		result := make(map[string]interface{})
		hasCurrentTime = false
		// Populate the map with the column values.
		for i, col := range columns {
			result[col] = values[i]
			// Get time for the sql_time_column.
			if col == input.SqlTimeColumn {
				var ok bool
				currentTime, ok = values[i].(time.Time)
				hasCurrentTime = ok || hasCurrentTime
			}
		}
		// Check if the sql_time_column is present.
		if !hasCurrentTime {
			errorMessage := fmt.Sprintf("%s column not found.", input.SqlTimeColumn)
			log.Fatalf("Error sql_time_column: %s", errorMessage)
			err = errors.New(errorMessage)
			return err
		}

		// Create a map for the Beat event fields.
		beat_event_fields := common.MapStr{
			"type":   b.Info.Name,
			"field":  input.Field,
			"thread": thread_no,
		}

		// Add the column values to the Beat event fields.
		for col, val := range result {
			beat_event_fields[input.FieldPrefix+col] = val
		}

		// Create a Beat event with the current timestamp and fields.
		event := beat.Event{
			Timestamp: currentTime,
			Fields:    beat_event_fields,
		}

		// Publish the event to the Beat client.
		bt.mu.Lock()
		bt.client.Publish(event)
		bt.mu.Unlock()

		// Update the maximum timestamp.
		if currentTime.After(maxTime) {
			maxTime = currentTime
		}
	}

	// Write the last timestamp to the file if it has been updated.
	if !maxTime.IsZero() {
		err = writeLastTime(input.Field, maxTime)
		if err != nil {
			return err
		}
	}
	db.Close()
	return nil
}
