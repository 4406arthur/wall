package developer

import (
	"context"
	"log"
	"time"

	"wall/pkg/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoRepository mongodb repo
type MongoRepository struct {
	Connection *mongo.Client
}

//NewMongoRepository create new repository
func NewMongoRepository(mongoEndpoint string) *MongoRepository {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoCli, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		fatal(err)
	}
	return &MongoRepository{
		Connection: mongoCli,
	}
}

//Find a Developer
func (r *MongoRepository) Find(name string) (*entity.Developer, error) {
	var result entity.Developer
	col := r.Connection.Database("ionian").Collection("Developer")
	err := col.FindOne(context.TODO(), bson.M{"name": name}).Decode(&result)
	switch err {
	case nil:
		return &result, nil
	default:
		return nil, err
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
