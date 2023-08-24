package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"spw/api/common"
	"spw/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sashabaranov/go-openai"
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
				requestModel, _ := parseRequestMsg(message)
				log.Println(requestModel)

				if requestModel.Type == "chatGPT" {
					RequestGPT(ws, mt, requestModel.Data)
				} else {
					ws.WriteMessage(mt, []byte("unsupport type"))
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

func RequestGPT(ws *websocket.Conn, mt int, prompt string) {

	c := openai.NewClient(config.Get().Openai.Apikey)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		// MaxTokens: 4096,
		Temperature: 0.8,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		// fmt.Printf(response.Choices[0].Delta.Content)
		ws.WriteMessage(mt, []byte(response.Choices[0].Delta.Content))
	}
}
