package embedding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"spw/config"
	"spw/model"
	"strconv"
	"strings"
)

const uri = "https://api.openai.com"
const defaultModel = "text-embedding-ada-002"

const pinekey = "7d6cbcf1-a9a6-4b1c-8bf0-c3d784aff34f"
const pineuri = "https://zhz-288566b.svc.us-west4-gcp-free.pinecone.io"

var apikey = config.Get().Openai.Apikey

type GPT struct {
}

type EmbedResult struct {
	Object string                 `json:"object"`
	Model  string                 `json:"model"`
	Usage  map[string]interface{} `json:"usage"`
	Data   []EmbedData            `json:"data"`
}

type EmbedData struct {
	Index     int       `json:"index"`
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
}

type EmbedQueryResult struct {
	Id       string            `json:"id"`
	Score    float64           `json:"score"`
	Values   []float64         `json:"values"`
	Metadata map[string]string `json:"metadata"`
}

type EmbedQueryMatch struct {
	Results   []string           `json:"results"`
	Matches   []EmbedQueryResult `json:"matches"`
	Namespace string             `json:"namespace"`
}

func (*GPT) Embedding(content string, model string) (*EmbedResult, error) {
	body := map[string]string{
		"input": content,
		"model": model,
	}
	dataBytes, _ := json.Marshal(body)
	request, err := http.NewRequest(http.MethodPost, uri+"/v1/embeddings", bytes.NewBuffer(dataBytes))

	if err != nil {
		log.Println("Failed to create request:", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apikey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Failed to send request:", err)
		return nil, err
	}
	defer response.Body.Close()

	data, _ := io.ReadAll(response.Body)

	var v EmbedResult
	err = json.Unmarshal(data, &v)

	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (*GPT) SaveChatEmbeddings(data []EmbedResult, ids []model.ChatContent) error {
	if len(data) == 0 || len(data) != len(ids) {
		return errors.New("data length params error")
	}

	if len(data) > 100 {
		return errors.New("exceed max limitation")
	}

	var embReq []map[string]interface{}
	for index, _ := range data {
		if len(data[index].Data) > 0 {
			embReq = append(embReq, map[string]interface{}{
				"id":     strconv.FormatUint(ids[index].Id, 10),
				"values": data[index].Data[0].Embedding,
				"metadata": map[string]string{
					"user":   strconv.FormatUint(ids[index].UserId, 10),
					"devid":  ids[index].DevId,
					"char":   strconv.FormatUint(ids[index].CharId, 10),
					"direct": ids[index].Direction,
				},
			})
		} else {
			log.Println("pinecone save why:::", data[index].Data)
		}
	}

	bytesData, _ := json.Marshal(map[string]interface{}{
		"vectors": embReq,
	})

	payload := strings.NewReader(string(bytesData))

	req, _ := http.NewRequest("POST", pineuri+"/vectors/upsert", payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", pinekey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	log.Println("pinecone upsert:::", string(body))

	return nil
}

func (ins *GPT) BatchUpsert(data []model.ChatContent) error {

	var emb []EmbedResult
	for _, v := range data {
		d, err := ins.Embedding(v.Content, defaultModel)
		if err != nil {
			return err
		}
		emb = append(emb, *d)
	}

	return ins.SaveChatEmbeddings(emb, data)
}

func (ins *GPT) Query(id string, question string, filter map[string]string, limitation int) ([]EmbedQueryResult, error) {
	r, err := ins.Embedding(question, defaultModel)
	if err != nil {
		return nil, err
	}
	queryCond := map[string]interface{}{
		"filter":          filter,
		"topK":            limitation,
		"includeMetadata": true,
	}
	if len(id) > 0 {
		queryCond["id"] = id
	}

	queryCond["vector"] = r.Data[0].Embedding
	queryCond["includeValues"] = false

	jsonCond, _ := json.Marshal(queryCond)

	payload := strings.NewReader(string(jsonCond))

	req, _ := http.NewRequest("POST", pineuri+"/query", payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", pinekey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var m EmbedQueryMatch
	err = json.Unmarshal(body, &m)

	return m.Matches, err
}
