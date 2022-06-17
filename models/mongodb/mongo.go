package mongodb

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var client *mongo.Client

func InitDB() {
	//设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	//连接到MongoDB
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Info("err = ", err)
		return
	}
	//检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Info("err = ", err)
		return
	}
	fmt.Println("Connected to MongoDB!")
 	s:=Student{
		Name: "小明",
		Age:  10,
	}
	insertOne(s)
	find()
}

type Student struct {
	Name string
	Age  int
}

//单条插入
func insertOne(s Student) {
	collection := client.Database("go_db").Collection("student")
	insertResult, err := collection.InsertOne(context.TODO(), s)
	if err != nil {
		logrus.Info("err = ", err)
		return
	}
	fmt.Println("Insert a single document: ", insertResult.InsertedID)
}

//多条插入
func insertMore(students []interface{}) {
	collection := client.Database("go_db").Collection("student")
	insertManyResult, err := collection.InsertMany(context.TODO(), students)
	if err != nil {
		logrus.Info("err = ", err)
		return
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

//查询
func find() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database("go_db").Collection("student")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		logrus.Info("err = ", err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			logrus.Info("err = ", err)
		}
		fmt.Printf("result:%v\n", result)
		fmt.Sprintf("result.Map(): %v\n", result.Map()["name"])
	}
	if err := cur.Err(); err != nil {
		logrus.Info("err = ", err)
	}
}

//更新
func update() {
	ctx := context.TODO()
	defer client.Disconnect(ctx)
	c := client.Database("go_db").Collection("Student")
	update := bson.D{{"$set", bson.D{{"Name", "big tom"}, {"Age", 22}}}}
	ur, err := c.UpdateMany(ctx, bson.D{{"name", "tom"}}, update)
	if err != nil {
		logrus.Info("err = ", err)
		return
	}
	fmt.Printf("ur.ModifiedCounnt: %v\n", ur.ModifiedCount)
}

//删除
func del() {
	c := client.Database("go_db").Collection("Student")
	ctx := context.TODO()

	dr, err := c.DeleteMany(ctx, bson.D{{"Name", "big kite"}})
	if err != nil {
		logrus.Info("err = ", err)
	}
	fmt.Printf("ur.ModifiedCount: %v\n", dr.DeletedCount)
}
