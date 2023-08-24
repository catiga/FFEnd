package sys

import (
	"fmt"

	"time"

	"spw/config"
	"spw/tool"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Ewriter struct {
	mlog *tool.Elog
}

func (m *Ewriter) Printf(format string, v ...interface{}) {
	logstr := fmt.Sprintf(format, v...)
	m.mlog.Info(logstr)
}

func NewWriter() *Ewriter {
	log := tool.Vlog
	//配置logrus
	// log.SetFormatter(&logrus.JSONFormatter{
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })

	return &Ewriter{mlog: log}
}

var DB *gorm.DB

func init() {
	cfg := config.Get()
	fmt.Println(cfg.Mysql.Username)

	newLogger := logger.New(
		// log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		NewWriter(),
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)

	// _ = newLogger

	dsn := cfg.Mysql.Username + ":" + cfg.Mysql.Password + "@tcp(" + cfg.Mysql.Address + ":" + cfg.Mysql.Port + ")/" + cfg.Mysql.Db + "?charset=utf8&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		tool.Vlog.Fatal(err)
	}
	DB = database
}

func GetDb() *gorm.DB {
	return DB
}
