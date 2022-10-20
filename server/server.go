package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CPU-commits/Intranet_BAnnoucements/controllers"
	"github.com/CPU-commits/Intranet_BAnnoucements/middlewares"
	"github.com/CPU-commits/Intranet_BAnnoucements/models"
	"github.com/CPU-commits/Intranet_BAnnoucements/res"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Server Internal Error: %s", err))
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, res.Response{
			Success: false,
			Message: "Server Internal Error",
		})
	}))
	// CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
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
