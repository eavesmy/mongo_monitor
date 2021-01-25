package main

/*
	Support websocket && unix socket .
*/

import (
	"bufio"
	"fmt"
	"github.com/eavesmy/mongo_monitor/lib/task"
	"net"
	"net/url"
	"os"
	"time"
)

const SOCKET_FILE = "/tmp/mongo_watch"

var MONGOURI = "mongodb://10.40.126.223:27017/test?compressors=disabled&gssapiServiceName=mongodb"
var Task = task.NewTask(MONGOURI)

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

	l, err := net.Listen("unix", SOCKET_FILE)
	if err != nil {
		panic(err)
	}

	fmt.Println("Mongo watch start at", l.Addr())

	Task.Sub()

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go readSocket(conn)
	}
}

func readSocket(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("connect error ", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println(b)
	}
}
