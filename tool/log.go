package tool

import (
	"github.com/sirupsen/logrus"
)

type Elog struct {
	*logrus.Logger
}

var Vlog *Elog // for general output

func init() {
	Vlog = &Elog{
		logrus.New(),
	}
	// mlog = logrus.New()
	Vlog.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})
}
