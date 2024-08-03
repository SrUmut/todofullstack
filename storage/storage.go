package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	AddTodo(int, string, string) error
	GetTodos() ([]*todo, error)
	GetNextID() int
	RemoveByID(string) error
}

type postgres struct {
	db     *sql.DB
	lastID *int
}

func NewPostgres() (*postgres, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}

	postgres := &postgres{db: db}
	id, err := postgres.getLastId()
	if err != nil {
		return nil, err
	}
	postgres.lastID = &id
	return postgres, nil
}

func (p *postgres) AddTodo(id int, title string, description string) error {
	query := `
	INSERT INTO todos
	VALUES ($1, $2, $3);`
	if _, err := p.db.Exec(query, id, title, description); err != nil {
		return err
	}
	return nil
}

func (p *postgres) GetTodos() ([]*todo, error) {
	query := `SELECT * FROM todos;`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]*todo, 0)
	for rows.Next() {
		var t todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Description); err != nil {
			return nil, err
		}
		todos = append(todos, &t)
	}
	return todos, nil
}

func (p *postgres) RemoveByID(id string) error {
	query := `DELETE FROM todos WHERE id = $1;`
	if _, err := p.db.Exec(query, id); err != nil {
		return err
	}
	return nil
}

// get the next available id
func (p *postgres) GetNextID() int {
	*p.lastID++
	return *p.lastID
}

func initDB() (*sql.DB, error) {
	connStr := "user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := createTable(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
	id SERIAL PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
	description VARCHAR(255) NOT NULL
	);`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func (p *postgres) getLastId() (int, error) {
	query := `SELECT MAX(id) from todos;`
	var id sql.NullInt64
	p.db.QueryRow(query).Scan(&id)
	if !id.Valid {
		return 0, nil
	}
	return int(id.Int64), nil
}
