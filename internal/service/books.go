package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/logger"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/repo"
	"github.com/chernyshev-alex/go-bookstore-oapi/pkg/domain"
	"github.com/mercari/go-circuitbreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const otelName = "github.com/chernyshev-alex/go-bookstore-oapi/internal/service"

//go:generate counterfeiter -generate
//counterfeiter:generate -o test/books.gen.go . BooksService

type BooksService interface {
	AddBook(context.Context, domain.Book) (domain.Book, error)
	DeleteBook(context.Context, string) error
	FindBooksByAuthor(context.Context, string) (domain.Books, error)
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

func (s *BookService) FindBooksByAuthor(ctx context.Context, authorId string) ([]domain.Book, error) {
	defer newOTELSpan(ctx, "Books.SearchByAuthor").End()
	if !s.cb.Ready() {
		return []domain.Book{}, NewServiceError(ErrorServiceNotAvailable, "service not available")
	}

	var (
		authorIdKey int
		err         error
		books       domain.Books = []domain.Book{}
	)

	defer func() {
		err = s.cb.Done(ctx, err)
	}()

	if authorIdKey, err = strconv.Atoi(authorId); err != nil {
		logger.Error("illegal value", err, zap.Int("authorId", authorIdKey))
		return books, err
	}
	if books, err = s.bookSearch.BooksByAuthorId(ctx, authorIdKey); err != nil {
		logger.Error("BooksByAuthorId failed", err)
		return books, err
	}
	return books, err
}

func (s *BookService) AddBook(ctx context.Context, b domain.Book) (book domain.Book, err error) {
	defer newOTELSpan(ctx, "Book.Create").End()
	if !s.cb.Ready() {
		return domain.Book{}, fmt.Errorf("service not available")
	}
	defer func() {
		err = s.cb.Done(ctx, err)
	}()

	book, err = s.bookCrud.AddBook(ctx, b)
	if err != nil {
		return domain.Book{}, err
	}
	return book, nil
}

func (s *BookService) DeleteBook(ctx context.Context, bookId string) (err error) {
	defer newOTELSpan(ctx, "Book.Remove").End()
	if !s.cb.Ready() {
		return fmt.Errorf("service not available")
	}

	defer func() {
		err = s.cb.Done(ctx, err)
	}()

	var bookIdKey int
	if bookIdKey, err = strconv.Atoi(bookId); err == nil {
		err = s.bookCrud.DeleteBook(ctx, bookIdKey)
	}
	return err
}

func newOTELSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)
	return span
}
