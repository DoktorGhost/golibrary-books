package providers

import (
	usecases "github.com/DoktorGhost/golibrary-books/internal/usecase"
)

type UseCaseProvider struct {
	BookUseCase *usecases.BookUseCase
}

func NewUseCaseProvider() *UseCaseProvider {
	return &UseCaseProvider{}
}

func (ucp *UseCaseProvider) RegisterDependencies(provider *ServiceProvider) {
	ucp.BookUseCase = usecases.NewBookUseCase(provider.bookService, provider.authorService)

}
