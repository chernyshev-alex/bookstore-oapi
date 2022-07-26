package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/gen"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/repo"
	"github.com/mercari/go-circuitbreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const otelName = "github.com/chernyshev-alex/go-bookstore-oapi/internal/service"

//go:generate counterfeiter -generate
//counterfeiter:generate -o test/books.gen.go . BooksService

type BooksService interface {
	AddBook(context.Context, gen.AddBookRequest) (gen.Book, error)
	DeleteBook(context.Context, int) error
	SearchBooksByAuthor(context.Context, gen.SearchBooksByAuthorParams) ([]gen.Book, error)
}

// Service error --
type ErrorCode int
type ServiceError struct {
	orig error
	code ErrorCode
}

func NewServiceError(errCode ErrorCode, msg string) ServiceError {
	return ServiceError{
		orig: errors.New(msg),
		code: errCode,
	}
}

// Error implements error
func (se ServiceError) Error() string {
	return se.orig.Error()
}

func (se ServiceError) ErrorCode() ErrorCode {
	return se.code
}

const (
	ErrorCodeGeneral ErrorCode = iota
	ErrorCodeNotFound
	ErrorServiceNotAvailable
)

type BookService struct {
	cb         *circuitbreaker.CircuitBreaker
	bookSearch repo.BooksSearchRepository
	bookCrud   repo.BooksCrudRepository
}

func NewBookService(crud repo.BooksCrudRepository, search repo.BooksSearchRepository) *BookService {
	return &BookService{
		bookSearch: search,
		bookCrud:   crud,
		cb: circuitbreaker.New(
			circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncConsecutiveFailures(3))),
	}
}

func (s *BookService) SearchBooksByAuthor(ctx context.Context, params gen.SearchBooksByAuthorParams) ([]gen.Book, error) {
	defer newOTELSpan(ctx, "Books.SearchByAuthor").End()
	if !s.cb.Ready() {
		return []gen.Book{}, NewServiceError(ErrorServiceNotAvailable, "service not available")
	}
	var err error
	defer func() {
		err = s.cb.Done(ctx, err)
	}()
	books, err := s.bookSearch.SearchByAuthor(ctx, params)
	if err != nil {
		return []gen.Book{}, err
	}
	return books, nil
}

func (s *BookService) AddBook(ctx context.Context, b gen.Book) (book gen.Book, err error) {
	defer newOTELSpan(ctx, "Book.Create").End()
	if !s.cb.Ready() {
		return gen.Book{}, fmt.Errorf("service not available")
	}
	defer func() {
		err = s.cb.Done(ctx, err)
	}()

	book, err = s.bookCrud.AddBook(ctx, b)
	if err != nil {
		return gen.Book{}, err
	}
	return book, nil
}

func (s *BookService) DeleteBook(ctx context.Context, bookId int) (err error) {
	defer newOTELSpan(ctx, "Book.Remove").End()
	if !s.cb.Ready() {
		return fmt.Errorf("service not available")
	}

	defer func() {
		err = s.cb.Done(ctx, err)
	}()

	err = s.bookCrud.DeleteBook(ctx, bookId)
	return
}

func newOTELSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)
	return span
}
