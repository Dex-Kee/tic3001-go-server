package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/service"
)

type userController struct{}

var UserController = new(userController)

func (controller userController) GenerateUserTestData(c *gin.Context) {
	service.UserService.GenerateTestData()
	c.JSON(http.StatusOK, dto.GetSuccessRespDto(""))
}

func (controller userController) FindMockUser(c *gin.Context) {
	service.UserService.FindMockUser()
	c.JSON(http.StatusOK, dto.GetSuccessRespDto(""))
}

func (controller userController) DeleteUser(c *gin.Context) {
	// used for the demo of invalid access of request, simply return success 200
	c.JSON(http.StatusOK, dto.GetSuccessRespDto(""))
}
