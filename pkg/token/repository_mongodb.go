package token

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

//InsertOne for audit
func (r *MongoRepository) InsertOne(tid, ownerID string) (bool, error) {
	col := r.Connection.Database("wall").Collection("TokenList")
	tokenInfo := &entity.TokenInfo{
		Tid:      tid,
		OwnerID:  ownerID,
		Enabled:  true,
		CreateAt: makeTimestamp(),
	}
	_, err := col.InsertOne(context.TODO(), tokenInfo)
	if err != nil {
		return false, err
	}
	return true, nil
}

//FindOne
func (r *MongoRepository) FindOne(tid string) (*entity.TokenInfo, error) {
	var result *entity.TokenInfo
	col := r.Connection.Database("wall").Collection("TokenList")
	err := col.FindOne(context.TODO(), bson.M{"tid": tid}).Decode(&result)
	switch err {
	case nil:
		return result, nil
	default:
		return nil, err
	}
}

//UpdateTokenStatus ...
func (r *MongoRepository) UpdateTokenStatus(tid, ownerID string) (bool, error) {
	coll := r.Connection.Database("wall").Collection("TokenList")

	update := bson.M{
		"$set": bson.M{
			"enabled": false,
		},
	}
	_, err := coll.UpdateOne(
		context.TODO(),
		bson.M{
			"tid":     tid,
			"ownerID": ownerID,
		},
		update,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
