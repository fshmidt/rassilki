package handler

import (
	_ "github.com/fshmidt/rassilki/docs"
	"github.com/fshmidt/rassilki/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	clients := router.Group("/clients")
	{

		clients.POST("/", h.createClient)
		clients.GET("/:id", h.getClient)
		clients.PUT("/:id", h.updateClient)
		clients.DELETE("/:id", h.deleteClient)

	}
	rassilki := router.Group("/rassilki")
	{
		rassilki.POST("/", h.createRassilka)
		rassilki.GET("/", h.getRassilkiReview)
		rassilki.GET("/:id", h.getRassilkaReviewById)
		rassilki.PUT("/:id", h.updateRassilka)
		rassilki.DELETE("/:id", h.deleteRassilka)
	}
	return router
}
