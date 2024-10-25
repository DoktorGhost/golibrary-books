package providers

import (
	"github.com/DoktorGhost/golibrary-books/internal/services"
)

type ServiceProvider struct {
	authorService *services.AuthorService
	bookService   *services.BookService
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) RegisterDependencies(provider *RepositoryProvider) {
	s.authorService = services.NewAuthorService(provider.usersRepositoryPostgres)
	s.bookService = services.NewBookService(provider.usersRepositoryPostgres)
}
