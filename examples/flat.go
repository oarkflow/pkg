package main

import (
	"fmt"
	"reflect"

	"github.com/oarkflow/pkg/konf"
)

type Config struct {
	JavaHome string `json:"java_home" env:"JAVA_HOME"`
	Lc       struct {
		Ctype string `json:"ctype"`
	} `json:"lc"`
	Login struct {
		Shell int `json:"shell,string"`
	} `json:"login"`
}

func main() {
	floatCompare()
	cfg := konf.New("config.json", true)
	err := cfg.Read()
	if err != nil {
		panic(err)
	}
	var config Config
	err = cfg.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}

func floatCompare() {
	var a, b any
	a = "1"
	b = 1
	fmt.Println(a == b)
	fmt.Println(reflect.DeepEqual(a, b))
}
