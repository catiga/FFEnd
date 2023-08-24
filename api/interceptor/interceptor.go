package interceptor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var exception = []string{""}

// http 请求拦截器
func HttpInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		signMap, _ := getSignParam(c, c.Request)

		fmt.Sprintln("签名参数:      ", signMap)
		c.Next()
	}
}

func getSignParam(c *gin.Context, request *http.Request) (map[string]string, bool) {

	retMap := map[string]string{}

	paramMap := c.Request.URL.Query()

	for k, v := range paramMap {
		retMap[k] = strings.Join(v, "")
	}

	var parms []byte
	if c.Request.Body != nil {
		parms, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(parms))
	}

	formParams := strings.Split(string(parms), "&")

	for i := 0; i < len(formParams); i++ {

		fp := strings.Split(formParams[i], "=")

		if fp[0] == "" {
			continue
		}

		if len(fp) > 1 {
			retMap[fp[0]] = fp[1]
		} else {
			retMap[fp[0]] = ""
		}
	}

	retMap["ts"] = c.GetHeader("ts")
	retMap["device"] = c.GetHeader("device")
	return retMap, true
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
