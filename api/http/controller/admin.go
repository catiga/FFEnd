package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"time"

	"spw/api/common"
	"spw/model"
	database "spw/system"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	res := common.Response{}
	username := c.PostForm("username")
	password := c.PostForm("password")

	var result []model.AdminUser
	db := database.GetDb()
	err := db.Model(&model.AdminUser{}).Where("user_name = ? and password = ?", username, password).Find(&result).Error
	log.Println(err)

	if len(result) == 0 {
		res.Code = 100
		res.Msg = "admin not found"
		c.JSON(http.StatusOK, res)
		return
	}
	if len(result) > 1 {
		res.Code = 101
		res.Msg = "admin config error"
		c.JSON(http.StatusOK, res)
		return
	}

	db.Model(&model.AdminToken{}).Where("admin_id = ? and expire_time <= ?", result[0].Id, time.Now()).Updates(map[string]interface{}{
		"flag": -1,
	})

	var token model.AdminToken
	db.Model(&model.AdminToken{}).Where("admin_id = ? and expire_time > ?", result[0].Id, time.Now()).Last(&token)

	now := time.Now()
	exp := now.Add(30 * 24 * time.Hour)

	if token.Id == 0 {
		token.AdminId = result[0].Id
		token.CreateTime = &now
		token.Flag = 0
		token.ExpireTime = &exp

		hashToken := sha256.Sum256([]byte(strconv.Itoa(int(token.AdminId)) + result[0].UserName + result[0].Password))
		token.Token = hex.EncodeToString(hashToken[:])
		err := db.Save(&token).Error
		log.Println(err)
	}

	res.Code = 0
	res.Msg = "success"
	res.Data = map[string]interface{}{
		"token":    token.Token,
		"admin_id": token.AdminId,
	}
	c.JSON(http.StatusOK, res)
}

func CharacterList(c *gin.Context) {
	var result []model.Character
	db := database.GetDb()
	err := db.Model(&model.Character{}).Find(&result).Error
	log.Println(err)

	res := common.Response{}

	res.Code = 0
	res.Msg = "success"
	res.Data = result
	c.JSON(http.StatusOK, res)
}

func MethodList(c *gin.Context) {
	var result []model.Method
	db := database.GetDb()
	err := db.Model(&model.Method{}).Find(&result).Error
	log.Println(err)

	res := common.Response{}

	res.Code = 0
	res.Msg = "success"
	res.Data = result
	c.JSON(http.StatusOK, res)
}

func CharSettingList(c *gin.Context) {
	res := common.Response{}

	code := c.PostForm("code")
	lan := c.PostForm("lan")

	var result []model.CharBack
	db := database.GetDb()
	db.Model(&model.CharBack{}).Where("code = ? and lan = ? and flag = ?", code, lan, 0).Order("seq asc").Find(&result)

	res.Code = 0
	res.Msg = "success"
	res.Data = result
	c.JSON(http.StatusOK, res)
}

func CharSettingSave(c *gin.Context) {
	res := common.Response{}

	charIdStr := c.PostForm("charId")
	role := c.PostForm("role")
	prompt := c.PostForm("prompt")
	answer := c.PostForm("answer")
	seqStr := c.PostForm("seq")

	if role != "system" && role != "assistant" {
		res.Code = common.CODE_ERR_CHAR_ROLE_CAT
		res.Msg = "character not found"
		c.JSON(http.StatusOK, res)
		return
	}

	if len(prompt) == 0 {
		res.Code = common.CODE_ERR_CHAR_PARAM
		res.Msg = "param error"
		c.JSON(http.StatusOK, res)
		return
	}

	seq, err := strconv.ParseInt(seqStr, 10, 64)
	if err != nil {
		seq = 0
	}

	var cha model.Character
	db := database.GetDb()
	db.Model(&model.Character{}).Where("id = ?", charIdStr).First(&cha)

	if cha.Id == 0 {
		res.Code = common.CODE_ERR_CHAR_NOTFOUND
		res.Msg = "params error"
		c.JSON(http.StatusOK, res)
		return
	}
	//判断一下有多少基础设定了，这个不允许过多
	var count int64
	db.Model(&model.CharBack{}).Where("flag != ? and char_id = ?", -1, cha.Id).Count(&count)
	if count > 15 {
		res.Code = common.CODE_ERR_CHARBACK_MAX
		res.Msg = "exceed max character background settings count"
		c.JSON(http.StatusOK, res)
		return
	}

	now := time.Now()
	data := model.CharBack{
		BaseModel: model.BaseModel{
			Code: cha.Code,
			Lan:  cha.Lan,
			Flag: 0,
		},
		CharId:  cha.Id,
		Role:    role,
		Prompt:  prompt,
		Seq:     int(seq),
		AddTime: &now,
		Answer:  answer,
	}
	db.Save(&data)

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	c.JSON(http.StatusOK, res)
}

