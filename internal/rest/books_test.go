package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/gen"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/service/test"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

const (
	CONTENT_TYPE_TAG = "Content-Type"
	JSON_APP_CONTENT = "application/json"
)

const baseURL = "http://127.0.0.1:8080"

func TestTasks_Search(t *testing.T) {
	t.Parallel()

	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

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
		setup  func(*test.FakeBooksService)
		output output
	}{{
		"OK: 200",
		func(s *test.FakeBooksService) { s.SearchBooksByAuthorReturns([]gen.Book{book}, nil) },
		output{
			http.StatusOK,
			&gen.SearchBooksResponse{book},
			&gen.SearchBooksResponse{},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fakeBooService := &test.FakeBooksService{}
			tt.setup(fakeBooService)

			handler := NewBooksHandler(fakeBooService)

			// params := rest.SearchBooksByAuthorParams{
			// 	AuthorId: 100000,
			// }

			rq, err := http.NewRequest("GET", "/search/books?authorId=1", nil)

			//	rq, err := rest.NewSearchBooksByAuthorRequest(baseURL, &params)
			//rq, err := rest.SearchBooksByAuthor(params)
			if err != nil {
				t.Errorf("request error %v", err.Error())
			}

			resp := sendRequest(rq, handler)

			assert.Equal(t, tt.output.expectedStatus, resp.StatusCode)
			assert.Equal(t, JSON_APP_CONTENT, resp.Header.Get(CONTENT_TYPE_TAG))
			assertResponse(t, resp, tt.output.expected, tt.output.target)
		})
	}
}

func sendRequest(req *http.Request, spi gen.ServerInterface) *http.Response {
	w := httptest.NewRecorder()
	ctx, e := gin.CreateTestContext(w)
	ctx.Request = req

	gen.RegisterHandlers(e, spi)
	e.ServeHTTP(w, req)
	return w.Result()
}

func assertResponse(t *testing.T, res *http.Response, expected interface{}, target interface{}) {
	t.Helper()

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		t.Fatalf("couldn't decode %s", err)
	}

	if !cmp.Equal(expected, target) {
		t.Fatalf("expected results don't match: %s", cmp.Diff(expected, target))
	}
}
