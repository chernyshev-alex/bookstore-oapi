package service

import (
	"context"
	"errors"
	"testing"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/gen"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/repo/test"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type output struct {
	expected interface{}
	target   interface{}
}

// After N consecutive failures, circuit breaker should returns Service Not Available error
func Test_SearchWithCircuitBreaker(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*test.FakeBooksSearchRepository)
		output output
	}{{"search by author",
		func(s *test.FakeBooksSearchRepository) { s.SearchByAuthorReturns(nil, errors.New("failed")) },
		output{nil, &gen.SearchBooksResponse{}},
	}}

	triggeredServiceNotAvailable := false
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeSearchApi := &test.FakeBooksSearchRepository{}
			tt.setup(fakeSearchApi)

			svc := NewBookService(nil, fakeSearchApi)
			for i := 0; i < 10; i++ {
				_, err := svc.SearchBooksByAuthor(context.Background(), gen.SearchBooksByAuthorParams{
					AuthorId: 1,
				})

				if err != nil && errors.As(err, &ServiceError{}) {
					if err.(ServiceError).ErrorCode() == ErrorServiceNotAvailable {
						triggeredServiceNotAvailable = true
					}
					break
				}
			}
		})
		assert.Equal(t, true, triggeredServiceNotAvailable)
	}
}

func Test_SearchBook(t *testing.T) {
	t.Parallel()

	books := []gen.Book{{
		Author:      "Author",
		AuthorId:    0,
		BookId:      0,
		Publisher:   "Publisher",
		PublisherId: 0,
		Title:       "Title",
		Year:        1901,
	}}

	tests := []struct {
		name   string
		setup  func(*test.FakeBooksSearchRepository)
		output output
	}{{
		"search by author",
		func(s *test.FakeBooksSearchRepository) { s.SearchByAuthorReturns(books, nil) },
		output{books, &gen.SearchBooksResponse{}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fakeSearchApi := &test.FakeBooksSearchRepository{}
			tt.setup(fakeSearchApi)

			svc := NewBookService(nil, fakeSearchApi)
			books, err := svc.SearchBooksByAuthor(context.Background(), gen.SearchBooksByAuthorParams{
				AuthorId: 1,
			})
			if err != nil {
				t.Fatalf("error %v", err.Error())
			}
			if !cmp.Equal(books, tt.output.expected) {
				t.Fatalf("don't match: %s", cmp.Diff(books, tt.output.expected))
			}
		})
	}
}

func Test_CreateBook(t *testing.T) {
	t.Parallel()

	book := gen.Book{
		Author:      "Author",
		AuthorId:    0,
		BookId:      0,
		Publisher:   "Publisher",
		PublisherId: 0,
		Title:       "Title",
		Year:        1901,
	}

	tests := []struct {
		name   string
		setup  func(*test.FakeBooksCrudRepository)
		output output
	}{{"create book",
		func(s *test.FakeBooksCrudRepository) { s.AddBookReturns(book, nil) },
		output{book, nil},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fakeCrudApi := &test.FakeBooksCrudRepository{}

			tt.setup(fakeCrudApi)

			svc := NewBookService(fakeCrudApi, nil)
			createdBook, err := svc.AddBook(context.Background(), book)
			if err != nil {
				t.Fatalf("error %v", err.Error())
			}
			if !cmp.Equal(createdBook, tt.output.expected) {
				t.Fatalf("don't match: %s", cmp.Diff(createdBook, tt.output.expected))
			}
		})
	}
}

func Test_RemoveBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(*test.FakeBooksCrudRepository)
		expected interface{}
	}{{
		"delete book",
		func(s *test.FakeBooksCrudRepository) { s.DeleteBookReturns(nil) },
		nil,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fakeCrudApi := &test.FakeBooksCrudRepository{}
			tt.setup(fakeCrudApi)

			svc := NewBookService(fakeCrudApi, nil)
			err := svc.DeleteBook(context.Background(), 1)
			if err != nil {
				t.Fatalf("error %v", err.Error())
			}
			assert.Equal(t, tt.expected, err)
		})
	}
}
