package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/logger"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/models"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/repo"
	"github.com/mercari/go-circuitbreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const otelName = "github.com/chernyshev-alex/go-bookstore-oapi/internal/service"

//go:generate counterfeiter -generate
//counterfeiter:generate -o test/books.gen.go . BooksService

type Book = models.Book
type Books = []*Book
type BooksService interface {
	AddBook(context.Context, Book) (*Book, error)
	DeleteBook(context.Context, string) error
	FindBooksByAuthor(context.Context, string) (Books, error)
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

func (s *BookService) FindBooksByAuthor(ctx context.Context, authorId string) (books Books, err error) {
	v, err := runAround(ctx, "Books.SearchByAuthor", s.cb,
		func() (*Books, error) {
			if books, err = s.bookSearch.BooksByAuthorId(ctx, strToInt64(authorId)); err != nil {
				logger.Error("BooksByAuthorId failed", err)
				return nil, err
			}
			return &books, err
		})
	return *v, err
}

func (s *BookService) AddBook(ctx context.Context, b Book) (book *Book, err error) {
	v, err := runAround(ctx, "Books.AddBook", s.cb,
		func() (*Book, error) {
			book, err := s.bookCrud.AddBook(ctx, b)
			if err != nil {
				return nil, err
			}
			return &book, nil
		})
	return v, err
}

func (s *BookService) DeleteBook(ctx context.Context, bookId string) (err error) {
	n, err := runAround(ctx, "Books.DeleteBook", s.cb,
		func() (*int64, error) {
			n, err := s.bookCrud.DeleteBook(ctx, strToInt64(bookId))
			if err != nil {
				return nil, err
			}
			return &n, nil
		})
	_ = n
	return err
}

func runAround[T any](ctx context.Context, spanName string, cb *circuitbreaker.CircuitBreaker,
	f func() (*T, error)) (*T, error) {

	defer newOTELSpan(ctx, spanName).End()
	var err error
	if !cb.Ready() {
		return nil, NewServiceError(ErrorServiceNotAvailable, "service not available")
	}
	defer func() {
		err = cb.Done(ctx, err)
	}()
	return f()
}

func newOTELSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)
	return span
}

func strToInt64(s string) (v int64) {
	var err error
	if v, err = strconv.ParseInt(s, 10, 64); err != nil {
		logger.Error("bad conversion str->i64", err)
	}
	return

}
