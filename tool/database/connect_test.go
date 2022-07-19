package database

import (
	"context"
	"fmt"
	"luogulog/declare"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
)

type test struct {
	PostID string
}

func TestDatabase(t *testing.T) {
	ConfingData := viper.New()
	ConfingData.SetConfigFile("../../config.yaml")
	var Configs declare.Config
	ConfingData.ReadInConfig()
	ConfingData.Unmarshal(&Configs)
	t.Log(Configs.Database.Url)
	cc, err := ConnectDatabase(Configs.Database)
	if err != nil {
		t.FailNow()
	}
//	fmt.Print()
	wow := bson.M{"PostID": "hi"}
	result, err := cc.Collection("test").InsertOne(context.Background(), wow)
	if err != nil {
		fmt.Print(err)
		t.FailNow()
	}
	fmt.Print(result)
}