package srv

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/rest"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////go:embed src/books.yaml
//var testSchema []byte

func doGet(t *testing.T, handler http.Handler, rawURL string) *httptest.ResponseRecorder {
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("Invalid url: %s", rawURL)
	}

	response := testutil.NewRequest().Get(u.RequestURI()).WithHost(u.Host).WithAcceptJson().GoWithHTTPHandler(t, handler)
	return response.Recorder
}

func doPost(t *testing.T, handler http.Handler, rawURL string, jsonBody interface{}) *httptest.ResponseRecorder {
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("Invalid url: %s", rawURL)
	}

	response := testutil.NewRequest().Post(u.RequestURI()).WithHost(u.Host).WithJsonBody(jsonBody).GoWithHTTPHandler(t, handler)
	return response.Recorder
}

var baseUrl = "http://127.0.0.1:8080"

func Test_OpenApiRequestValidator(t *testing.T) {
	swagger, err := openapi3.NewLoader().LoadFromFile("../../books.yaml")
	require.NoError(t, err, "Error initializing swagger")

	validator, _ := rest.NewValidator()
	handler := middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: rest.NewAuthenticator(validator),
			},
		})

	g := gin.New()
	g.Use(handler)

	// create fake token
	tokenStr := rest.CreateToken(jwt.MapClaims{"foo": "bar"})

	rec := testutil.NewRequest().WithAcceptJson().WithJWSAuth(tokenStr).
		Get(baseUrl+"/search/books?authorId=1000").
		GoWithHTTPHandler(t, g)

	assert.Equal(t, http.StatusForbidden, rec.Code())

}
