package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CPU-commits/Intranet_BAnnoucements/controllers"
	"github.com/CPU-commits/Intranet_BAnnoucements/docs"
	"github.com/CPU-commits/Intranet_BAnnoucements/middlewares"
	"github.com/CPU-commits/Intranet_BAnnoucements/models"
	"github.com/CPU-commits/Intranet_BAnnoucements/res"
	"github.com/CPU-commits/Intranet_BAnnoucements/settings"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go.uber.org/zap"
)

var settingsData = settings.GetSettings()

func Init() {
	router := gin.New()
	// Zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/api/annoucements/swagger"},
	}))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Server Internal Error: %s", err))
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, res.Response{
			Success: false,
			Message: "Server Internal Error",
		})
	}))

	router.GET("/api/annoucements/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong"+fmt.Sprint(time.Now().Unix()))
	})
	// Docs
	docs.SwaggerInfo.BasePath = "/api/annoucements"
	docs.SwaggerInfo.Version = "v1"
	docs.SwaggerInfo.Host = "localhost:8080"
	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{settingsData.CLIENT_URL},
		AllowMethods:     []string{"GET", "OPTIONS", "PUT", "DELETE", "POST"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Secure
	sslUrl := "ssl." + settingsData.CLIENT_URL
	router.Use(secure.New(secure.Config{
		AllowedHosts:         []string{settingsData.CLIENT_URL, sslUrl},
		SSLHost:              sslUrl,
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		IENoOpen:             true,
		ReferrerPolicy:       "strict-origin-when-cross-origin",
		SSLProxyHeaders: map[string]string{
			"X-Fowarded-Proto": "https",
		},
	}))
	// Routes
	defaultRoles := []string{
		models.DIRECTIVE,
		models.DIRECTOR,
	}
	annoucements := router.Group(
		"api/annoucements",
		middlewares.JWTMiddleware(),
	)
	{
		// Init controllers
		annoucementsController := new(controllers.AnnoucementController)
		// Define routes
		// Annoucements
		annoucements.GET(
			"/get_annoucements",
			annoucementsController.GetAnnoucements,
		)
		annoucements.POST(
			"/upload_annoucement",
			middlewares.RolesMiddleware(defaultRoles),
			annoucementsController.UploadAnnoucement,
		)
		annoucements.DELETE(
			"/delete_annoucement/:idAnnoucement",
			middlewares.RolesMiddleware(defaultRoles),
			annoucementsController.DeleteAnnoucement,
		)
	}
	// Route docs
	router.GET("/api/annoucements/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// No route
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, res.Response{
			Success: false,
			Message: "Not found",
		})
	})
	// Init server
	if err := router.Run(); err != nil {
		log.Fatalf("Error init server")
	}
}
