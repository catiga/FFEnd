package common

const (
	CODE_SUCCESS              = 0
	CODE_ERR_METHOD_UNSUPPORT = 1
	CODE_ERR_REQFORMAT        = 2
	CODE_ERR_LAN              = 901
	CODE_ERR_CHAR_UNKNOWN     = 101
	CODE_ERR_CHAR_NOTFOUND    = 102
	CODE_ERR_CHAR_ROLE_CAT    = 103
	CODE_ERR_CHAR_PARAM       = 104
	CODE_ERR_CHARBACK_MAX     = 105

	CODE_ERR_GPT_COMPLETE   = 201
	CODE_ERR_GPT_STREAM     = 202
	CODE_ERR_GPT_STREAM_EOF = 203
)

type Response struct {
	Code      int64
	Msg       string
	Timestamp int64
	Data      interface{}
}
