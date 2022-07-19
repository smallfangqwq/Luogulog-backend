package database

import (
	"context"
	"luogulog/declare"

	"github.com/qiniu/qmgo"
)

func ConnectDatabase(Config declare.ConfigDatabase) (*qmgo.Database, error) {
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: Config.Url})
	if err != nil {
		return nil, err
	}
	db := client.Database(Config.Name)
	return db, nil
}