package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/env"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/gen"
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
	"xorm.io/xorm"
)

type serverConfig struct {
	env *env.EnvConfig
	db  *xorm.Engine
	//	Memcached     *memcache.Client
	metrics     http.Handler
	middlewares gin.HandlersChain
	logger      *zap.Logger
	handler     gen.ServerInterface
}

func NewServer(conf *serverConfig) *http.Server {
	r := gin.Default()
	for _, mw := range conf.middlewares {
		r.Use(mw)
	}

	r = gen.RegisterHandlers(r, conf.handler)
	registerSwaggerApi(conf.logger, r)

	lmt := tollbooth.NewLimiter(conf.env.MAX_LIMITER, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})
	lmtmw := tollbooth.LimitHandler(lmt, r)

	conf.logger.Info("starting server on", zap.String("address", conf.env.HTTP_ADDRESS))
	return &http.Server{
		Handler:           lmtmw,
		Addr:              conf.env.HTTP_ADDRESS,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
}

func initDb(env *env.EnvConfig) *xorm.Engine {
	db_zap_info := zap.Strings("xorm", []string{env.DB_DRIVER, env.DATABASE_NAME, env.DATABASE_USERNAME})

	engine, err := xorm.NewEngine(env.DB_DRIVER, env.DATABASE_NAME)
	if err != nil {
		log.Panic("db connect failed", db_zap_info, zap.Error(err))
	}
	err = repo.MayBeMigrate(engine)
	if err != nil {
		log.Panic("db migraton failed", db_zap_info, zap.Error(err))
	}
	return engine
}

func registerSwaggerApi(logger *zap.Logger, router *gin.Engine) {
	swagger, err := api.GetSwagger()
	if err != nil {
		logger.Error("failed loading swagger file ", zap.Error(err))
	}

	router.GET("/openapi3.json", func(c *gin.Context) {
		data, _ := yaml.Marshal(&swagger)
		w := c.Writer
		_, _ = w.Write(data)
		c.JSON(http.StatusOK, gin.H{"Content-Type": "application/x-yaml"})
	})
}

func startServices(env *env.EnvConfig) {
	logger, _ := zap.NewProduction()

	swagger, err := api.GetSwagger()
	if err != nil {
		logger.Panic("Error loading swagger spec", zap.Error(err))
	}
	swagger.Servers = nil

	logger.Info("init databases")
	db_engine := initDb(env)

	logger.Info("init repositories")
	repo := repo.NewRepository(db_engine)

	logger.Info("init services")
	svc := service.NewBookService(repo, repo)

	logger.Info("init controllers")
	bookStoreAPI := rest.NewBooksHandler(svc)

	logger.Info("init telemetry")
	promExporter, err := otel.NewOTExporter(env)

	if err != nil {
		logger.Panic("telemetry failed", zap.Error(err))
	}

	// log all requests
	var logRequests gin.HandlerFunc = func(g *gin.Context) {
		logger.Info(g.Request.Method, zap.Time("time", time.Now()),
			zap.String("url", g.Request.URL.String()),
		)
	}

	srv_conf := serverConfig{
		env:     env,
		logger:  logger,
		db:      db_engine,
		metrics: promExporter,
		handler: bookStoreAPI,
		middlewares: gin.HandlersChain{
			logRequests,
			middleware.OapiRequestValidator(swagger),
		},
	}

	err = NewServer(&srv_conf).ListenAndServe()
	if err != http.ErrServerClosed {
		logger.Fatal("http server failed", zap.Error(err))
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
		env.HTTP_ADDRESS = address
	}
	startServices(env)
}
