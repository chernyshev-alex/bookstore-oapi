package srv

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/rest"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

var (
	swagger_path = "../../books.yaml"
	swagger      *openapi3.T
	baseUrl      = "http://127.0.0.1:8080"
	testToken    string
	g            *gin.Engine
	r            gin.IRoutes
)

func init() {
	var err error
	swagger, err = openapi3.NewLoader().LoadFromFile(swagger_path)
	if err != nil {
		panic(err)
	}
	testToken = rest.CreateToken(jwt.MapClaims{"foo": "bar"})

	validator, _ := rest.NewValidator()
	handler := middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: rest.NewAuthenticator(validator),
			},
		})

	g = gin.New()
	r = g.Use(handler)
}

func Test_OpenApiRequestValidator(t *testing.T) {
	var q = "/search/books"
	var qbyAuthor = baseUrl + q + "?authorId=1000"

	r.GET(q, func(c *gin.Context) {})

	rec := testutil.NewRequest().WithAcceptJson().WithJWSAuth(testToken).Get(qbyAuthor).GoWithHTTPHandler(t, g)
	assert.Equal(t, http.StatusOK, rec.Code())

	rec = testutil.NewRequest().WithAcceptJson().Get(qbyAuthor).GoWithHTTPHandler(t, g)
	assert.Equal(t, http.StatusBadRequest, rec.Code())

}
