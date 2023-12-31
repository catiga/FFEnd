package router

import (
	"spw/api/http/controller"

	"github.com/gin-gonic/gin"

	"spw/api/interceptor"
)

func Routers(e *gin.RouterGroup) {
	adminGroup := e.Group("/admin", interceptor.CheckAdmin())
	adminGroup.POST("login", controller.Login)
	adminGroup.GET("characters", controller.CharacterList)
	adminGroup.GET("methods", controller.MethodList)
	adminGroup.POST("charsettings", controller.CharSettingList)
	adminGroup.POST("charsetsave", controller.CharSettingSave)
	adminGroup.POST("charsetdel", controller.CharSettingDel)
	adminGroup.POST("characteradd", controller.CharacterAdd)
	adminGroup.POST("sts", controller.Sts)
	adminGroup.POST("catalog/list", controller.CatalogList)
	adminGroup.POST("catalog/add", controller.CatalogAdd)
	adminGroup.POST("catalog/del", controller.CatalogDel)

	userGroup := e.Group("/user")
	userGroup.GET("character", controller.Characters)
	userGroup.POST("characters", controller.Characters)
	userGroup.POST("character", controller.CharacterWithCode)
	userGroup.POST("chat/history", controller.ChatHistory)
	userGroup.POST("samplechats", controller.ChatSamples)
	userGroup.POST("samplechat", controller.ChatSampleById)

	sysGroup := e.Group("/sys")
	sysGroup.POST("artags", controller.SystemMethods)
	sysGroup.POST("catags", controller.SystemCatalogs)
}
