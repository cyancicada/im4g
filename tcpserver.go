package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"sync"

	"github.com/yakaa/im4g/config"
	"github.com/yakaa/im4g/handler"
	"github.com/yakaa/im4g/models"
)

func main() {
	c := LoadConfig()
	SafeUsers := new(sync.Map)
	SafeGroup := new(sync.Map)
	messageHandler := handler.NewMessageHandler(
		&models.UserChannelView{
			SafeUsers: SafeUsers,
		},
		&models.DefaultGroup{
			SafeGroup: SafeGroup,
		},
		c,
	)
	listener, err := net.Listen("tcp", c.ListenAddress) // 打开监听接口
	if err != nil {
		fmt.Println("im server start error")
	}

	defer listener.Close()
	fmt.Println("im server is wating .... on :" + c.ListenAddress)

	for {
		conn, err := listener.Accept() // 收到来自客户端发来的消息
		if err != nil {
			fmt.Println("conn fail ...")
		}
		fmt.Println(conn.RemoteAddr(), "connect successed")
		go messageHandler.Handler(conn) // 创建线程
	}
}
func LoadConfig() *config.Config {
	configFile := flag.String("f", "config/conf.json", "the config file")
	flag.Parse()
	c := new(config.Config)
	json.Unmarshal([]byte(*configFile), c)
	return c
}
