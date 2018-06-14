package main

import(
	"fmt"
)

func main(){
	var s1 string
	var s2 string
	s1 = "消息不合规范，最多只能有两段。"
	s2 = s1

	if s1 == s2 {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
	
}