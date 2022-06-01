package main

import (
	"fmt"
	"net/http"
)

func handler(writer http.ResponseWriter, request *http.Request) {
	// Pathはそのままでドメイン以下の文字列になっている localhost/XXXX /が0でXXXXXが1:
	fmt.Fprintf(writer, "Hello World Test Go!! %s", request.URL.Path[1:])
}

// ちなみにCLIでgo run main.goを実行するのはいいがCtrl + zで停止した際にプロセスが残るため、再度実行してもすぐ落ちる。
// kill -9 PIDでプロセスをkillする必要がある

// 上記がめんどくさい場合はCtrl + cで停止しよう
func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
