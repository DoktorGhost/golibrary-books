package grpcUC

import (
	"context"
	proto "github.com/DoktorGhost/external-api/src/go/pkg/grpc/clients/api/grpc/protobuf/books_v1"
	"github.com/DoktorGhost/golibrary-books/internal/repositories/postgres/dao"
	usecases "github.com/DoktorGhost/golibrary-books/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookGRPCServer struct {
	uc *usecases.BookUseCase
	proto.BooksServiceServer
}

func NewBookGRPCServer(uc *usecases.BookUseCase) *BookGRPCServer {
	return &BookGRPCServer{uc: uc}
}

// / controllers
func (s *BookGRPCServer) AddBook(ctx context.Context, req *proto.AddBookRequest) (*proto.Response, error) {
	book := dao.BookTable{
		Title:    req.Title,
		AuthorID: int(req.AuthorId),
	}

	// Вызов метода юзкейса
	id, err := s.uc.AddBook(book)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add book: %v", err)
	}

	return &proto.Response{Id: int64(id)}, nil
}

func (s *BookGRPCServer) AddAuthor(ctx context.Context, req *proto.AddAuthorRequest) (*proto.Response, error) {
	id, err := s.uc.AddAuthor(req.Name, req.Surname, req.Patronymic)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add author: %v", err)
	}

	return &proto.Response{Id: int64(id)}, nil
}

func (s *BookGRPCServer) GetAllBookWithAuthor(ctx context.Context, req *proto.Empty) (*proto.BooksResponse, error) {
	books, err := s.uc.GetAllBookWithAuthor()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "not found: %v", err)
	}

	booksList := make([]*proto.BookWithAuthor, len(books))
	for i, b := range books {
		booksList[i] = &proto.BookWithAuthor{
			Id:    int64(b.ID),
			Title: b.Title,
			Author: &proto.Author{
				Id:       int64(b.Author.ID),
				FullName: b.Author.Name,
			},
		}
	}

	return &proto.BooksResponse{Books: booksList}, nil
}

func (s *BookGRPCServer) GetBookWithAuthor(ctx context.Context, req *proto.AuthorID) (*proto.BookWithAuthor, error) {
	book, err := s.uc.GetBookWithAuthor(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "book not found: %v", err)
	}

	return &proto.BookWithAuthor{
		Id:    int64(book.ID),
		Title: book.Title,
		Author: &proto.Author{
			Id:       int64(book.Author.ID),
			FullName: book.Author.Name,
		},
	}, nil
}

func (s *BookGRPCServer) GetAllAuthorWithBooks(ctx context.Context, req *proto.Empty) (*proto.AuthorResponse, error) {
	authors, err := s.uc.GetAllAuthorWithBooks()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "authors not found: %v", err)
	}

	authorList := make([]*proto.AuthorWithBooks, len(authors))
	for i, a := range authors {
		books := make([]*proto.Book, len(a.Books))
		for j, b := range a.Books {
			books[j] = &proto.Book{
				Id:       int64(b.ID),
				Title:    b.Title,
				AuthorId: int64(b.AuthorID),
			}
		}

		authorList[i] = &proto.AuthorWithBooks{
			Id:       int64(a.ID),
			FullName: a.Name,
			Books:    books,
		}
	}

	return &proto.AuthorResponse{Authors: authorList}, nil
}
