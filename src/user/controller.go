package user

import (
	"net/http"
	"ws/src/common"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Repo *Repository
}

func NewController(repo *Repository) *Controller {
	return &Controller{Repo: repo}
}

func (ctrl Controller) Register(ctx *gin.Context) {
	var input User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	hashed, _ := common.HashPassword(input.Password)
	input.Password = hashed

	// Call Repo
	if err := ctrl.Repo.Create(&input); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "username or email existed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Register success!"})
}