func CharSettingDel(c *gin.Context) {
	res := common.Response{}

	settingId := c.PostForm("setid")

	var cha model.CharBack
	db := database.GetDb()
	db.Model(&model.CharBack{}).Where("id = ?", settingId).First(&cha)

	if cha.Id == 0 {
		res.Code = common.CODE_ERR_CHAR_NOTFOUND
		res.Msg = "params error"
		c.JSON(http.StatusOK, res)
		return
	}
	cha.Flag = -1
	db.Updates(cha)

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	c.JSON(http.StatusOK, res)
}

func CharacterAdd(c *gin.Context) {
	res := common.Response{}

	idStr := c.PostForm("id")
	name := c.PostForm("name")
	info := c.PostForm("info")
	birth := c.PostForm("birth")
	age := c.PostForm("age")
	gender := c.PostForm("gender")
	place := c.PostForm("place")
	profile := c.PostForm("profile")
	region := c.PostForm("region")
	nature := c.PostForm("nature")

	avatar := c.PostForm("avatar")
	bodyImg := c.PostForm("bodyimg")

	code := c.PostForm("code")
	lan := c.PostForm("lan")

	if len(code) == 0 || len(lan) == 0 {
		res.Code = common.CODE_ERR_CHAR_BASPARAM
		res.Msg = "Basic params error"
		c.JSON(http.StatusOK, res)
		return
	}

	if lan != "zh-CN" && lan != "en" {
		res.Code = common.CODE_ERR_CHAR_BASPARAM
		res.Msg = "Basic params error"
		c.JSON(http.StatusOK, res)
		return
	}

	data := model.Character{}
	db := database.GetDb()

	//判断code和lan对应数据是否存在
	var existChar model.Character
	db.Model(&model.Character{}).Where("code = ? and lan = ? and flag != ?", code, lan, -1).Last(&existChar)
	if existChar.Id > 0 {
		res.Code = common.CODE_ERR_CHAR_EXIST
		res.Msg = "character exist"
		res.Data = existChar
		c.JSON(http.StatusOK, res)
		return
	}

	update := false
	if len(idStr) > 0 && idStr != "0" {
		db.Model(&data).Where("id = ?", idStr).Last(&data)
		if data.Id == 0 {
			res.Code = common.CODE_ERR_CHAR_NOTFOUND
			res.Msg = "character not found"
			c.JSON(http.StatusOK, res)
			return
		}
		update = true
	}

	natureInt, _ := strconv.Atoi(nature)

	data.CharName = name
	data.CharInfo = info
	data.CharBirth = birth
	data.CharAge = age
	data.CharGender = gender
	data.CharPlace = place
	data.CharProfile = profile
	data.CharNature = natureInt
	data.CharRegion = region
	data.CharAvatar = avatar
	data.CharFullBody = bodyImg

	data.Code = code
	data.Lan = lan
	data.Flag = 0

	if !update {
		err := db.Model(&model.Character{}).Create(&data).Error
		log.Println(err)
	} else {
		db.Updates(&data)
	}

	res.Code = common.CODE_SUCCESS
	res.Msg = "success"
	c.JSON(http.StatusOK, res)
}
