package providers

import (
	"github.com/DoktorGhost/golibrary-books/internal/repositories/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryProvider struct {
	db                      *pgxpool.Pool
	usersRepositoryPostgres *postgres.BookRepository
}

func NewRepositoryProvider(db *pgxpool.Pool) *RepositoryProvider {
	return &RepositoryProvider{db: db}
}

func (r *RepositoryProvider) RegisterDependencies() {
	r.usersRepositoryPostgres = postgres.NewPostgresRepository(r.db)
}
