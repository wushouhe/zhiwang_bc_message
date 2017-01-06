package main

import (
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
	"zhiwang_bc_message/geth/config"
)

func main() {
	var config config.Config
	source, err := ioutil.ReadFile("D:/mygo/src/zhiwang_bc_message/cfg/config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal([]byte(source), &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.Mysql.Username)
	fmt.Println(config.RPC.Protocol)
	fmt.Println(config.ThreadSize)
}
