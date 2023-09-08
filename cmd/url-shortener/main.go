package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/alexzhirkov/url-shortener/internal/config"
	http_handler "github.com/alexzhirkov/url-shortener/internal/handlers/http"
	http_server "github.com/alexzhirkov/url-shortener/internal/lib/http-server"
	"github.com/alexzhirkov/url-shortener/internal/lib/logger/sl"
	"github.com/alexzhirkov/url-shortener/internal/repository/url_repository/sqlite"
	"github.com/alexzhirkov/url-shortener/internal/usecases"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {

	var r io.Reader
	r = os.Stdin
	r = bufio.NewReader(r)
	//r = new(bytes.Buffer)

	//init.go config
	cfg := config.MustLoad()

	//init.go logger
	logger := sl.SetupLogger(cfg.Env)

	//logger.With(slog.String("env", cfg.Env)).Info("info message")
	logger.Info("start http-shortener", slog.String("env", cfg.Env))
	//logger.Debug("debug message")

	//adapters: sqlite
	urlStore, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("adapters initialization failed", sl.Err(err))
		return fmt.Errorf("adapters initialization failed: %v", err)
	}
	useCase, err := usecases.New(usecases.WithUrlRepository(urlStore))
	//useCase, err := usecases.New(usecases.WithUrlRepository(in_memory.New()))
	//useCase, err := usecases.New(usecases.WithMemoryRepository())
	if err != nil {
		return fmt.Errorf("usecases initialization failed: %v", err)
	}
	HTTPHandler := http_handler.NewHTTPHandler(useCase)

	//init.go router: gin
	gin.SetMode(gin.ReleaseMode) //switch off warnings
	//router := gin.New()	//without recoverer and logger
	router := gin.Default()

	//middleware
	//router.Use()
	router.Use(requestid.New())

	//router.Use(gin.BasicAuth(gin.Accounts{"alex": "pass"}))//авторизация. встроенная только basic
	//выключим пока, чтобы смотреть в браузере
	//router.Use(middleware.EnsureLoggedIn()) //Bearer авторизация
	//router.Use(middleware.Logger)	//потом

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "id:"+requestid.Get(c))
	})
	router.GET("/get/:alias", HTTPHandler.GetUrl)
	router.POST("/create", HTTPHandler.CreateUrl)
	router.GET("/count", HTTPHandler.CountUrls)
	//run server
	//router.Run(cfg.Address)

	// Create the HTTP server
	httpServer := http_server.NewHttpServer(
		router,
		cfg.HttpServer,
	)

	// Start the HTTP server
	httpServer.Start()
	defer httpServer.Stop()
	// Listen for OS signals to perform a graceful shutdown
	//c := make(chan os.Signal, 1)
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	_ = cancel //call it to initiate shutdown
	log.Println("listening signals...")
	<-ctx.Done()
	log.Println("graceful shutdown...")

	return nil
}
