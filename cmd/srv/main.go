package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/env"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/logger"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/otel"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/repo"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/rest"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/service"
	"github.com/deepmap/oapi-codegen/examples/petstore-expanded/gin/api"
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	_ "github.com/mattn/go-sqlite3"
)

type serverConfig struct {
	env *env.Config
	//	Memcached     *memcache.Client
	metrics     http.Handler
	middlewares gin.HandlersChain
	handler     rest.ServerInterface
}

func NewServer(conf *serverConfig) *http.Server {
	r := gin.Default()
	for _, mw := range conf.middlewares {
		r.Use(mw)
	}

	r = rest.RegisterHandlers(r, conf.handler)
	registerSwaggerApi(r)

	lmt := tollbooth.NewLimiter(conf.env.Server.HttpLimiter,
		&limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})
	lmtmw := tollbooth.LimitHandler(lmt, r)

	address := fmt.Sprintf("%s:%d", conf.env.Server.Host, conf.env.Server.Port)
	logger.Info("starting server on", zap.String("address", address))

	return &http.Server{
		Handler:           lmtmw,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
}

func initDb(conf *env.DbConfig) *sql.DB {
	logger.Info("starting  db  : ", zap.Any("conf", conf))

	db, err := sql.Open(conf.Driver, conf.Dbname)
	if err != nil {
		logger.Panic("db connect failed", zap.Error(err))
	}
	return db
}

func registerSwaggerApi(router *gin.Engine) {
	swagger, err := api.GetSwagger()
	if err != nil {
		logger.Error("failed loading swagger file ", err)
	}

	router.GET("/openapi3.json", func(c *gin.Context) {
		data, _ := yaml.Marshal(&swagger)
		w := c.Writer
		_, _ = w.Write(data)
		c.JSON(http.StatusOK, gin.H{"Content-Type": "application/x-yaml"})
	})
}

func startServices(conf *env.Config) {
	swagger, err := api.GetSwagger()
	if err != nil {
		logger.Panic("Error loading swagger spec", zap.Error(err))
	}
	swagger.Servers = nil

	logger.Info("init db", zap.Any("conf", conf.Db))
	db := initDb(&conf.Db)

	logger.Info("init repositories")
	repo := repo.NewRepository(db)

	logger.Info("init services")
	svc := service.NewBookService(repo, repo)

	logger.Info("init controllers")
	bookStoreAPI := rest.NewBooksHandler(svc)

	logger.Info("init telemetry")
	promExporter, err := otel.NewOTExporter(&conf.Server)
	if err != nil {
		logger.Panic("telemetry failed", zap.Error(err))
	}

	// TODO : log all requests
	var logRequests gin.HandlerFunc = func(g *gin.Context) {
		logger.Info(g.Request.Method, zap.Time("time", time.Now()),
			zap.String("url", g.Request.URL.String()),
		)
	}

	srv_conf := serverConfig{
		env:     conf,
		metrics: promExporter,
		handler: bookStoreAPI,
		middlewares: gin.HandlersChain{
			logRequests,
			middleware.OapiRequestValidator(swagger),
		},
	}

	err = NewServer(&srv_conf).ListenAndServe()
	if err != http.ErrServerClosed {
		logger.Panic("http server failed", zap.Error(err))
	}
}

func main() {
	var confPath, address string

	flag.StringVar(&confPath, "conf", ".env", "config file")
	flag.StringVar(&address, "address", "", "HTTP Server Address")
	flag.Parse()

	env, err := env.LoadConfig(confPath)
	if err != nil {
		log.Fatal("failed load config")
	}

	if len(address) > 0 { // override
		env.Server.Host = address
	}
	startServices(env)
}
