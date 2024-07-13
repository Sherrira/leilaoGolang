package auction

import (
	"context"

	"github.com/Sherrira/leilaoGolang/configuration/logger"
	"github.com/Sherrira/leilaoGolang/internal/entity/auction_entity"
	"github.com/Sherrira/leilaoGolang/internal/internal_error"
	"go.mongodb.org/mongo-driver/mongo"
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

	return nil
}
