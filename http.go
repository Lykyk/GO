package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/silenceper/wechat" //go get github.com/silenceper/wechat 执行后才能运行
	"github.com/silenceper/wechat/message"
)

//以 []string 形式返回一行数据，并将指针移到下一行
// func GetOne(rows *sql.Rows) []string {
// 	if rows == nil {
// 		return nil
// 	}

// 	cols, err := rows.Columns()

// 	rawResult := make([][]byte, len(cols))
// 	result := make([]string, len(cols))
// 	dest := make([]interface{}, len(cols))
// 	for i, _ := range rawResult {
// 		dest[i] = &rawResult[i]
// 	}

// 	if rows.Next() {
// 		err = rows.Scan(dest...)

// 		for i, raw := range rawResult {
// 			if raw == nil {
// 				result[i] = ""
// 			} else {
// 				result[i] = string(raw)
// 			}
// 		}

// 		//fmt.Printf("%#v\n", result)

// 		//break
// 	} else {
// 		return nil
// 	}

// 	_ = err
// 	return result
// }

func hello(rw http.ResponseWriter, req *http.Request) {

	//配置微信参数
	config := &wechat.Config{
		AppID: "wx8781b9f815a3ca4c",
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
		// return_message := "你好呀"

		// var msg_content string
		var split_message []string
		var return_message string

		// log.Println("openid =", msg.FromUserName)		//打印日志信息
		log.Println("用户信息<-：" + msg.Content) //打印日志信息:用户发送的信息

		dbtemp, err := sql.Open("mysql", "root:baixuewuyaak47@/wx?charset=utf8")
		db := dbtemp
		if err != nil {
			fmt.Println(err)
		}

		//查询
		// rows, err := db.Query("select * from student")
		// fmt.Println(err)
		// for i:=0; i<30; i++{
		// 	row := GetOne(rows)
		// 	log.Println(row[0] + " " + row[1])	//输出
		// }

		// rows, err := db.Query("select * from student")
		// fmt.Println(err)
		// _ = rows
		// _ = err

		// msg_content = string(msg.Content)
		split_message = strings.Split(msg.Content, " ") //以空格为标志分解用户信息
		//生成回复信息
		return_message = "openid = " + msg.FromUserName + "\n动作标志：" + split_message[0] + "\n附带信息：" + split_message[1]

		if split_message[0] == "BDXH" {
			id := split_message[1]     //学号
			openid := msg.FromUserName //openid

			sql := "INSERT INTO `wx`.`user` (`id`, `openid`) VALUES ('" + id + "', '" + openid + "');"
			rows, err := db.Query(sql)
			_ = rows
			_ = err
			// if err != nil {
			// 	fmt.Println("db.Query(sql) 执行发生错误:", err.Error())
			// }

			if err == nil {
				return_message += "\n\n绑定学号成功。"
			}else if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
				return_message += "\n\n此学号已经绑定过"
			}
		}

		// db.Close()
		log.Println("回复信息->：" + return_message) //打印日志信息:回复给用户的信息
		text := message.NewText(return_message)
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
