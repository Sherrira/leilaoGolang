package auction

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Sherrira/leilaoGolang/configuration/database/mongodb"
	"github.com/Sherrira/leilaoGolang/internal/entity/auction_entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuctionIntegration(t *testing.T) {
	const (
		AUCTION_INTERVAL           = "3s"
		TIME_SLEEP_TO_FIND_AUCTION = 5 * time.Second
	)

	ctx := context.Background()
	os.Setenv("AUCTION_INTERVAL", AUCTION_INTERVAL)

	database, err := mongodb.NewMongoDBConnectionIntegrationTest(ctx)
	if err != nil {
		log.Fatal(`Integration Test Error - No MongoDB connection. \
		Check if MongoDB is running. \
		Try 'make run-mongo' to start MongoDB.\
		Error: `, err)
		return
	}
	defer database.Drop(ctx)

	auctionRepo := NewAuctionRepository(database)

	testAuction := &auction_entity.Auction{
		Id:          "testId",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	err = auctionRepo.CreateAuction(ctx, testAuction)
	assert.Nil(t, err)

	time.Sleep(TIME_SLEEP_TO_FIND_AUCTION)

	result, err := auctionRepo.FindAuctionById(ctx, "testId")
	assert.Nil(t, err)
	assert.Equal(t, auction_entity.Completed, result.Status)
}
