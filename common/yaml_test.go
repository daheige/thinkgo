package common

import (
	"log"
	"testing"
)

func TestYaml(t *testing.T) {
	conf := NewConf()
	conf.LoadConf("test.yaml")
	log.Println(conf.data)

	// var v interface{}
	// conf.GetStruct("RedisCommon", v)
	var v = &RedisConf{}
	conf.GetStruct("RedisCommon", v)
	log.Println(v)
	log.Println(conf.data["RedisCommon"])

}
