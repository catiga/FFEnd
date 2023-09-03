package controller

import (
	"net/http"
	"spw/api/common"
	"spw/model"
	database "spw/system"
	tool "spw/tool"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SystemMethods(c *gin.Context) {
	res := common.Response{}
	lan := c.PostForm("lan")
	if len(lan) == 0 {
		lan = c.Request.Header.Get("lan")
	}
	if !tool.IsSupportLan(lan) {
		res.Code = common.CODE_ERR_LAN
		res.Msg = "unsupport lan"
		c.JSON(http.StatusOK, res)
		return
	}

	var result []model.MysticalArt
	db := database.GetDb()

	db.Model(&model.MysticalArt{}).Where("lan = ? and flag != ?", lan, -1).Find(&result)

	var data []map[string]interface{}

	for _, v := range result {
		data = append(data, map[string]interface{}{
			"code":  v.Code,
			"lan":   v.Lan,
			"name":  v.Name,
			"brief": v.Brief,
		})
	}

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	res.Data = data

	c.JSON(http.StatusOK, res)
}

func SystemCatalogs(c *gin.Context) {
	res := common.Response{}
	lan := c.PostForm("lan")
	if len(lan) == 0 {
		lan = c.Request.Header.Get("lan")
	}
	if !tool.IsSupportLan(lan) {
		res.Code = common.CODE_ERR_LAN
		res.Msg = "unsupport lan"
		c.JSON(http.StatusOK, res)
		return
	}

	pidStr := c.PostForm("pid")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		pid = 0
	}

	var result []model.Catalog
	db := database.GetDb()

	if pid == 0 {
		db.Model(&model.Catalog{}).Where("lan = ? and flag != ?", lan, -1).Order("seq asc").Find(&result)
	} else if pid > 0 {
		db.Model(&model.Catalog{}).Where("lan = ? and flag != ? and parent = ?", lan, -1, pid).Order("seq asc").Find(&result)
	}

	var data []map[string]interface{}

	for _, v := range result {
		data = append(data, map[string]interface{}{
			"code": v.Code,
			"lan":  v.Lan,
			"name": v.Name,
			"id":   v.Id,
			"pid":  v.Parent,
		})
	}

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	res.Data = data

	c.JSON(http.StatusOK, res)
}
