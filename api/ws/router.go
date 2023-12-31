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
	"math"
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
	"spw/embedding"
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
	err := db.Model(&model.Character{}).Where("lan = ? and code = ? and flag != ?", language, ascode, -1).Last(&character).Error

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

	defaultModelKey := config.Get().Openai.Apikey
	defaultModelName := openai.GPT3Dot5Turbo

	if len(character.ModelKey) > 0 && len(character.ModelName) > 0 {
		defaultModelKey = character.ModelKey
		defaultModelName = character.ModelName
		log.Println("replace default model：", defaultModelName)
	}
	c := openai.NewClient(defaultModelKey)
	ctx := context.Background()

	background := buildPrompt(&character, chatType, request, question)

	defaultTemp := 0.5
	if character.CharNature >= 0 && character.CharNature <= 100 {
		vs := float64(character.CharNature) / 100
		defaultTemp = math.Round(vs*10) / 10
	}
	req := openai.ChatCompletionRequest{
		Model: defaultModelName, //openai.GPT3Dot5Turbo,
		// MaxTokens: 4096,
		// Temperature: 0.8,
		Temperature: float32(defaultTemp),
		// Messages: []openai.ChatCompletionMessage{
		// 	{
		// 		Role:    openai.ChatMessageRoleUser,
		// 		Content: prompt,
		// 	},
		// },
		Messages: background,
		Stream:   true,
	}

	chatIn := time.Now()
	//Save chat data
	chat := model.ChatContent{
		Flag:     0,
		DevId:    request.DevId,
		UserId:   request.UserId,
		CharId:   character.Id,
		Question: question,
		Reply:    "",
		AddTime:  &chatIn,
		CharCode: character.Code,
	}
	// db.Save(&chat)

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

	replyMsg := ""

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("\nStream EOF finished")
			chat.Reply = replyMsg
			db.Save(&chat)

			go func(chat *model.ChatContent) {
				gpt := &embedding.GPT{}
				gpt.BatchUpsert(&embedding.EmbededUpsertData{
					QuestionId: chat.Id,
					Question:   chat.Question,
					Reply:      chat.Reply,
					UserId:     chat.UserId,
					DevId:      chat.DevId,
					CharId:     chat.CharId,
					CharCode:   chat.CharCode,
				})
			}(&chat)
			return
		}

		if err != nil {
			log.Println("\nStream error:", err)
			rp := makeReply(common.CODE_ERR_GPT_STREAM, err.Error(), timeNowHs, "", request.Timestamp, "")
			ws.WriteJSON(rp)
			return
		}

		rp := makeReply(common.CODE_SUCCESS, "success", timeNowHs, chatHash, request.Timestamp, response.Choices[0].Delta.Content)
		replyMsg += response.Choices[0].Delta.Content
		ws.WriteJSON(rp)
	}
}

func buildPrompt(chars *model.Character, chatType string, request common.Request, question string) []openai.ChatCompletionMessage {
	var back []openai.ChatCompletionMessage

	db := database.GetDb()

	var result []model.CharBack
	db.Model(&model.CharBack{}).Where("code = ? and lan = ? and flag = ?", chars.Code, chars.Lan, 0).Order("seq asc").Find(&result)

	gpt := &embedding.GPT{}
	metaFilter := map[string]string{
		"charid": strconv.FormatUint(chars.Id, 10),
	}
	if request.UserId > 0 {
		metaFilter["user"] = strconv.FormatUint(request.UserId, 10)
	}
	if len(request.DevId) > 0 {
		metaFilter["devid"] = request.DevId
	}
	embResults, err := gpt.Query("", question, metaFilter, 3)

	backgroundContext := ""
	if err == nil && len(embResults) > 0 {
		var ids []uint64
		for _, v := range embResults {
			if v.Metadata["user"] == strconv.FormatUint(request.UserId, 10) || v.Metadata["devid"] == request.DevId {
				if v.Score > float64(0.66) {
					// if len(ids) > 5 {
					// 	break
					// }
					idint, err := strconv.ParseUint(v.Id, 10, 64)
					if err == nil {
						ids = append(ids, idint)
					}
				}
			}
		}
		var result_1 []model.ChatContent
		db.Where("id IN (?)", ids).Order("add_time desc").Find(&result_1)

		//同时取最近的三条聊天记录补充进去
		result_recent := findRecentChats(3, request.DevId, request.UserId, chars)
		if len(result_recent) > 0 {
			log.Println("add new chat content: ", len(result_recent))
			result_1 = append(result_1, result_recent...)
		}

		if len(result_1) > 0 { // here is related chat history data
			log.Println("find appendix user data:", len(result_1), ids, " and start to build question background")
			for _, v := range result_1 {
				backgroundContext += ("Q:`" + v.Question + "`;A:`" + v.Reply + "` \n")
			}
		}
	}

	log.Println("开始构建角色设定：", len(result))

	if len(result) > 0 {
		for _, v := range result {
			log.Println(v.Role, v.Prompt)
			roleType := ""
			if v.Role == "system" {
				roleType = openai.ChatMessageRoleSystem
				back = append(back, openai.ChatCompletionMessage{
					Role:    roleType,
					Content: v.Prompt,
				})
			} else if v.Role == "assistant" {
				back = append(back, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: v.Prompt,
				})
				back = append(back, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: v.Answer,
				})
			} else if v.Role == "user" {
				roleType = openai.ChatMessageRoleUser
			}
		}
	}
	if len(backgroundContext) > 0 {
		backgroundContext = "Context: \n" + backgroundContext + "\n"
	}
	log.Println("Question with Context:", backgroundContext+question)
	back = append(back, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: backgroundContext + question,
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

func findRecentChats(count int, dev string, user uint64, character *model.Character) []model.ChatContent {
	db := database.GetDb()

	var result_recent []model.ChatContent
	var params []interface{}
	sql := "char_code = ? and flag != ?"
	params = append(params, character.Code)
	params = append(params, -1)

	canFind := false
	if user > 0 {
		sql = sql + " and user_id = ?"
		params = append(params, user)
		canFind = true
	} else {
		if len(dev) > 0 {
			canFind = true
			sql = sql + " and dev_id = ?"
			params = append(params, dev)
		}
	}

	if !canFind {
		return result_recent
	}

	err := db.Where(sql, params...).Order("add_time desc").Limit(3).Find(&result_recent).Error
	if err != nil {
		log.Println("find recent chats error:", err)
	}
	return result_recent
}
