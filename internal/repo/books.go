package repo

import (
	"context"
	"database/sql"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/logger"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

//go:generate counterfeiter -generate
//counterfeiter:generate -o test/books-search.gen.go . BooksSearchRepository
//counterfeiter:generate -o test/books-crud.gen.go . BooksCrudRepository

type BooksSearchRepository interface {
	BooksByAuthorId(context.Context, int64) ([]*models.Book, error)
}

type BooksCrudRepository interface {
	AddBook(context.Context, models.Book) (models.Book, error)
	DeleteBook(context.Context, int64) (int64, error)
}

type DbaRepository struct {
}

func NewRepository(x *sql.DB) *DbaRepository {
	boil.SetDB(x) // set as global
	return &DbaRepository{}
}

var (
	_ BooksCrudRepository   = (NewRepository)(nil)
	_ BooksSearchRepository = (NewRepository)(nil)
)

func (r *DbaRepository) AddBook(ctx context.Context, book models.Book) (models.Book, error) {
	err := book.InsertG(ctx, boil.Infer())
	return book, err
}

func (r *DbaRepository) DeleteBook(ctx context.Context, Id int64) (int64, error) {
	var book models.Book = models.Book{
		ID: null.Int64From(Id),
	}
	return book.DeleteG(ctx)
}

func (r *DbaRepository) BooksByAuthorId(ctx context.Context, authorId int64) (booksResult []*models.Book, e error) {
	books, err := models.Books(models.BookWhere.Authorid.EQ(int64(authorId))).AllG(ctx)
	if err != nil {
		logger.Error("SearchByAuthor failed", err)
		return nil, err
	}
	return books, nil
}
