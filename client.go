// one sever to more client chat room
// This is chat client
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type (
	Message struct {
		FromUID string `json:"fromUid"`
		ToUID   string `json:"toUid"`
		Role    string `json:"role,optional"`
		Type    string `json:"type"`
		Content string `json:"content"`
		RoomID  string `json:"roomId"`
	}
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		log.Fatal("args length is 2")
	}
	conn, err := net.Dial("tcp", "127.0.0.1:7272") // 打开监听端口
	if err != nil {
		fmt.Println("conn fail...")
	}
	log.Println("client start...")
	defer conn.Close()
	var msg string
	toUid := args[2]
	fromUid := args[1]
	log.Println(toUid, fromUid)
	msg = `{"fromUid":"` + fromUid + `","toUid":"` + toUid + `","role":"s","type":"LOGIN","content":"hello LOGIN"}`
	conn.Write([]byte(msg + "\n")) // 将信息发送给服务器端

	go Handle(conn) // 创建线程

	for {
		msg = `{"fromUid":"` + fromUid + `","toUid":"` + toUid + `","role":"s","type":"SAY","content":"hello say"}`
		log.Println(msg)
		time.Sleep(5 * time.Second)
		conn.Write([]byte(msg + "\n")) // 三段字节流 say | 昵称 | 发送的消息
	}
}

func Handle(conn net.Conn) {

	for {

		data := make([]byte, 2048)       // 创建一个字节流
		msg_read, err := conn.Read(data) // 将读取的字节流赋值给msg_read和err
		if msg_read == 0 || err != nil { // 如果字节流为0或者有错误
			log.Println("no data input from server...")
		}
		m := new(Message)
		json.Unmarshal(data[0:msg_read], m)
		log.Println("I HAS GET THE MESSAGE FROM SERVER IS :" + string(data[0:msg_read])) // 把字节流转换成字符串
		log.Println(m.FromUID + " SAY to " + m.ToUID + ":" + m.Content)
	}
}
