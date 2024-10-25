package handlers

import (
	"github.com/DoktorGhost/golibrary-books/internal/providers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(provider *providers.UseCaseProvider) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/author/add", handlerAddAuthor(provider))
	r.Post("/books/add", handlerAddBook(provider))
	r.Get("/books", handlerGetAllBooks(provider))
	r.Get("/authors", handlerGetAllAuthors(provider))

	return r
}
