package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tic3001-go-server/common/constant"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/controller"
)

func Register(engine *gin.Engine) {
	// register no route
	registerNoRoute(engine)
	// register api
	api := engine.Group("/api")
	registerNotesAPI(api)
}

func registerNotesAPI(group *gin.RouterGroup) {
	notes := group.Group("/notes")
	{
		notes.GET("/list", controller.NotesController.List)
		notes.POST("/create", controller.NotesController.Create)
		notes.PUT("/update", controller.NotesController.Update)
		notes.DELETE("/delete", controller.NotesController.Delete)
	}
}

func registerNoRoute(engine *gin.Engine) {
	engine.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, dto.ResponseDto{
			Code: constant.RespCodeResourceNotFound,
			Msg:  constant.RespMsgResourceNotFound,
			Data: "",
		})
	})
}
