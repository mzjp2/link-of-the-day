package storage

import (
	"database/sql"
	"fmt"
	"time"

	// This loads the postgres drivers.
	_ "github.com/lib/pq"
)

// New returns a postgres backed storage service.
func New(connect string) (Service, error) {
	// Connect postgres
	db, err := sql.Open("postgres", connect)
	if err != nil {
		return nil, fmt.Errorf("could not open connection: %v", err)
	}

	// Ping to connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping connection: %v", err)
	}

	// Create table if does not already exists
	query := "CREATE TABLE IF NOT EXISTS links (id serial NOT NULL, url VARCHAR not NULL, " +
		"count INTEGER DEFAULT 0, date_added DATE not NULL, scheduled DATE UNIQUE not NULL);"

	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("could not create table: %v", err)
	}

	return &postgres{db}, nil
}

type postgres struct{ db *sql.DB }

func (p *postgres) Save(url string, s, t time.Time) (int64, error) {
	var id int64
	query := "INSERT INTO links(url,count,date_added,scheduled) VALUES($1,$2,$3,$4) returning id;"

	err := p.db.QueryRow(query, url, 0, t, s).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not save row: %v", err)
	}

	return id, nil
}

func (p *postgres) UpdateCount(id int) error {
	query := "UPDATE links SET count = count+1 WHERE id = $1;"
	_, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not update count value in db: %v", err)
	}
	return nil
}

func (p *postgres) Load(id int) (*Record, error) {
	var record Record
	query := "SELECT * FROM links where id=$1 limit 1;"

	err := p.db.QueryRow(query, id).Scan(&record.ID, &record.URL, &record.Count, &record.DateAdded, &record.Scheduled)
	if err != nil {
		return nil, fmt.Errorf("could not load row: %v", err)
	}

	return &record, nil
}

func (p *postgres) LoadLast() (*Record, error) {
	rows, err := p.db.Query("SELECT * FROM links ORDER BY id DESC LIMIT 1;")
	if err != nil {
		return nil, fmt.Errorf("could not load rows: %v", err)
	}
	defer rows.Close()

	var record Record
	if rows.Next() {
		err = rows.Scan(&record.ID, &record.URL, &record.Count, &record.DateAdded, &record.Scheduled)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
	} else {
		err = rows.Err()
		if err != nil {
			return nil, fmt.Errorf("could not iterate rows: %v", err)
		}
		return nil, nil
	}

	return &record, nil
}

func (p *postgres) LoadScheduled(s time.Time) (*Record, error) {
	rows, err := p.db.Query("SELECT * FROM links where scheduled=$1 limit 1;", s)
	if err != nil {
		return nil, fmt.Errorf("could not load rows: %v", err)
	}
	defer rows.Close()

	var record Record
	if rows.Next() {
		err = rows.Scan(&record.ID, &record.URL, &record.Count, &record.DateAdded, &record.Scheduled)
		if err != nil {
			return nil, fmt.Errorf("could not scan row %v", err)
		}
	} else {
		err = rows.Err()
		if err != nil {
			return nil, fmt.Errorf("could not iterate rows: %v", err)
		}
		return nil, nil
	}

	return &record, nil
}

func (p *postgres) Close() error { return p.db.Close() }
