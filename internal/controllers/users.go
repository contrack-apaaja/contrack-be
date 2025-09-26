package controllers

import (
	"context"

	"github.com/gin-gonic/gin"

	"contrack-be/internal/repository"
	"contrack-be/internal/utils"
)

func ListUsers(c *gin.Context) {
	repo := repository.NewUserRepo()
	users, err := repo.List(context.Background())
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve users")
		return
	}
	
	utils.OKResponse(c, "Users retrieved successfully", gin.H{
		"users": users,
		"count": len(users),
	})
}
