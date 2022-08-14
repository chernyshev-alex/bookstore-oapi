package cli

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

const (
	baseURL string = "http://localhost:8080"
)

func Test_FindBooksByAuthor(t *testing.T) {
	book := BookJson{
		Author:    "A. Conan Doyle",
		AuthorId:  1,
		Publisher: "George Newnes",
		//	PublisherId: 1000,
		Title: "The Adventures of Sherlock Holmes",
		Year:  0,
	}

	ctx := context.Background()

	c, _ := NewClientWithResponses(baseURL)
	authorId := "1000"

	booksOfAuthor, err := c.BooksByAuthorIdWithResponse(ctx, &BooksByAuthorIdParams{AuthorId: authorId})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, booksOfAuthor.StatusCode())
	assert.Equal(t, true, cmp.Equal([]BookJson{book}, *(booksOfAuthor.JSON200)))
}

func Test_AddSearchBook(t *testing.T) {
	book := BookJson{
		Author:    "A. Conan Doyle",
		AuthorId:  1,
		Publisher: "George Newnes",
		//	PublisherId: 1000,
		Title: "The Adventures of Sherlock Holmes",
		Year:  0,
	}

	ctx := context.Background()

	c, _ := NewClientWithResponses(baseURL)
	resp, err := c.AddBookWithResponse(ctx, book)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())

	authorId := strconv.Itoa(resp.JSON201.Book.AuthorId)

	booksOfAuthor, err := c.BooksByAuthorIdWithResponse(ctx, &BooksByAuthorIdParams{AuthorId: authorId})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, booksOfAuthor.StatusCode())
	assert.Equal(t, true, cmp.Equal([]BookJson{book}, *(booksOfAuthor.JSON200)))
}
