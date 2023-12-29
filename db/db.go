package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func NewPqClient(user, pass, host, port, database, schema string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", user, pass, host, port, database, schema)
	return NewPqClientByDSN(dsn)
}

func NewPqClientByDSN(dsnString string) (*sql.DB, error) {
	var err error
	dataSourceName := fmt.Sprintf("postgres://" + dsnString)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 5. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(5)
	// Set the maximum number of concurrently idle connections to 5. Setting this
	// to less than or equal to 0 will mean that no idle connections are retained.
	db.SetMaxIdleConns(5)
	// Sets the maximum amount of time a connection may be idle.
	db.SetConnMaxIdleTime(5 * time.Minute)
	// Set the maximum lifetime of a connection. Setting it to 0
	// means that there is no maximum lifetime and the connection is reused
	// forever (which is the default behavior).
	db.SetConnMaxLifetime(1 * time.Hour)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
