package repo

import (
	"context"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/gen"
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

var (
	_ BooksCrudRepository   = (NewRepository)(nil)
	_ BooksSearchRepository = (NewRepository)(nil)
)

func (r *DbaRepository) AddBook(ctx context.Context, book gen.Book) (gen.Book, error) {
	return book, nil
}

func (r *DbaRepository) DeleteBook(ctx context.Context, bookId int) error {
	return nil
}

func (r *DbaRepository) SearchByAuthor(ctx context.Context, params gen.SearchBooksByAuthorParams) ([]gen.Book, error) {
	// db save
	// ...
	return []gen.Book{{
		Author:      "author",
		AuthorId:    0,
		BookId:      0,
		Publisher:   "publisher",
		PublisherId: 0,
		Title:       "title from db",
		Year:        1999,
	}}, nil
}
