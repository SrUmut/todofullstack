package storage

type todo struct {
	ID          int
	Title       string
	Description string
}

func NewTodo(id int, title string, description string) *todo {
	return &todo{
		ID:          id,
		Title:       title,
		Description: description,
	}
}
