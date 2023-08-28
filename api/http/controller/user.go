package controller

import (
	"log"
	"net/http"
	"spw/api/common"
	"spw/model"
	database "spw/system"
	tool "spw/tool"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Characters(c *gin.Context) {
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

	var result []model.Character
	db := database.GetDb()
	err := db.Model(&model.Character{}).Where("lan = ?", lan).Find(&result).Error
	log.Println(err)

	var data []map[string]interface{}
	if len(result) > 0 {
		for _, v := range result {
			data = append(data, map[string]interface{}{
				"id":         v.Id,
				"name":       v.CharName,
				"avatar":     v.CharAvatar,
				"info":       v.CharInfo,
				"birth":      v.CharBirth,
				"age":        v.CharAge,
				"gender":     v.CharGender,
				"place":      v.CharPlace,
				"bodyImg":    v.CharFullBody,
				"profile":    v.CharProfile,
				"natureCode": v.CharNature,
				"code":       v.Code,
				"region":     v.CharRegion,
			})
		}
	}

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	res.Data = data
	c.JSON(http.StatusOK, res)
}

func CharacterWithCode(c *gin.Context) {
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

	code := c.PostForm("code")

	var v model.Character
	db := database.GetDb()
	err := db.Model(&model.Character{}).Where("lan = ? and code = ?", lan, code).Last(&v).Error
	log.Println(err)
	if v.Id == 0 {
		res.Code = common.CODE_ERR_CHAR_NOTFOUND
		res.Msg = "character not found"
		c.JSON(http.StatusOK, res)
		return
	}

	data := map[string]interface{}{
		"id":         v.Id,
		"name":       v.CharName,
		"avatar":     v.CharAvatar,
		"info":       v.CharInfo,
		"birth":      v.CharBirth,
		"age":        v.CharAge,
		"gender":     v.CharGender,
		"place":      v.CharPlace,
		"bodyImg":    v.CharFullBody,
		"profile":    v.CharProfile,
		"natureCode": v.CharNature,
		"code":       v.Code,
		"region":     v.CharRegion,
	}

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	res.Data = data
	c.JSON(http.StatusOK, res)
}

func ChatHistory(c *gin.Context) {
	res := common.Response{}

	useridstr := c.PostForm("userid")
	devId := c.Request.Header.Get("devid")
	charCode := c.PostForm("code")

	userid, err := strconv.ParseUint(useridstr, 10, 64)

	ql := "char_code = ?"
	var params []interface{}
	params = append(params, charCode)
	if err == nil && userid > 0 {
		ql += " and user_id = ?"
		params = append(params, userid)
	} else {
		if len(devId) > 0 {
			ql += " and dev_id = ?"
			params = append(params, devId)
		}
	}

	var result []model.ChatContent
	db := database.GetDb()

	if len(params) > 1 {
		db.Model(&model.ChatContent{}).Where(ql, params).Order("add_time asc").Find(&result)
	}

	var data []map[string]interface{}

	for _, v := range result {
		data = append(data, map[string]interface{}{
			"id":        v.Id,
			"userid":    v.UserId,
			"content":   v.Content,
			"direction": v.Direction,
			"charid":    v.CharId,
			"time":      v.AddTime,
		})
	}

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	res.Data = data

	c.JSON(http.StatusOK, res)
}
