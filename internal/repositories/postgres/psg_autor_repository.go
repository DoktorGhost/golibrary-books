package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DoktorGhost/golibrary-books/internal/repositories/postgres/dao"
)

func (s *BookRepository) CreateAuthor(name string) (int, error) {
	var id int
	query := `INSERT INTO authors (name) VALUES ($1) RETURNING id`
	err := s.db.QueryRow(context.Background(), query, name).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка добавления записи: %v", err)
	}

	return id, nil
}

func (s *BookRepository) GetAuthorByID(id int) (dao.AuthorTable, error) {
	var result dao.AuthorTable
	query := `SELECT name FROM authors WHERE id = $1`
	err := s.db.QueryRow(context.Background(), query, id).Scan(&result.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return dao.AuthorTable{}, fmt.Errorf("автор с ID %d не найден", id)
		}
		return dao.AuthorTable{}, fmt.Errorf("ошибка получения автора: %v", err)
	}
	result.ID = id
	return result, nil
}

func (s *BookRepository) GetAllAuthors() ([]dao.AuthorTable, error) {
	query := `SELECT id, name FROM authors;`
	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var authors []dao.AuthorTable

	for rows.Next() {
		var author dao.AuthorTable

		err := rows.Scan(&author.ID, &author.Name)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		authors = append(authors, author)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %v", err)
	}

	return authors, nil
}
