package app

import (
	"context"
	"go-shopping-cart/internal/config"
	"go-shopping-cart/internal/db"
	"go-shopping-cart/internal/db/sqlc"
	"go-shopping-cart/internal/routes"
	"go-shopping-cart/internal/validation"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

type ModuleContext struct {
	DB sqlc.Querier
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Validator init failed %v", err)
	}

	loadEnv()

	r := gin.Default()

	if err := db.InitDB(); err != nil {
		log.Fatalf("unable to connnect to db: %v", err)
	}

	ctx := &ModuleContext{
		DB: db.DB,
	}

	modules := []Module{
		NewUserModule(ctx),
	}

	routes.RegisterRoutes(r, getModulRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
	}

	//Graceful Shutdown for API
	quit := make(chan os.Signal, 1)
	//syscall.SIGINT -> Ctrl+C
	//syscall.SIGTERM -> kill service
	//syscall.SIGHUP -> Reload service
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		log.Printf("Server is running at %s", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutdown signal received ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")

	return nil
}

func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}

func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found")
	}
}
