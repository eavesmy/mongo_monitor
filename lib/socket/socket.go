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
	WaitTime   int
	DisConnect bool
}

func (s *Socket) Listen() {
	reader := bufio.NewReader(s.Conn)

	for {
		b, err := reader.ReadBytes('\n')

		if err != nil {

			if s.DisConnect {
				return
			}

			fmt.Println("connect error ", err)
			time.Sleep(5 * time.Second)
			s.WaitTime++
			if s.WaitTime > 5 {
				s.Destory()
			}
			continue
		}

		s.WaitTime = 0

		req := &task.SubscribeReq{}
		err = json.Unmarshal(b, &req)

		if err != nil {
			fmt.Println("invalid struct", err)
			s.Conn.Write([]byte("invalid struct"))
			continue
		}

		sub := task.NewSub(req)

		if sub == nil {
			// error handler
			s.Conn.Write([]byte("connect failed"))
			defer s.Destory()
			return
		}

		sub.Conn = s.Conn
		if s.Subscribes[req.DB+"."+req.Collection] == nil {
			s.Subscribes[req.DB+"."+req.Collection] = sub
		}
		// 存入
		go sub.Listen()
	}
}

func (s *Socket) Destory() {
	s.Conn.Close()
	s.DisConnect = true
	fmt.Println("destory connection")
}
