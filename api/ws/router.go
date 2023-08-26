package router

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"spw/api/common"
	"spw/config"
	"spw/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sashabaranov/go-openai"

	apicommon "spw/api/common"
	database "spw/system"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Routers(e *gin.RouterGroup) {
	userGroup := e.Group("/ws")

	userGroup.GET("", func(c *gin.Context) {
		ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Fatalln(err)
		}
		defer ws.Close()
		go func() {
			<-c.Done()
			log.Println("ws lost connection")
		}()

		timeNowHs := time.Now().UnixNano() / int64(time.Millisecond)

		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println("read error")
				fmt.Println(err)
				break
			}
			if string(message) == "ping" { //heart beat
				message = []byte("pong")
				err = ws.WriteMessage(mt, message)
				if err != nil {
					log.Println(err)
					break
				}
			} else {
				requestModel, err := parseRequestMsg(message)
				if err != nil {
					rp := makeReply(common.CODE_ERR_REQFORMAT, err.Error(), timeNowHs, "", requestModel.Timestamp, "")
					ws.WriteJSON(rp)
					return
				}

				if requestModel.Method == apicommon.METHOD_GPT {
					RequestGPT(ws, mt, requestModel, timeNowHs)
				} else {
					rp := makeReply(common.CODE_ERR_METHOD_UNSUPPORT, err.Error(), timeNowHs, "", requestModel.Timestamp, "")
					ws.WriteJSON(rp)
				}
			}

		}
	})
}

func parseRequestMsg(body []byte) (c common.Request, e error) {

	defer func() {
		if r := recover(); r != nil {
			e = errors.New("invalid request data format")
		}
	}()

	err := json.Unmarshal(body, &c)
	if err != nil {
		panic(err)
	}

	return c, nil
}

func RequestGPT(ws *websocket.Conn, mt int, request common.Request, timeNowHs int64) {
	ascode := request.Ascode
	language := request.Lan
	chatType := request.Type
	question := request.Data

	db := database.GetDb()
	var character model.Character
	err := db.Model(&model.Character{}).Where("lan = ? and code = ?", language, ascode).Last(&character).Error

	if err != nil {
		log.Println("chat error:", err)
		rp := makeReply(common.CODE_ERR_CHAR_UNKNOWN, err.Error(), timeNowHs, "", request.Timestamp, "")
		ws.WriteJSON(rp)
		return
	}
	if character.Id == 0 {
		rp := makeReply(common.CODE_ERR_CHAR_NOTFOUND, "character not found", timeNowHs, "", request.Timestamp, "")
		ws.WriteJSON(rp)
		return
	}

	c := openai.NewClient(config.Get().Openai.Apikey)
	ctx := context.Background()

	background := buildPrompt(&character, chatType, question)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		// MaxTokens: 4096,
		Temperature: 0.8,
		// Messages: []openai.ChatCompletionMessage{
		// 	{
		// 		Role:    openai.ChatMessageRoleUser,
		// 		Content: prompt,
		// 	},
		// },
		Messages: background,
		Stream:   true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Println("ChatCompletionStream error:", err)

		rp := makeReply(common.CODE_ERR_GPT_COMPLETE, err.Error(), timeNowHs, "", request.Timestamp, "")

		ws.WriteJSON(rp)
		return
	}
	defer stream.Close()

	log.Println("Stream response: ")

	chatHash := generateChatHash(timeNowHs, request)

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			log.Println("\nStream error:", err)

			rp := makeReply(common.CODE_ERR_GPT_STREAM, err.Error(), timeNowHs, "", request.Timestamp, "")

			ws.WriteJSON(rp)
			return
		}

		rp := makeReply(common.CODE_SUCCESS, "success", timeNowHs, chatHash, request.Timestamp, response.Choices[0].Delta.Content)

		ws.WriteJSON(rp)
	}
}

func buildPrompt(chars *model.Character, chatType string, question string) []openai.ChatCompletionMessage {
	var back []openai.ChatCompletionMessage

	db := database.GetDb()

	var result []model.CharBack
	db.Model(&model.CharBack{}).Where("code = ? and lan = ? and flag = ?", chars.Code, chars.Lan, 0).Order("seq asc").Find(&result)

	if len(result) > 0 {
		for _, v := range result {
			roleType := ""
			if v.Role == "system" {
				roleType = openai.ChatMessageRoleSystem
			} else if v.Role == "assistant" {
				roleType = openai.ChatMessageRoleAssistant
			}
			if len(roleType) > 0 {
				back = append(back, openai.ChatCompletionMessage{
					Role:    roleType,
					Content: v.Prompt,
				})
			}
		}
	}
	back = append(back, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: question,
	})
	return back
}

func generateChatHash(timeHs int64, request common.Request) string {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(100000)
	chatHash := strconv.FormatInt(timeHs, 10) + "-" + strconv.FormatInt(request.Timestamp, 10) + "-" + strconv.FormatInt(int64(randomInt), 10)

	hashByte := sha256.Sum256([]byte(chatHash))

	return hex.EncodeToString(hashByte[:])
}

func makeReply(code int64, msg string, timeHs int64, chatId string, replyTs int64, content string) *common.Response {
	return &common.Response{
		Code:      code,
		Msg:       msg,
		Timestamp: timeHs,
		Data: map[string]interface{}{
			"Id":      chatId,
			"ReplyTs": replyTs,
			"Content": content,
		},
	}
}
