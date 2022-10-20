package controllers

import (
	"net/http"
	"strconv"

	"github.com/CPU-commits/Intranet_BAnnoucements/forms"
	"github.com/CPU-commits/Intranet_BAnnoucements/res"
	"github.com/CPU-commits/Intranet_BAnnoucements/services"
	"github.com/gin-gonic/gin"
)

type AnnoucementController struct{}

// Services
var annoucementService = services.NewAnnoucementService()

func (a *AnnoucementController) GetAnnoucements(c *gin.Context) {
	skipQuery := c.DefaultQuery("skip", "0")
	limitQuery := c.DefaultQuery("limit", "20")
	totalQuery := c.DefaultQuery("total", "false")
	// Transform query params
	skipInt, err := strconv.Atoi(skipQuery)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &res.Response{
			Message: "Skip must be a int",
		})
		return
	}
	limitInt, err := strconv.Atoi(limitQuery)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &res.Response{
			Message: "Limit must be a int",
		})
	}
	totalBool := totalQuery == "true"
	// Get annoucement
	annoucements, total, errRes := annoucementService.GetAnnoucements(skipInt, limitInt, totalBool)
	if errRes != nil {
		c.AbortWithStatusJSON(errRes.StatusCode, &res.Response{
			Message: errRes.Err.Error(),
		})
		return
	}

	// Responses
	response := make(map[string]interface{})
	response["annoucements"] = annoucements
	response["total"] = total
	c.JSON(http.StatusOK, &res.Response{
		Success: true,
		Data:    response,
	})
}

func (a *AnnoucementController) UploadAnnoucement(c *gin.Context) {
	var annoucement forms.AnnoucementForm

	if err := c.BindJSON(&annoucement); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &res.Response{
			Message: err.Error(),
		})
		return
	}
	claims, _ := services.NewClaimsFromContext(c)

	insertedId, err := annoucementService.UploadAnnoucement(annoucement, claims.ID)
	if err != nil {
		c.AbortWithStatusJSON(err.StatusCode, &res.Response{
			Message: err.Err.Error(),
		})
		return
	}
	// Response
	response := make(map[string]interface{})
	response["_id"] = insertedId.Hex()
	c.JSON(http.StatusOK, &res.Response{
		Success: true,
		Data:    response,
	})
}

func (a *AnnoucementController) DeleteAnnoucement(c *gin.Context) {
	idAnnoucement := c.Param("idAnnoucement")
	// Delete
	err := annoucementService.DeleteAnnoucement(idAnnoucement)
	if err != nil {
		c.AbortWithStatusJSON(err.StatusCode, &res.Response{
			Message: err.Err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &res.Response{
		Success: true,
	})
}
