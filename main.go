package main

/*
	Support websocket && unix socket .
*/

import (
	"fmt"
	"github.com/eavesmy/mongo_monitor/lib/db"
	"github.com/eavesmy/mongo_monitor/lib/socket"
	"net"
	"net/url"
	"os"
)

const SOCKET_FILE = "/tmp/mongo_watch"

var MONGOURI = "mongodb://10.40.126.223:27017/test?compressors=disabled&gssapiServiceName=mongodb"

// var Task = task.NewTask(MONGOURI)

// 启动 Unix Socket 服务。 其他服务向该服务请求订阅信息，将收到的订阅信息发送至订阅方。

func main() {

	if _, err := os.Stat(SOCKET_FILE); err == nil || os.IsExist(err) {
		os.Remove(SOCKET_FILE)
	}

	if len(os.Args) >= 2 {
		MONGOURI = os.Args[1]
	}

	_, err := url.ParseRequestURI(MONGOURI)
	if err != nil {
		panic("Invalid mongo uri.")
	}

	db.URI = MONGOURI

	l, err := net.Listen("unix", SOCKET_FILE)
	if err != nil {
		panic(err)
	}

	fmt.Println("Mongo watch start at", l.Addr())

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		socket := &socket.Socket{Conn: conn}
		go socket.Listen()
	}
}
