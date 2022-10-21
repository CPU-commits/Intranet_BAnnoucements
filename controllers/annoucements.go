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

// GetAnnoucements godoc
// @Summary     Get annoucements of home
// @Description Get genneral annoucements
// @Tags        annoucements
// @Accept      json
// @Produce     json
// @Param       skip  query    uint                                           false "int valid >= 0"    minimun(0) default(0)
// @Param       limit query    uint                                           false "int valid >= 0"    minimun(1) default(20) maximum(20)
// @Param       total query    bool                                           false "total length data" default(false)
// @Success     200   {object} res.Response{body=[]models.AnnoucementWLookUp} "desc"
// @Failure     400   {object} res.Response{}                                 "Skip must be a int"
// @Failure     400   {object} res.Response{}                                 "Limit must be a int"
// @Failure     503   {object} res.Response{}                                 "Service Unavailable - NATS || Mongo Service Unavailable"
// @Router      /get_annoucements [get]
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

// UploadAnnoucement godoc
// @Summary     Upload annoucement in home
// @Description Upload a annoucement in home, ROLS=[Director,Directive]
// @Tags        annoucements
// @Accept      json
// @Produce     json
// @Param       annoucement body     forms.AnnoucementForm true "Add annoucement - Files must be MongoId[]"
// @Success     200         {object} res.Response{}
// @Failure     400         {object} res.Response{} "Bad request - Bad body"
// @Failure     409         {object} res.Response{} "Todos los archivos a publicar deben ser públicos"
// @Failure     500         {object} res.Response{} "Server Internal Error - Maybe a bad response of NATS"
// @Failure     503         {object} res.Response{} "Service Unavailable - NATS || DB Service Unavailable"
// @Router      /upload_annoucement [post]
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

// DeleteAnnoucement godoc
// @Summary     Delete annoucement
// @Description Delete a annoucement in home, ROLS=[Director,Directive]
// @Tags        annoucements
// @Accept      json
// @Produce     json
// @Param       idAnnoucement path     string         true "MongoID Annoucement"
// @Success     200           {object} res.Response{} "desc"
// @Failure     400           {object} res.Response{} "idAnnoucement Must be a MongoID"
// @Failure     503           {object} res.Response{} "Service Unavailable - NATS || DB Service Unavailable"
// @Router      /delete_annoucement/{idAnnoucement} [delete]
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
