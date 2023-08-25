package router

import (
	"spw/api/http/controller"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {
	adminGroup := e.Group("/admin")
	adminGroup.POST("login", controller.Login)
	adminGroup.GET("characters", controller.CharacterList)
	adminGroup.GET("methods", controller.MethodList)

	userGroup := e.Group("/user")
	userGroup.GET("character", controller.Characters)
	userGroup.POST("character", controller.CharacterWithCode)
}
