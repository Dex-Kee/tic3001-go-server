package controller

import (
	"github.com/gin-gonic/gin"
	"math"
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
	mockUsers := service.UserService.FindMockUser()
	// to prevent the crash of browser, only take the first 1k element
	size := int(math.Min(float64(len(mockUsers)), 1000))
	c.JSON(http.StatusOK, dto.GetSuccessRespDto(mockUsers[0:size+1]))
}

func (controller userController) DeleteUser(c *gin.Context) {
	// used for the demo of invalid access of request, simply return success 200
	c.JSON(http.StatusOK, dto.GetSuccessRespDto(""))
}
