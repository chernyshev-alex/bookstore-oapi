package repo

import (
	"context"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/gen"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/logger"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

//go:generate counterfeiter -generate
//counterfeiter:generate -o test/books-search.gen.go . BooksSearchRepository
//counterfeiter:generate -o test/books-crud.gen.go . BooksCrudRepository

type BooksSearchRepository interface {
	SearchByAuthor(context.Context, gen.SearchBooksByAuthorParams) ([]gen.Book, error)
}

type BooksCrudRepository interface {
	AddBook(context.Context, gen.Book) (gen.Book, error)
	DeleteBook(context.Context, int) error
}

type DbaRepository struct {
	xorm *xorm.Engine
}

func NewRepository(x *xorm.Engine) *DbaRepository {
	return &DbaRepository{
		xorm: x,
	}
}

const (
	TBL_BOOKS = "BOOKS"
)

var (
	_ BooksCrudRepository   = (NewRepository)(nil)
	_ BooksSearchRepository = (NewRepository)(nil)
)

// @see https://github.com/zzdboy/GoCMS/blob/master/app/models/article.go

func (r *DbaRepository) AddBook(ctx context.Context, book gen.Book) (gen.Book, error) {
	_, err := r.xorm.Insert(book)
	if err != nil {
		logger.Error("AddBook failed", err)
		return gen.Book{}, err
	}
	logger.Debug("AddBook OK", zap.Any("book", book))
	return book, nil
}

// TODO : bookId should be string
// TODO : use Fluent SQL
func (r *DbaRepository) DeleteBook(ctx context.Context, bookId int) error {
	_, err := r.xorm.Table(TBL_BOOKS).Delete(bookId)
	if err != nil {
		logger.Error("deleteBook failed", err)
		return err
	}
	return nil
}

// TODO : remove dependency SearchBooksByAuthorParams
func (r *DbaRepository) SearchByAuthor(ctx context.Context, params gen.SearchBooksByAuthorParams) ([]gen.Book, error) {
	authorId := params.AuthorId
	books := []gen.Book{}
	err := r.xorm.Find(books, authorId)
	if err != nil {
		logger.Error("SearchByAuthor failed", err)
		return nil, err
	}
	logger.Info("SearchByAuthor OK", zap.Any("[]books", books))
	return books, nil
}
