package handlers

import (
	"encoding/json"
	"github.com/DoktorGhost/golibrary-books/internal/entities"
	"github.com/DoktorGhost/golibrary-books/internal/providers"
	"github.com/DoktorGhost/golibrary-books/internal/repositories/postgres/dao"
	"io"
	"net/http"
	"strconv"
)

func handlerAddAuthor(useCaseProvider *providers.UseCaseProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
			return
		}

		// Чтение тела запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Декодирование JSON из тела запроса
		var author entities.AuthorRequest
		if err := json.Unmarshal(body, &author); err != nil {
			http.Error(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		id, err := useCaseProvider.BookUseCase.AddAuthor(author.Name, author.Surname, author.Patronymic)
		if err != nil {
			http.Error(w, "Ошибка при добавлении автора: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusCreated)
		responseMessage := "Автор успешно добавлен, ID: " + strconv.Itoa(id)
		w.Write([]byte(responseMessage))
	}
}

func handlerAddBook(useCaseProvider *providers.UseCaseProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
			return
		}

		// Чтение тела запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Декодирование JSON из тела запроса
		var book entities.BookRequest
		if err := json.Unmarshal(body, &book); err != nil {
			http.Error(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		authorID, err := strconv.Atoi(book.AuthorID)
		if err != nil {
			http.Error(w, "Ошибка при добавлении книги: "+err.Error(), http.StatusInternalServerError)
			return
		}
		id, err := useCaseProvider.BookUseCase.AddBook(dao.BookTable{Title: book.Title, AuthorID: authorID})
		if err != nil {
			http.Error(w, "Ошибка при добавлении книги: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusCreated)
		responseMessage := "Книга успешно добавлена, ID: " + strconv.Itoa(id)
		w.Write([]byte(responseMessage))
	}
}

func handlerGetAllBooks(useCaseProvider *providers.UseCaseProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
			return
		}

		books, err := useCaseProvider.BookUseCase.GetAllBookWithAuthor()
		if err != nil {
			http.Error(w, "Ошибка получения книг: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Установите заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")

		// Преобразуйте книги в JSON
		response, err := json.Marshal(books)
		if err != nil {
			http.Error(w, "Ошибка при преобразовании данных в JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(response)
	}
}

func handlerGetAllAuthors(useCaseProvider *providers.UseCaseProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Неправильный метод", http.StatusMethodNotAllowed)
			return
		}

		books, err := useCaseProvider.BookUseCase.GetAllAuthorWithBooks()
		if err != nil {
			http.Error(w, "Ошибка получения авторов: "+err.Error(), http.StatusInternalServerError)

			return
		}

		// Установите заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")

		// Преобразуйте книги в JSON
		response, err := json.Marshal(books)
		if err != nil {
			http.Error(w, "Ошибка при преобразовании данных в JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(response)
	}
}
