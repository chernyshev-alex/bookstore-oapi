package rest

// Generate types and interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/gen"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/service"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	booksSpi service.BooksService
}

var _ gen.ServerInterface = (NewBooksHandler)(nil)

func NewBooksHandler(spi service.BooksService) *BookHandler {
	return &BookHandler{booksSpi: spi}
}

// -- ServerInterface --

func (h *BookHandler) AddBook(c *gin.Context) {
	var book gen.Book

	if err := json.NewDecoder(c.Request.Body).Decode(&book); err != nil {
		httpError(c.Writer, err, http.StatusBadRequest)
		return
	}
	defer c.Request.Body.Close()

	b, err := h.booksSpi.AddBook(c.Request.Context(), book)
	if err != nil {
		httpError(c.Writer, err, http.StatusInternalServerError)
		return
	}
	responseJson(c.Writer, http.StatusOK, b)
}

func (h *BookHandler) DeleteBook(c *gin.Context, bookId int) {
	if err := h.booksSpi.DeleteBook(c.Request.Context(), bookId); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	responseJson(c.Writer, http.StatusOK, nil)
}

func (h *BookHandler) SearchBooksByAuthor(c *gin.Context, params gen.SearchBooksByAuthorParams) {
	books, err := h.booksSpi.SearchBooksByAuthor(c.Request.Context(), params)
	if err != nil {
		httpError(c.Writer, err, http.StatusInternalServerError)
		return
	}
	responseJson(c.Writer, http.StatusOK, books)
}

// TODO move to httpUtils

func httpError(w http.ResponseWriter, e error, httpStatus int) {
	http.Error(w, e.Error(), httpStatus)
}

func responseJson(w http.ResponseWriter, httpStatus int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(body)
}
