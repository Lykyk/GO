package main

import (
	"fmt"
	"net/http"
	"log"

	"github.com/silenceper/wechat"	//go get github.com/silenceper/wechat 执行后才能运行
	"github.com/silenceper/wechat/message"
)

func hello(rw http.ResponseWriter, req *http.Request) {

	//配置微信参数
	config := &wechat.Config{
		AppID:          "wx8781b9f815a3ca4c",
		// AppSecret:      "your app secret",
		Token:          "baixuewuyaak47",
		EncodingAESKey: "YOqOkJTgYeEtPHC77hRpCiE3qiLUa4cjskst60QHlsu",
	}
	wc := wechat.NewWechat(config)

	// 传入request和responseWriter
	server := wc.GetServer(req, rw)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {

		//回复消息：演示回复用户发送的消息
		// text := message.NewText(msg.Content)

		text := message.NewText("你好呀")
		log.Println("用户信息：" + msg.Content)	//打印日志信息
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
