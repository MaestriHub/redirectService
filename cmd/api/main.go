package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
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

	migrate := migrations.NewMigration_1(db)
	if err := migrate.Up(); err != nil {
		log.Fatal(err)
	}

	server := gin.New()

	setupMiddlewares(server)

	// Repositories
	linkRepo := repository.NewLinkRepo(db)
	salonInfoRepo := repository.NewSalonInfoRepo(db)
	empInfoRepo := repository.NewEmployeeInfoRepo(db)
	fingerprintRepo := repository.NewFingerprintRepo(db)

	// Services
	linkService := service.NewLinkService(linkRepo)
	renderService := service.NewRenderService(cfg.AppStoreLinksConfig, salonInfoRepo, empInfoRepo)
	fingerprintService := service.NewFingerprintService(fingerprintRepo, linkRepo)

	// Handlers
	linkHandler := rest.NewLinkHandler(linkService, renderService)
	mainScreenHandler := rest.NewMainScreenHandler(linkService, renderService)
	fingerprintHandler := rest.NewFingerprintHandler(linkService, fingerprintService)

	// Routes
	setupLinkRoutes(server, linkHandler)
	setupFingerprint(server, fingerprintHandler)
	setupMainScreen(server, mainScreenHandler)

	setupHealthRoutes(server)
	setupStaticFiles(server)
	setupSwagger(server)

	server.Run()
}

func setupLinkRoutes(router *gin.Engine, h rest.LinkHandler) {
	createLink := router.Group("link")
	{
		createLink.POST("/salon", h.CreateInviteSalon)
		createLink.POST("/employee", h.CreateInviteEmployee)
		createLink.POST("/client", h.CreateInviteClient)
	}
}

func setupHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
}

// TODO: протестировать
func setupStaticFiles(router *gin.Engine) {
	router.Static("/static", "./static")
}

func setupMainScreen(router *gin.Engine, h rest.MainScreenHandler) {
	router.GET("/:linkId", h.MainScreen)
}

func setupSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func setupFingerprint(router *gin.Engine, h rest.FingerprintHandler) {
	fp := router.Group("fingerprint")
	{
		fp.POST("/:linkId", h.Create)
		fp.POST("/find/*linkId", h.Find)
	}
}

func setupMiddlewares(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORS())
}
