package mongodata

import (
	"context"
	"github.com/alperhankendi/golang-api/internal/basket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	collection *mongo.Collection
}

func (receiver mongoRepository) Create(ctx context.Context, basket *basket.Basket) error {

	_, err := receiver.collection.InsertOne(ctx, basket)
	return err
}
func (receiver mongoRepository) Get(ctx context.Context, id string) (basket *basket.Basket, err error) {

	err = receiver.collection.FindOne(ctx, getId(id)).Decode(&basket)
	return
}

//SELECT * FROM baskets where Id = {id}
func getId(id string) interface{} {
	return bson.M{"id": id}
}

func NewRepository(db *mongo.Database) basket.Repository {

	col := db.Collection("baskets")

	return &mongoRepository{
		collection: col,
	}
}

func (receiver mongoRepository) Update(ctx context.Context, basket *basket.Basket) error {

	replaceOptions := options.Replace().SetUpsert(true)

	_, err := receiver.collection.ReplaceOne(ctx, getId(basket.Id), basket, replaceOptions)

	return err
}

func (receiver mongoRepository) Count(ctx context.Context) (int64, error) {

	return receiver.collection.CountDocuments(ctx, bson.M{})
}
