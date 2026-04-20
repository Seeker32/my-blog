package router

import (
	"net/http"

	"github.com/Seeker32/my-blog/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, userHandler handler.UserHandler) {
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	apiV1 := engine.Group("/api/v1")
	{
		users := apiV1.Group("/users")
		{
			users.POST("/register", userHandler.Register)
		}
	}
}