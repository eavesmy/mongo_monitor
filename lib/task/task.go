package task

import (
	"context"
	"fmt"
	"net"
	// "github.com/eavesmy/golang-lib/crypto"
	"encoding/json"
	"github.com/eavesmy/mongo_monitor/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "time"
)

type Subscribe struct {
	Name   string
	Cancel context.CancelFunc
	Stream *mongo.ChangeStream
	Ctx    context.Context
	Conn   net.Conn
}

type SubscribeReq struct {
	DB         string      `json:"db"`
	Collection string      `json:"collection"`
	Match      interface{} `json:"match"`
}

func NewSub(req *SubscribeReq) *Subscribe {
	// 获取 db

	ctx, cancel := context.WithCancel(context.Background())

	_db := db.Register(req.DB, req.Collection, ctx)

	sub := &Subscribe{Cancel: cancel, Ctx: ctx}

	fmt.Println(req.Match)

	// match := bson.D{{"operationType", "update"}, {"updateDescription.updatedFields.cash", bson.D{{"$exists", true}}}}

	stream, err := _db.Watch(ctx, mongo.Pipeline{bson.D{{"$match", req.Match}}}, options.ChangeStream().SetFullDocument(options.UpdateLookup))

	if err != nil {
		fmt.Println(err)
		return nil
	}

	sub.Stream = stream

	return sub
}

func (s *Subscribe) Listen() {

	for s.Stream.Next(s.Ctx) {

		// 直接返回结果，交给 对应 task 处理
		b, _ := bson.Marshal(s.Stream.Current)
		m := map[string]interface{}{}
		bson.Unmarshal(b, &m)
		r, _ := json.Marshal(m)

		s.Conn.Write(r)
	}
}
