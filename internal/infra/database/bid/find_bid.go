package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/Sherrira/leilaoGolang/configuration/logger"
	"github.com/Sherrira/leilaoGolang/internal/entity/bid_entity"
	"github.com/Sherrira/leilaoGolang/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *BidRepository) FindBidByAuctionId(ctx context.Context, auctionId string) ([]bid_entity.Bid, *internal_error.InternalError) {
	errMsg := fmt.Sprintf("error finding bid by auction id %s", auctionId)

	cursor, err := b.Collection.Find(ctx, bson.M{"auction_id": auctionId})
	if err != nil {
		logger.Error(errMsg, err)
		return nil, internal_error.NewInternalServerError(errMsg)
	}

	var bids []BidEntityMongo
	if err := cursor.All(ctx, &bids); err != nil {
		logger.Error(errMsg, err)
		return nil, internal_error.NewInternalServerError(errMsg)
	}

	var bidEntities []bid_entity.Bid
	for _, bid := range bids {
		bidEntities = append(bidEntities, bid_entity.Bid{
			Id:        bid.Id,
			UserId:    bid.UserId,
			AuctionId: bid.AuctionId,
			Amount:    bid.Amount,
			Timestamp: time.Unix(bid.Timestamp, 0),
		})
	}

	return bidEntities, nil
}

func (b *BidRepository) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*bid_entity.Bid, *internal_error.InternalError) {
	errMsg := fmt.Sprintf("error finding winning bid by auction id %s", auctionId)

	var bidEntityMongo BidEntityMongo
	filter := bson.M{"auction_id": auctionId}
	opts := options.FindOne().SetSort(bson.D{{"amount", -1}})
	if err := b.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo); err != nil {
		logger.Error(errMsg, err)
		return nil, internal_error.NewInternalServerError(errMsg)
	}

	return &bid_entity.Bid{
		Id:        bidEntityMongo.Id,
		UserId:    bidEntityMongo.UserId,
		AuctionId: bidEntityMongo.AuctionId,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}, nil
}
