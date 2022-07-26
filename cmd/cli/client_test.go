package cli

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

const (
	baseURL string = "http://localhost:8080"
)

func Test_AddSearchBook(t *testing.T) {
	book := Book{
		Author:      "A. Conan Doyle",
		AuthorId:    1,
		BookId:      10,
		Publisher:   "George Newnes",
		PublisherId: 1000,
		Title:       "The Adventures of Sherlock Holmes",
		Year:        0,
	}

	ctx := context.Background()

	c, _ := NewClientWithResponses(baseURL)
	resp, err := c.AddBookWithResponse(ctx, book)

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode(), http.StatusCreated)

	authorId := resp.JSON201.Book.AuthorId

	searchResp, err := c.SearchBooksByAuthorWithResponse(ctx, &SearchBooksByAuthorParams{AuthorId: authorId})

	assert.NoError(t, err)
	assert.Equal(t, searchResp.StatusCode(), http.StatusOK)
	assert.Equal(t, true, cmp.Equal([]Book{book}, *(searchResp.JSON200)))
}
