package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/service"
)

type authController struct{}

var AuthController = new(authController)

func (controller authController) Login(c *gin.Context) {
	var form dto.LoginForm
	err := c.BindJSON(&form)
	if err != nil {
		log.Error("err when bind login form: ", err.Error())
		c.JSON(http.StatusInternalServerError, dto.GetServerErrorRespDto())
		return
	}

	token, err := service.AuthService.Login(form)
	if err != nil {
		log.Error("err when call Login api: ", err.Error())
		c.JSON(http.StatusBadRequest, dto.GetClientParamErrorRespDto(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.GetSuccessRespDto(token))
}
