package bid_entity

import (
	"context"
	"time"

	"github.com/Sherrira/leilaoGolang/internal/internal_error"
	"github.com/google/uuid"
)

type Bid struct {
	Id        string
	UserId    string
	AuctionId string
	Amount    float64
	Timestamp time.Time
}

func CreateBid(userId, auctionId string, amount float64) (*Bid, *internal_error.InternalError) {
	bid := &Bid{
		Id:        uuid.New().String(),
		UserId:    userId,
		AuctionId: auctionId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.Validate(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (b *Bid) Validate() *internal_error.InternalError {
	if err := uuid.Validate(b.UserId); err != nil {
		return internal_error.NewBadRequestError("Invalid user id")
	}

	if err := uuid.Validate(b.AuctionId); err != nil {
		return internal_error.NewBadRequestError("Invalid auction id")
	}

	if b.Amount <= 0 {
		return internal_error.NewBadRequestError("Invalid amount")
	}

	return nil
}

type BidEntityRepository interface {
	CreateBid(ctx context.Context, bidEntities []Bid) *internal_error.InternalError
	FindBidByAuctionId(ctx context.Context, auctionId string) ([]Bid, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*Bid, *internal_error.InternalError)
}
