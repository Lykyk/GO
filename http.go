package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"math/rand"
	"time"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"		//第三方包，需要安装才能运行，例如执行下面的语句
	"github.com/silenceper/wechat" 			//go get github.com/silenceper/wechat
	"github.com/silenceper/wechat/message"
)

//以 []string 形式返回一行数据，并将指针移到下一行
func GetOne(rows *sql.Rows) []string {
	if rows == nil {
		return nil
	}

	cols, err := rows.Columns()

	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))
	dest := make([]interface{}, len(cols))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}

	if rows.Next() {
		err = rows.Scan(dest...)

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = ""
			} else {
				result[i] = string(raw)
			}
		}

		//fmt.Printf("%#v\n", result)

		//break
	} else {
		return nil
	}

	_ = err
	return result
}

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

		var split_message   []string
		var action 			string		//动作标记
		var info 			string		//附带信息
		var err_info 		string		//回复的错误信息
		var return_message 	string 		//回复的信息
		var sql_str 		string		//SQL 语句
		var rows 			*sql.Rows	//执行 SQL 语句后返回的结果
		var err 			error		//执行 SQL 语句后返回的错误
		var row 			[]string	//SQL 查询行数据
		_ = row
		var true_code 		string 		//生成的签到码
		_ = true_code


		log.Println("用户信息<-：" + msg.Content) //打印日志信息:用户发送的信息

		//连接数据库，实际上在查询时才真正连接
		dbtemp, err := sql.Open("mysql", "root:baixuewuyaak47@/wx?charset=utf8")
		db := dbtemp
		if err != nil {
			fmt.Println(err)
		}


		//以空格为标志分解用户信息
		split_message = strings.Split(msg.Content, " ")
		if len(split_message) == 1 {
			action = split_message[0]
			info = ""
		} else if len(split_message) == 2 {
			action = split_message[0]
			info = split_message[1]
		} else {
			err_info = "消息不合规范，最多只能有两段。"
			return_message = err_info
		}


		//没有错误信息则继续
		if err_info == "" {

			//生成回复信息
			return_message = "openid = " + msg.FromUserName + "\n动作标志：" + action + "\n附带信息：" + info
			
			if action == "BDXH" {
				/*绑定学号*/
				id := info     //学号
				openid := msg.FromUserName //openid
	
				//查询 student 表是否存在此学号
				sql_str = "SELECT * FROM  `student` WHERE  `id` =  '" + id + "'"
				rows, err = db.Query(sql_str)
				_ = rows
				_ = err
	
				if GetOne(rows) == nil {
					return_message += "\n\n学生名单不存在此学号"
				} else {
					//将 学号 和 openid 添加到 user 表
					sql_str = "INSERT INTO `wx`.`user` (`openid`, `id`) VALUES ('" + openid + "', '" + id + "');"
					rows, err = db.Query(sql_str)
	
					// if err != nil {
					// 	fmt.Println("db.Query(sql) 执行发生错误:", err.Error())
					// }
					
					//未发生错误表示添加成功
					if err == nil {
						return_message += "\n\n绑定学号成功。"
					} else if strings.Contains(err.Error(), "Error 1062: Duplicate entry") { //错误提示:主键重复。表示记录已经存在
						return_message += "\n\n此微信号已经绑定过"
					}
				}
			} else if action == "KSQD" && msg.FromUserName == "oT_d3096S2XEn34jDGUbbqRCf0ng" { //只有特定用户可以生成签到码
				/*开始签到,生成签到码*/
	
				rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
				true_code = fmt.Sprintf("%06v", rnd.Int31n(1000000))	//生成实际签到码
				// fmt.Println("实际签到码:", true_code)
	
				//将签到码存到表中,并记录生成时间
				sql_str = "UPDATE `acode` SET `code`='" + true_code + "', `start_time` = CURRENT_TIMESTAMP WHERE 1"
				rows, err = db.Query(sql_str)
				// if err != nil{
				// 	fmt.Println(err.Error())
				// }
				
				return_message += "\n\n签到码：" + true_code
	
			} else if action == "QD" {
				/*签到*/
				user_code := info //用户发送的签到码
				sql_str = "SELECT * FROM  `acode` LIMIT 0 , 1"
				rows, err = db.Query(sql_str)
				_ = rows
				_ = err
				row = GetOne(rows)
				true_code = row[0] //实际签到码
				
				
				//验证签到码
				if user_code == true_code {
					//将签到时间与openid写入签到表
					sql_str = "INSERT INTO `wx`.`attendence` (`openid`, `atime`) VALUES ('" + msg.FromUserName + "', CURRENT_TIMESTAMP);"
					rows, err = db.Query(sql_str)
					if err == nil {
						return_message += "\n\n签到成功"
					} else {
						return_message += "\n\n签到失败"
						fmt.Println(err.Error())
					}
				} else {
					return_message += "\n\n签到码不正确"
				}
	
			} else if action == "QQ" {
				/*缺勤*/
				sql_str = "SELECT * FROM `student` WHERE `student`.`id` IN (SELECT `user`.`id` FROM `user` WHERE `user`.`openid` NOT IN (SELECT `openid` FROM `attendence`))"
				rows, err = db.Query(sql_str)
	
				//生成缺勤名单
				if rows != nil {
					return_message += "\n\n缺勤名单："
					for {
						row := GetOne(rows)
						if row == nil {
							break
						}
						return_message += "\n" + row[0] + " " + row[1]
					}
				}
			}
		}


		// db.Close()
		log.Println("回复信息->：" + return_message) //打印日志信息:回复给用户的信息
		text := message.NewText(return_message)		//回复的信息
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
