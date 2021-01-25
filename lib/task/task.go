package task

import (
	"context"
	// "fmt"
	"github.com/eavesmy/mongo_monitor/lib/db"
	"github.com/robertkrimen/otto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "time"
)

type Subscribe struct {
	Name   string
	Cancel context.CancelFunc
	Script string
}

type Task struct {
	db        *mongo.Client
	Subscribe map[string]*Subscribe
}

func NewTask(uri string) *Task {
	return &Task{db: db.NewClient(uri), Subscribe: map[string]*Subscribe{}}
}

// 传入参数
func (t *Task) Sub() {

	ctx, cancel := context.WithCancel(context.Background())
	sub := &Subscribe{Cancel: cancel, Name: "test", Script: "console.log(result); console.log(JSON.stringify(result))"}

	t.Subscribe[sub.Name] = sub

	// opts := options.ChangeStream().SetMaxAwaitTime(2 * time.Second)

	if err := t.db.Connect(ctx); err != nil {
		panic(err)
	}

	match := bson.D{{"$match", bson.D{{"operationType", "update"}, {"updateDescription.updatedFields.cash", bson.D{{"$exists", true}}}}}}

	stream, err := t.db.Database("test").Collection("a").Watch(context.Background(), mongo.Pipeline{match}, options.ChangeStream().SetFullDocument(options.Default))

	if err != nil {
		panic(err)
	}

	for stream.Next(context.Background()) {
		// 直接返回结果，交给 对应 task 处理

		go func() {
			b, _ := bson.Marshal(stream.Current)
			m := map[string]interface{}{}
			bson.Unmarshal(b, &m)

			vm := otto.New()
			vm.Set("result", m)
			vm.Run(sub.Script)
		}()
	}
}

func (t *Task) UnSub() {

}
