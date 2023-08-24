package common

type Request struct {
	Type   string `json:"type"`
	Method string `json:"method"`
	Data   string `json:"data"`
}
