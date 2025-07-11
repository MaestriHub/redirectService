package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"redirectServer/configs"
	_ "redirectServer/docs"
	"redirectServer/internal/database"
	"redirectServer/internal/database/migrations"
	"redirectServer/internal/database/repo"
	"redirectServer/internal/service"
	"redirectServer/internal/transport/rest"
	"redirectServer/middlewares"
)

//	@title			Invite API
//	@version		1.0
//	@description	API for generating links and catch fingerprints.
//	@host			localhost:8080
//  @BasePath       /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := configs.NewConfig()

	db := database.InitDB(cfg.DbConfig)

	migrate := migrations.NewMigration1(db)
	if err := migrate.Up(); err != nil {
		log.Fatal(err)
	}

	app := gin.New()

	setupMiddlewares(app)

	// Repositories
	linkRepo := repository.NewLinkRepo(db)
	salonInfoRepo := repository.NewSalonInfoRepo(db)
	fingerprintRepo := repository.NewFingerprintRepo(db)

	// Services
	urlGenerator := service.NewUrlGeneratorService(cfg.UrlConfig.RedirectLink)
	linkService := service.NewLinkService(linkRepo)
	renderService := service.NewRenderService(cfg.AppStoreLinksConfig, salonInfoRepo, urlGenerator)
	fingerprintService := service.NewFingerprintService(fingerprintRepo, linkRepo)

	// Handlers
	linkHandler := rest.NewLinkHandler(linkService)
	mainScreenHandler := rest.NewMainScreenHandler(linkService, renderService)
	fingerprintHandler := rest.NewFingerprintHandler(linkService, fingerprintService)

	// Routes
	rest.MapLinkRoutes(app, linkHandler)
	rest.MapFPRoutes(app, fingerprintHandler)
	rest.MapMainScreenRoutes(app, mainScreenHandler)
	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	app.GET("/apple-app-site-association", func(c *gin.Context) {
		c.File("static/apple-app-site-association")
	})

	app.GET("/assetlinks.json", func(c *gin.Context) {
		c.File("static/assetlinks.json")
	})

	app.Static("/static", "./static")

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	address := "0.0.0.0:8080"

	bootWithGracefulShutdown(app, address, 5*time.Second)
}

func setupMiddlewares(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORS())
	router.Use(middlewares.ErrorHandlerMiddleware())
}

func bootWithGracefulShutdown(app *gin.Engine, address string, shutdownTimeoutS time.Duration) {
	srv := &http.Server{
		Addr:    address,
		Handler: app.Handler(),
	}

	log.Println("Listening on " + address)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutS)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Printf("timeout of %s seconds.", shutdownTimeoutS)
	log.Println("Server exiting")
}
