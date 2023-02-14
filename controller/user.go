package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tic3001-go-server/common/dto"
)

type userController struct{}

var UserController = new(userController)

func (controller userController) DeleteUser(c *gin.Context) {
	// used for the demo of invalid access of request, simply return success 200
	c.JSON(http.StatusOK, dto.GetSuccessRespDto(""))
}
