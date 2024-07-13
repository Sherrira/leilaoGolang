package bid_usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Sherrira/leilaoGolang/configuration/logger"
	"github.com/Sherrira/leilaoGolang/internal/entity/bid_entity"
	"github.com/Sherrira/leilaoGolang/internal/internal_error"
)

type BidInputDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type BidUseCase struct {
	BidRepository bid_entity.BidEntityRepository

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bid_entity.Bid
}

func NewBidUseCase(bidRepository bid_entity.BidEntityRepository) BidUseCaseInterface {
	maxSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()
	bidUC := &BidUseCase{
		BidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxSizeInterval,
		timer:               time.NewTimer(maxSizeInterval),
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
	}

	bidUC.triggerCreateRoutine(context.Background())

	return bidUC
}

var bidBatch []bid_entity.Bid

type BidUseCaseInterface interface {
	CreateBid(ctx context.Context, bitInputDTO BidInputDTO) *internal_error.InternalError
	FindBidByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*BidOutputDTO, *internal_error.InternalError)
}

func (uc *BidUseCase) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(uc.bidChannel)

		for {
			select {
			case bidEntity, ok := <-uc.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						if err := uc.BidRepository.CreateBid(ctx, bidBatch); err != nil {
							logger.Error("Error creating bid", err)
						}
					}
					return
				}

				bidBatch = append(bidBatch, bidEntity)

				if len(bidBatch) >= uc.maxBatchSize {
					if err := uc.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("Error creating bid", err)
					}
					bidBatch = nil
					uc.timer.Reset(uc.batchInsertInterval)
				}
			case <-uc.timer.C:
				if err := uc.BidRepository.CreateBid(ctx, bidBatch); err != nil {
					logger.Error("Error creating bid", err)
				}
				bidBatch = nil
				uc.timer.Reset(uc.batchInsertInterval)
			}
		}
	}()
}

func (uc *BidUseCase) CreateBid(ctx context.Context, bitInputDTO BidInputDTO) *internal_error.InternalError {

	bid, err := bid_entity.CreateBid(bitInputDTO.UserId, bitInputDTO.AuctionId, bitInputDTO.Amount)
	if err != nil {
		return err
	}

	uc.bidChannel <- *bid

	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}

	return value
}
