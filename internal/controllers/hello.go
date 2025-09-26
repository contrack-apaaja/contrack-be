package controllers

import (
	"net/http"

	"contrack-be/internal/services/supabase"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	ok := supabase.IsConfigured()
	c.JSON(http.StatusOK, gin.H{
		"message":             "Hello, world!",
		"supabase_configured": ok,
	})
}
