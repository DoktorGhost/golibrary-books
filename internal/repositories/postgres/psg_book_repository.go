package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DoktorGhost/golibrary-books/internal/repositories/postgres/dao"
)

func (s *BookRepository) CreateBook(book dao.BookTable) (int, error) {
	var id int
	query := `INSERT INTO books (title, author_id) VALUES ($1, $2) RETURNING id`
	err := s.db.QueryRow(context.Background(), query, book.Title, book.AuthorID).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка добавления книги с названием '%s' и автором ID %d: %v", book.Title, book.AuthorID, err)
	}

	return id, nil
}

func (s *BookRepository) GetBookByID(id int) (dao.BookTable, error) {
	var result dao.BookTable
	query := `SELECT title, author_id FROM books WHERE id = $1`
	err := s.db.QueryRow(context.Background(), query, id).Scan(&result.Title, &result.AuthorID)

	if err != nil {
		if err == sql.ErrNoRows {
			return dao.BookTable{}, fmt.Errorf("книга с ID %d не найдена", id)
		}
		return dao.BookTable{}, fmt.Errorf("ошибка получения книги: %v", err)
	}
	result.ID = id

	return result, nil
}

func (s *BookRepository) GetAllBooks() ([]dao.BookTable, error) {
	query := `SELECT id, title, author_id FROM books;`
	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var books []dao.BookTable

	for rows.Next() {
		var book dao.BookTable

		err := rows.Scan(&book.ID, &book.Title, &book.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %v", err)
	}

	return books, nil
}
