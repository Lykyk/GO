package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func main() {
	// http.HandleFunc("/", sayhelloName)     //设置访问的路由
	http.HandleFunc("/", checkSignature)   //设置访问的路由
	err := http.ListenAndServe(":80", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func str2sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func checkSignature(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var token string = "baixuewuyaak47"
	var signature string = strings.Join(r.Form["signature"], "")
	var timestamp string = strings.Join(r.Form["timestamp"], "")
	var nonce string = strings.Join(r.Form["nonce"], "")
	var echostr string = strings.Join(r.Form["echostr"], "")
	tmps := []string{token, timestamp, nonce}
	sort.Strings(tmps)
	tmpStr := tmps[0] + tmps[1] + tmps[2]
	tmp := str2sha1(tmpStr)
	if tmp == signature {
		log.Println("signature is OK")
		fmt.Fprintf(w, echostr)
	} else{
		log.Println("signature is fail")
	}

}
