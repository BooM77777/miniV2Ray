package main

import (
	"fmt"
	"log"
	"math/bits"
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

	// 载入配置文件
	globalConfig, err := parseConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(globalConfig)
	fmt.Println(bits.Len32(4))
	fmt.Println(bits.Len32(8))
	fmt.Println(bits.Len32(7))
	fmt.Println(bits.Len32(65535))

}

func printVersion() {
	log.Printf("miniV2Ray version : %s\n", version)
}
