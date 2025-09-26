package router

import (
	"contrack-be/internal/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/hello", controllers.Hello)
		api.GET("/users", controllers.ListUsers)
	}
}
