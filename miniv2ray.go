package main

import (
	"log"
	"os"

	"./vmessTest/client"
	"./vmessTest/server"
)

const (
	version      = "v1.0.0"
	proxyPortTCP = 1080
	configFile   = "./config/config.json"
)

var logger *log.Logger

func main() {

	// 打印版本号
	printVersion()

	if len(os.Args) > 1 {
		if os.Args[1] == "-c" {
			client.Client()
		} else if os.Args[1] == "-s" {
			server.Server()
		}
	}

	// 载入配置文件
	// globalConfig, err := parseConfigFile(configFile)

}

func printVersion() {
	log.Printf("miniV2Ray version : %s\n", version)
}
