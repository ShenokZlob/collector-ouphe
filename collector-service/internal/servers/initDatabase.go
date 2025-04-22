package servers

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitDataBase(config *viper.Viper) *mongo.Client {
	connString := config.GetString("database.conn_string")
	client, err := mongo.Connect(options.Client().ApplyURI(connString))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	return client
}
