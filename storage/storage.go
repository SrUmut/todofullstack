package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	AddTodo(username, title string, description string) (id int64, err error)
	GetTodos(username string) ([]*todo, error)
	//GetNextID() int
	RemoveTodo(username, tid string) error
	AddUser(username string, enc_password string) error
	GetPassword(username string) (string, error)
}

type postgres struct {
	db *sql.DB
}

func (p *postgres) AddUser(username, enc_password string) error {
	query := `
	SELECT COUNT(*) FROM "user"
	WHERE username = $1;
	`
	var userExist bool
	if err := p.db.QueryRow(query, username).Scan(&userExist); err != nil {
		return err
	}

	if userExist {
		return fmt.Errorf("username is taken")
	}

	query = `
	INSERT INTO "user" (username, password)
	VALUES ($1, $2);
	`

	if _, err := p.db.Exec(query, username, enc_password); err != nil {
		return err
	}

	return nil
}

func NewPostgres() (*postgres, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}

	postgres := &postgres{db: db}
	return postgres, nil
}

// adds a todo and return its todoid
func (p *postgres) AddTodo(username, title string, description string) (int64, error) {
	uid, err := p.getID(username)
	if err != nil {
		return -1, err
	}
	if uid == "" {
		return -1, fmt.Errorf("no account with username: %v", username)
	}

	query := `
	INSERT INTO todo (id, title, description)
	VALUES ($1, $2, $3);`
	if _, err := p.db.Exec(query, uid, title, description); err != nil {
		return -1, err
	}

	var todoid sql.NullInt64
	query = `SELECT "todoid" FROM todo 
	WHERE id = $1 AND title = $2;`
	p.db.QueryRow(query, uid, title).Scan(&todoid)
	if !todoid.Valid {
		return -1, fmt.Errorf("an error occured while adding a todo")
	}

	return todoid.Int64, nil
}

func (p *postgres) GetTodos(username string) ([]*todo, error) {
	id, err := p.getID(username)
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, fmt.Errorf("no account with username: %v", username)
	}

	query := `
	SELECT * FROM todo
	WHERE id = $1;`
	rows, err := p.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]*todo, 0)
	for rows.Next() {
		var t todo
		if err := rows.Scan(&t.ID, &t.TodoID, &t.Title, &t.Description); err != nil {
			return nil, err
		}
		todos = append(todos, &t)
	}
	return todos, nil
}

func (p *postgres) RemoveTodo(username, tid string) error {
	uid, err := p.getID(username)
	if err != nil {
		return err
	}
	query := `DELETE FROM todo WHERE id = $1 and todoid = $2;`
	if _, err := p.db.Exec(query, uid, tid); err != nil {
		return err
	}
	return nil
}

func initDB() (*sql.DB, error) {
	dbpass := os.Getenv("DBPASS")
	connStr := fmt.Sprintf("user=postgres dbname=postgres password=%s sslmode=disable", dbpass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// create user and todo table
func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS "user" (
	id SERIAL PRIMARY KEY,
	username VARCHAR(64) NOT NULL,
	password VARCHAR(100) NOT NULL
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `
	CREATE TABLE IF NOT EXISTS todo (
	id BIGINT NOT NULL,
	todoid SERIAL NOT NULL,
	title VARCHAR(100) NOT NULL,
	description VARCHAR(255) NOT NULL,
	PRIMARY KEY (id, title)
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

func (p *postgres) GetPassword(username string) (string, error) {
	query := `
	SELECT password from "user"
	WHERE username = $1;
	`
	var pass sql.NullString
	if err := p.db.QueryRow(query, username).Scan(&pass); err != nil {
		return "", err
	}
	if !pass.Valid {
		return "", nil
	}

	return pass.String, nil
}

func (p *postgres) getID(username string) (string, error) {
	query := `
	SELECT id FROM "user"
	WHERE username = $1;
	`
	var id sql.NullString
	if err := p.db.QueryRow(query, username).Scan(&id); err != nil {
		return "", err
	}

	if !id.Valid {
		return "", nil
	}

	return id.String, nil
}
