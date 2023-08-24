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
