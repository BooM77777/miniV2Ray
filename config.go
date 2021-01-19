package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//GlobalConfig 保存全局配置，包含代理所需要的相关信息
type GlobalConfig struct {
	config struct {
		LocalhostConfig struct {
			LocalProxy []struct {
				Protocol string `json:"protocol"`
				Port     uint16 `json:"port"`
			} `json:"proxy"`
		} `json:"local"`
		RemoteServerConfig struct {
			Protocol string `json:"protocol"`
			Address  string `json:"address"`
			Port     uint16 `json:"port"`
			UserID   string `json:"userID"`
			AlterID  uint8  `json:"alterID"`
			Level    uint8  `json:"level"`
			Security string `json:"Security"`
			PAC      bool   `json:"pac"`
		} `json:"remote"`
	}
}

func parseConfigFile(filePath string) (globalConfig GlobalConfig, err error) {

	configFile, err := os.Open(filePath)
	defer configFile.Close() // 函数退出时关闭

	if err != nil {
		return
	}
	jsonContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		return
	}

	// 读取成功，解析json
	err = json.Unmarshal([]byte(jsonContent), &globalConfig.config)
	if err != nil {
		return
	}

	return globalConfig, nil
}

// GetRemoteServerURL 获取string类型的URL
// 例如：vmess://a684455c-b14f-11ea-bf0d-42010aaa0003@127.0.0.1:9527?alterID=4
func (c *GlobalConfig) GetRemoteServerURL() string {
	return fmt.Sprintf("%s://%s@%s:%d?alterID=%d", c.config.RemoteServerConfig.Protocol, c.config.RemoteServerConfig.UserID, c.config.RemoteServerConfig.Address, c.config.RemoteServerConfig.Port, c.config.RemoteServerConfig.AlterID)
}

// GetLocalListenerURL 获取本地监听的URL，由于存在多种协议，返回的是注册过的每一个协议对应的URL
// 例如："socks5://127.0.0.1:1081"
func (c *GlobalConfig) GetLocalListenerURL() ([]string, error) {
	format := "%s://127.0.0.1:%d"
	ret := make([]string, len(c.config.LocalhostConfig.LocalProxy))
	for i, p := range c.config.LocalhostConfig.LocalProxy {
		ret[i] = fmt.Sprintf(format, strings.ToLower(p.Protocol), p.Port)
	}
	if len(ret) == 0 {
		return nil, errors.New("没有设置程序的本地监听端口")
	}
	return ret, nil
}
