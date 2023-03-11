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
	// api.Use(middleware.AuthFileter.ValidateResource)
	registerAuthAPI(api)
	registerUserAPI(api)
	registerNotesAPI(api)
}

func registerAuthAPI(group *gin.RouterGroup) {
	auth := group.Group("/auth")
	{
		auth.POST("/login", controller.AuthController.Login)
	}
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

func registerUserAPI(group *gin.RouterGroup) {
	user := group.Group("/user")
	{
		user.GET("/generate-mock-data", controller.UserController.GenerateUserTestData)
		user.GET("/list/mock", controller.UserController.FindMockUser)
		user.DELETE("/delete", controller.UserController.DeleteUser)
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
