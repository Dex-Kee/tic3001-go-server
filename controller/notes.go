package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/service"
	"tic3001-go-server/validation"
)

type notesController struct{}

var NotesController = notesController{}

func (controller *notesController) List(c *gin.Context) {
	keyword := c.Query("keyword")
	log.Info("keyword: ", keyword)

	list := service.NotesService.List(keyword)

	c.JSON(http.StatusOK, dto.GetSuccessRespDto(list))
}

func (controller *notesController) Create(c *gin.Context) {
	form := new(dto.NotesForm)
	err := c.BindJSON(form)
	if err != nil {
		log.Error("err when bind notes form: ", err.Error())
		c.JSON(http.StatusInternalServerError, dto.GetServerErrorRespDto())
		return
	}

	log.Info("notes form: ", form)
	err = validation.NotesValidationService.FormChecker(*form)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetClientParamErrorRespDto(err.Error()))
		return
	}

	service.NotesService.Create(*form)

	c.JSON(http.StatusOK, dto.GetSuccessRespDto(nil))
}

func (controller *notesController) Update(c *gin.Context) {
	form := new(dto.NotesForm)
	err := c.BindJSON(form)
	if err != nil {
		log.Error("err when bind notes form: ", err.Error())
		c.JSON(http.StatusInternalServerError, dto.GetServerErrorRespDto())
		return
	}

	log.Info("notes form: ", form)

	err = validation.NotesValidationService.FormChecker(*form)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetClientParamErrorRespDto(err.Error()))
		return
	}

	err = validation.NotesValidationService.EntityExistedChecker(form.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetClientParamErrorRespDto(err.Error()))
		return
	}

	service.NotesService.Update(*form)

	c.JSON(http.StatusOK, dto.GetSuccessRespDto(nil))
}

func (controller *notesController) Delete(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, dto.GetClientParamErrorRespDto("error: deletion id is not specified"))
		return
	}

	log.Info("id: ", id)

	err := validation.NotesValidationService.EntityExistedChecker(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetClientParamErrorRespDto(err.Error()))
		return
	}

	service.NotesService.Delete(id)

	c.JSON(http.StatusOK, dto.GetSuccessRespDto(nil))
}
