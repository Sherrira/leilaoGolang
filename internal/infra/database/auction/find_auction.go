package auction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Sherrira/leilaoGolang/configuration/logger"
	"github.com/Sherrira/leilaoGolang/internal/entity/auction_entity"
	"github.com/Sherrira/leilaoGolang/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ar *AuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var auctionEntityMongo AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			msg := fmt.Sprintf("auction with id %s not found", id)
			logger.Error(msg, err)
			return nil, internal_error.NewNotFoundError(msg)
		}

		msg := "error finding auction by id"
		logger.Error(msg, err)
		return nil, internal_error.NewInternalServerError(msg)
	}

	return &auction_entity.Auction{
		Id:          auctionEntityMongo.Id,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
	}, nil
}

func (ar *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {

	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["product_name"] = primitive.Regex{Pattern: productName, Options: "i"}
	}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		msg := "error finding auctions"
		logger.Error(msg, err)
		return nil, internal_error.NewInternalServerError(msg)
	}
	defer cursor.Close(ctx)

	var auctionEntityMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionEntityMongo); err != nil {
		msg := "error decoding auctions"
		logger.Error(msg, err)
		return nil, internal_error.NewInternalServerError(msg)
	}

	var auctionEntity []auction_entity.Auction
	for _, a := range auctionEntityMongo {
		auctionEntity = append(auctionEntity, auction_entity.Auction{
			Id:          a.Id,
			ProductName: a.ProductName,
			Category:    a.Category,
			Description: a.Description,
			Condition:   a.Condition,
			Status:      a.Status,
			Timestamp:   time.Unix(a.Timestamp, 0),
		})
	}

	return auctionEntity, nil
}
