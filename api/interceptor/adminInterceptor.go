package interceptor

import (
	"log"
	"net/http"
	"time"

	"spw/api/common"
	"spw/model"
	database "spw/system"

	"github.com/gin-gonic/gin"
)

// http 请求拦截器
func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/spwapi/admin/login" || path == "/spwapi/admin/sts" {
			c.Next()
		} else {
			log.Println(path)

			token := c.Request.Header.Get("Auth-Token")

			verifyResult := validateToken(token)
			if !verifyResult {
				c.Abort()
				c.JSON(http.StatusOK, common.Response{
					Code:      common.CODE_ERR_ADMINTOKEN,
					Msg:       "token invalid",
					Timestamp: time.Now().Unix(),
				})
				return
			}
			c.Next()
		}
	}
}

func validateToken(token string) bool {
	if len(token) == 0 {
		return false
	}
	db := database.GetDb()

	var result []model.AdminToken
	db.Model(&model.AdminToken{}).Where("token = ? and flag != ?", token, -1).Find(&result)

	ver := false
	for _, v := range result {
		if v.ExpireTime.After(time.Now()) {
			ver = true
		} else {
			err := db.Model(&model.AdminToken{}).Where("id = ?", v.Id).Update("flag", -1)
			log.Println(err)
		}
	}

	return ver
}

// func getSignParam(c *gin.Context, request *http.Request) (map[string]string, bool) {

// 	retMap := map[string]string{}

// 	paramMap := c.Request.URL.Query()

// 	for k, v := range paramMap {
// 		retMap[k] = strings.Join(v, "")
// 	}

// 	var parms []byte
// 	if c.Request.Body != nil {
// 		parms, _ = io.ReadAll(c.Request.Body)
// 		c.Request.Body = io.NopCloser(bytes.NewBuffer(parms))
// 	}

// 	formParams := strings.Split(string(parms), "&")

// 	for i := 0; i < len(formParams); i++ {

// 		fp := strings.Split(formParams[i], "=")

// 		if fp[0] == "" {
// 			continue
// 		}

// 		if len(fp) > 1 {
// 			retMap[fp[0]] = fp[1]
// 		} else {
// 			retMap[fp[0]] = ""
// 		}
// 	}

// 	retMap["ts"] = c.GetHeader("ts")
// 	retMap["device"] = c.GetHeader("device")
// 	return retMap, true
// }

// type bodyLogWriter struct {
// 	gin.ResponseWriter
// 	body *bytes.Buffer
// }

// func (w bodyLogWriter) Write(b []byte) (int, error) {
// 	w.body.Write(b)
// 	return w.ResponseWriter.Write(b)
// }
