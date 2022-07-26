package rest

// Generate types and interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/gen"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/service"
	"github.com/gin-gonic/gin"
)

/*
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=conf/server-gen.conf  ../../books.yaml
//go:generate counterfeiter -generate
//counterfeiter:generate -o test/books-service.gen.go . BookService
*/

/*
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=conf/client-gen.conf  ../../books.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=conf/types-gen.conf   ../../books.yaml
*/

// type BooksService interface {
// 	AddBook(context.Context, AddBookRequest) (Book, error)
// 	DeleteBook(context.Context, int) error
// 	SearchBooksByAuthor(context.Context, SearchBooksByAuthorParams) ([]Book, error)
// }

type BookHandler struct {
	booksSpi service.BooksService
}

var _ gen.ServerInterface = (NewBooksHandler)(nil)

func NewBooksHandler(spi service.BooksService) *BookHandler {
	return &BookHandler{booksSpi: spi}
}

// AddBook implements ServerInterface
func (h *BookHandler) AddBook(c *gin.Context) {
	var book gen.Book
	if err := json.NewDecoder(c.Request.Body).Decode(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer c.Request.Body.Close()

	responseJson(c.Writer, http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c *gin.Context, bookId int) {
	h.booksSpi.DeleteBook(c.Request.Context(), bookId)
}

func (h *BookHandler) SearchBooksByAuthor(c *gin.Context, params gen.SearchBooksByAuthorParams) {
	books, err := h.booksSpi.SearchBooksByAuthor(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	responseJson(c.Writer, http.StatusOK, books)
}

// TODO move to httpUtils
func responseJson(w http.ResponseWriter, httpStatus int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(body)
}
