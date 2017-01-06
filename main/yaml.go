package main

import (
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mysql map[string]string `yaml:mysql`
}

func main() {
	var config Config
	source, err := ioutil.ReadFile("D:/mygo/src/zhiwang_bc_message/cfg/config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal([]byte(source), &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf(config.Mysql["ip"])
}
