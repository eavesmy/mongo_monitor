package task

import (
	"context"
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
	DB         string `json:"db"`
	Collection string `json:"collection"`
	Match      bson.D `json:"match"`

	// match := bson.D{{"operationType", "update"}, {"updateDescription.updatedFields.cash", bson.D{{"$exists", true}}}}
}

func NewSub(req *SubscribeReq) *Subscribe {
	// 获取 db

	_db := db.Register(req.DB, req.Collection)

	ctx, cancel := context.WithCancel(context.Background())

	sub := &Subscribe{Cancel: cancel, Ctx: ctx}

	stream, err := _db.Watch(ctx, mongo.Pipeline{bson.D{{"$match", req.Match}}}, options.ChangeStream().SetFullDocument(options.UpdateLookup))

	if err != nil {
		return nil
	}

	sub.Stream = stream

	return sub
}

func (s *Subscribe) Listen() {
	for s.Stream.Next(context.Background()) {
		// 直接返回结果，交给 对应 task 处理
		b, _ := bson.Marshal(s.Stream.Current)
		m := map[string]interface{}{}
		bson.Unmarshal(b, &m)
		r, _ := json.Marshal(m)
		s.Conn.Write(r)
	}
}

/*
// 传入参数
func (t *Task) Sub(req *SubscribeReq) {

	ctx, cancel := context.WithCancel(context.Background())
	sub := &Subscribe{Cancel: cancel, Name: "test"}

	t.Subscribe[sub.Name] = sub

	// opts := options.ChangeStream().SetMaxAwaitTime(2 * time.Second)

	if err := t.db.Connect(ctx); err != nil {
		panic(err)
	}

	stream, err := t.db.Database("test").Collection("a").Watch(context.Background(), mongo.Pipeline{bson.D{{"$match", req.Match}}}, options.ChangeStream().SetFullDocument(options.UpdateLookup))

	if err != nil {
		panic(err)
	}

	sub.Stream = *stream

	go sub.Listen()
}

func (t *Task) UnSub() {

}

*/
