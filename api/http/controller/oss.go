package controller

import (
	"net/http"
	"spw/api/common"
	"spw/util"

	"github.com/gin-gonic/gin"
)

func Sts(c *gin.Context) {
	res := common.Response{}
	sts := util.GetSts()

	if sts == "" {

		res.Code = 100
		res.Msg = "No sts acquired"
		res.Data = map[string]interface{}{}
		c.JSON(http.StatusOK, res)
		return
	}

	res.Code = 0
	res.Msg = "success"
	res.Data = sts
	c.JSON(http.StatusOK, res)
}
