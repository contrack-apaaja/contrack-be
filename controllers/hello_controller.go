package controllers

import (
	"net/http"

	"contrack-be/services/supabase"

	"github.com/gin-gonic/gin"
)

// Hello godoc
// @Summary Hello endpoint
// @Produce json
func Hello(c *gin.Context) {
	ok := supabase.IsConfigured()
	c.JSON(http.StatusOK, gin.H{
		"message":             "Hello, world!",
		"supabase_configured": ok,
	})
}
