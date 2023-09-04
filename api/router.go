package router

import (
	"log"
	"net/http"

	general "spw/api/http"
	"spw/api/interceptor"
	ws "spw/api/ws"

	"github.com/gin-gonic/gin"
)

type Option func(*gin.RouterGroup)

var options = []Option{}

func Include(opts ...Option) {
	options = append(options, opts...)
}

func Init() *gin.Engine {
	Include(ws.Routers)
	Include(general.Routers)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/index", helloHandler) //Default welcome api

	// staticTplPath := "./api/http/templates/**/*"

	// if configFilePathFromEnv := os.Getenv("DALINK_TPL_PATH"); configFilePathFromEnv != "" {
	// 	staticTplPath = configFilePathFromEnv + "/**/*"
	// }

	// r.LoadHTMLGlob(staticTplPath)

	apiGroup := r.Group("/spwapi", interceptor.HttpInterceptor()) // total interceptor stack
	for _, opt := range options {
		opt(apiGroup)
	}
	r.Run(":18080")
	return r
}

func helloHandler(c *gin.Context) {
	log.Println("hello")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello dalink",
	})
}
