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
	adminGroup.POST("charsettings", controller.CharSettingList)
	adminGroup.POST("charsetsave", controller.CharSettingSave)
	adminGroup.POST("charsetdel", controller.CharSettingDel)
	adminGroup.POST("characteradd", controller.CharacterAdd)
	adminGroup.POST("sts", controller.Sts)

	userGroup := e.Group("/user")
	userGroup.GET("character", controller.Characters)
	userGroup.POST("character", controller.CharacterWithCode)
	userGroup.POST("chat/history", controller.ChatHistory)
}
