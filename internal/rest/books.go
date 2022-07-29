package rest

// Generate types and interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/service"
	"github.com/gin-gonic/gin"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server-gen.conf  ../../books.yaml

type BookHandler struct {
	booksSpi service.BooksService
}

var _ ServerInterface = (NewBooksHandler)(nil)

func NewBooksHandler(spi service.BooksService) *BookHandler {
	return &BookHandler{booksSpi: spi}
}

// -- ServerInterface --

func (h *BookHandler) AddBook(c *gin.Context) {
	var rq AddBookRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&rq); err != nil {
		httpError(c.Writer, err, http.StatusBadRequest)
		return
	}
	defer c.Request.Body.Close()

	b, err := h.booksSpi.AddBook(c.Request.Context(), addBookRequestToDomain(&rq))
	if err != nil {
		httpError(c.Writer, err, http.StatusInternalServerError)
		return
	}
	responseJson(c.Writer, http.StatusOK, AddBookResponse{Book: domainToJsonBook(&b)})
}

func (h *BookHandler) DeleteBook(c *gin.Context, bookId BookId) {
	if err := h.booksSpi.DeleteBook(c.Request.Context(), bookId); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	responseJson(c.Writer, http.StatusOK, nil)
}

func (h *BookHandler) BooksByAuthorId(c *gin.Context, params BooksByAuthorIdParams) {
	books, err := h.booksSpi.FindBooksByAuthor(c.Request.Context(), params.AuthorId)
	if err != nil {
		httpError(c.Writer, err, http.StatusInternalServerError)
		return
	}
	responseJson(c.Writer, http.StatusOK, domainSliceToJsonBooks(books))
}

func httpError(w http.ResponseWriter, e error, httpStatus int) {
	http.Error(w, e.Error(), httpStatus)
}

// helpers
func responseJson(w http.ResponseWriter, httpStatus int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(body)
}
