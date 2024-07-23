package auction

import (
	"context"
	"os"
	"time"

	"github.com/Sherrira/leilaoGolang/configuration/logger"
	"github.com/Sherrira/leilaoGolang/internal/entity/auction_entity"
	"github.com/Sherrira/leilaoGolang/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError {
	auctionMongo := &AuctionEntityMongo{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   auction.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionMongo)
	if err != nil {
		msg := "error creating auction"
		logger.Error(msg, err)
		return internal_error.NewInternalServerError(msg)
	}

	interval := ar.GetAuctionInterval()

	go func(auctionId string) {
		select {
		case <-time.After(interval):
			filter := bson.M{"_id": auctionId}
			update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
			_, err := ar.Collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(false))
			if err != nil {
				logger.Error("Failed to close auction", err)
			}
		case <-ctx.Done():
			return
		}
	}(auction.Id)

	return nil
}

func (ar *AuctionRepository) GetAuctionInterval() time.Duration {
	durationStr := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Error("error parsing duration, returning default...", err)
		return time.Minute * 5
	}
	return duration
}
