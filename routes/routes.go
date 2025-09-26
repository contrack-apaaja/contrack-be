package routes

import (
	"contrack-be/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/hello", controllers.Hello)
	}
}
