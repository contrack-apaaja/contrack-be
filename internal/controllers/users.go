package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"contrack-be/internal/repository"
)

func ListUsers(c *gin.Context) {
	repo := repository.NewUserRepo()
	users, err := repo.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
