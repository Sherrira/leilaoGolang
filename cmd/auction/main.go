package main

import (
	"context"
	"log"

	"github.com/Sherrira/leilaoGolang/configuration/database/mongodb"
	"github.com/Sherrira/leilaoGolang/internal/infra/api/web/controller/auction_controller"
	"github.com/Sherrira/leilaoGolang/internal/infra/api/web/controller/bid_controller"
	"github.com/Sherrira/leilaoGolang/internal/infra/api/web/controller/user_controller"
	"github.com/Sherrira/leilaoGolang/internal/infra/database/auction"
	"github.com/Sherrira/leilaoGolang/internal/infra/database/bid"
	"github.com/Sherrira/leilaoGolang/internal/infra/database/user"
	"github.com/Sherrira/leilaoGolang/internal/usecase/auction_usecase"
	"github.com/Sherrira/leilaoGolang/internal/usecase/bid_usecase"
	"github.com/Sherrira/leilaoGolang/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error loading env variables")
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
		return
	}

	router := gin.Default()

	userController, bidController, auctionController := initDependencies(databaseConnection)

	router.GET("/auctions", auctionController.FindAuctions)
	router.GET("/auctions/winner/:auctionId", auctionController.FindAuctionById)
	router.POST("/auctions", auctionController.CreateAuction)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController) {
	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(user_usecase.NewUserUseCase(userRepository))
	auctionController = auction_controller.NewAuctionController(auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository))
	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}
