package controllers

import (
	"contrack-be/internal/services/supabase"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	ok := supabase.IsConfigured()
	utils.OKResponse(c, "Hello from Contrack API!", gin.H{
		"supabase_configured": ok,
		"version":            "1.0.0",
		"service":            "authentication",
	})
}
