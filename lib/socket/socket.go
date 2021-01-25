package socket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/eavesmy/mongo_monitor/lib/task"
	"net"
	"time"
)

type Socket struct {
	Conn       net.Conn
	Subscribes map[string]*task.Subscribe
}

func (s *Socket) Listen() {
	reader := bufio.NewReader(s.Conn)

	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("connect error ", err)
			time.Sleep(5 * time.Second)
			continue
		}

		req := &task.SubscribeReq{}
		json.Unmarshal(b, &req)

		sub := task.NewSub(req)
		sub.Conn = s.Conn
		// 存入
	}
}
