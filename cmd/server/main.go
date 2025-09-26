package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"contrack-be/internal/config"
	"contrack-be/internal/router"
	sup "contrack-be/internal/services/supabase"
)

func main() {
	cfg := config.Load()

	// initialize supabase (internal service reads env too; keep for clarity)
	sup.Init()

	r := gin.Default()
	router.Setup(r)

	fmt.Printf("starting server on :%s\n", cfg.Port)
	r.Run(":" + cfg.Port)
}
