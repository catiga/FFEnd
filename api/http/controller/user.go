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

	catCode := c.PostForm("cat")
	methodCode := c.PostForm("meth")

	db := database.GetDb()

	filter := false
	var ids []uint64
	if len(catCode) > 0 {
		filter = true
		err := db.Model(&model.CharacterPos{}).Select("char_id").Where("type_cat = ? and type_code = ? and type_lan = ? and flag != ?", "100", catCode, lan, -1).Find(&ids).Error
		if err != nil {
			log.Println("find cats error:", err)
		}
	}
	if len(methodCode) > 0 {
		filter = true
		var tmpids []uint64
		err := db.Model(&model.CharacterPos{}).Select("char_id").Where("type_cat = ? and type_code = ? and type_lan = ? and flag != ?", "200", catCode, lan, -1).Find(&tmpids).Error
		if err != nil {
			log.Println("find methods error:", err)
		}
		for _, v := range tmpids {
			isIn := false
			for _, w := range ids {
				if v == w {
					isIn = true
					break
				}
			}
			if !isIn {
				ids = append(ids, v)
			}
		}
	}

	var params []interface{}
	sql := "lan = ? and flag != ?"
	params = append(params, lan)
	params = append(params, -1)

	if filter {
		sql = sql + " and char_id IN (?)"
		params = append(params, ids[:])
	}

	var result []model.Character

	// err := db.Model(&model.Character{}).Where("lan = ?", lan).Find(&result).Error
	err := db.Model(&model.Character{}).Where(sql, params[:]).Find(&result).Error
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
	devId := c.Request.Header.Get("Devid")
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

	log.Println(useridstr, devId, charCode, params)

	var result []model.ChatContent
	db := database.GetDb()

	if len(params) > 1 {
		err := db.Model(&model.ChatContent{}).Where(ql, params...).Order("add_time asc").Find(&result).Error
		log.Println("find history:", err, ql)
	}

	var data []map[string]interface{}

	for _, v := range result {
		data = append(data, map[string]interface{}{
			"id":        v.Id,
			"userid":    v.UserId,
			"content":   v.Question,
			"direction": "1",
			"charid":    v.CharId,
			"time":      v.AddTime,
		})
		data = append(data, map[string]interface{}{
			"id":        v.Id,
			"userid":    v.UserId,
			"content":   v.Reply,
			"direction": "2",
			"charid":    v.CharId,
			"time":      v.AddTime,
		})
	}

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	res.Data = data

	c.JSON(http.StatusOK, res)
}

func ChatSamples(c *gin.Context) {
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

	var result []model.SampleChat
	db := database.GetDb()

	db.Model(&model.SampleChat{}).Where("lan = ? and flag != ?", lan, -1).Find(&result)

	var data []map[string]interface{}

	for _, v := range result {
		data = append(data, map[string]interface{}{
			"code":   v.Code,
			"lan":    v.Lan,
			"Q":      v.Samq,
			"A":      v.Sama,
			"charid": v.CharId,
			"id":     v.Id,
		})
	}

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	res.Data = data

	c.JSON(http.StatusOK, res)
}

func ChatSampleById(c *gin.Context) {
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

	idStr := c.PostForm("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		res.Code = common.CODE_SUCCESS
		res.Msg = "id not found"
		c.JSON(http.StatusOK, res)
		return
	}

	var result model.SampleChat
	db := database.GetDb()

	db.Model(&model.SampleChat{}).Where("id = ? and flag != ?", id, -1).Find(&result)

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	res.Data = map[string]interface{}{
		"code":   result.Code,
		"lan":    result.Lan,
		"Q":      result.Samq,
		"A":      result.Sama,
		"charid": result.CharId,
		"id":     result.Id,
	}

	c.JSON(http.StatusOK, res)
}
