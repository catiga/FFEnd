package config

import (
	"fmt"
	"testing"
)

func TestConfigRead(t *testing.T) {
	conf := Get()
	fmt.Println(conf)
}
