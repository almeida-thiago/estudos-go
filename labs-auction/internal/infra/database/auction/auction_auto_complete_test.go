package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuctionAutoComplete(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	assert.NoError(t, err, "Error connecting to MongoDB")
	defer client.Disconnect(context.TODO())

	db := client.Database(os.Getenv("MONGODB_DB") + "_test")
	collection := db.Collection("auctions")

	_, err = collection.DeleteMany(context.TODO(), bson.M{})
	assert.NoError(t, err, "Error cleaning collection before the test")

	auction := auction_entity.Auction{
		Id:     primitive.NewObjectID().Hex(),
		Status: auction_entity.Active,
	}
	_, err = collection.InsertOne(context.TODO(), auction)
	assert.NoError(t, err, "Error inserting test auction")

	time.Sleep(getAuctionInterval() + 2*time.Second)

	var updatedAuction auction_entity.Auction
	err = collection.FindOne(context.TODO(), bson.M{"_id": auction.Id}).Decode(&updatedAuction)
	assert.NoError(t, err, "Error fetching updated auction")
	assert.Equal(t, auction_entity.Completed, updatedAuction.Status, "The auction was not automatically closed")

	_, err = collection.DeleteMany(context.TODO(), bson.M{})
	assert.NoError(t, err, "Error cleaning collection after the test")
}
