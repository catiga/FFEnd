package config

import (
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	"spw/tool"
)

type Config struct {
	Mysql struct {
		Port     string `yaml:"port"`
		Address  string `yaml:address`
		Username string `yaml:username`
		Password string `yaml:password`
		Db       string `yaml:db`
	}
	Log struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
		Node  string `yaml:"node"`
	}
	Openai struct {
		Apikey string `yaml:"apikey"`
	}
}

var cfg *Config

func getRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(filename))
	return root
}

func init() {
	rootPath := getRoot()

	env := "local"

	confFilePath := rootPath + "/config/param-" + strings.ToLower(env) + ".yaml"

	if configFilePathFromEnv := os.Getenv("DALINK_GO_CONFIG_PATH"); configFilePathFromEnv != "" {
		confFilePath = configFilePathFromEnv
	}

	configFile, err := os.ReadFile(confFilePath)
	if err != nil {
		log.Fatal(err)
	}
	var data Config
	err2 := yaml.Unmarshal(configFile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}
	cfg = &data

	writer2 := os.Stdout
	writer3, err := os.OpenFile(cfg.Log.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	tool.Vlog.SetOutput(io.MultiWriter(writer2, writer3))
	if cfg.Log.Level == "debug" {
		tool.Vlog.SetLevel(logrus.DebugLevel)
	} else if cfg.Log.Level == "info" {
		tool.Vlog.SetLevel(logrus.InfoLevel)
	} else if cfg.Log.Level == "error" {
		tool.Vlog.SetLevel(logrus.ErrorLevel)
	} else if cfg.Log.Level == "warn" {
		tool.Vlog.SetLevel(logrus.WarnLevel)
	}
	tool.Vlog.Error("sys config done.")
}
func Get() *Config {
	return cfg
}
