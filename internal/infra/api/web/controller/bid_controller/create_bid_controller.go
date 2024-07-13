package bid_controller

import (
	"context"
	"net/http"

	"github.com/Sherrira/leilaoGolang/configuration/rest_err"
	"github.com/Sherrira/leilaoGolang/internal/infra/api/web/validation"
	"github.com/Sherrira/leilaoGolang/internal/usecase/bid_usecase"
	"github.com/gin-gonic/gin"
)

type BidController struct {
	bidUseCase bid_usecase.BidUseCaseInterface
}

func NewBidController(bidUseCase bid_usecase.BidUseCaseInterface) *BidController {
	return &BidController{
		bidUseCase: bidUseCase,
	}
}

func (u *BidController) CreateBid(c *gin.Context) {
	var bidInputDTO bid_usecase.BidInputDTO

	if err := c.ShouldBindJSON(&bidInputDTO); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	if err := u.bidUseCase.CreateBid(context.Background(), bidInputDTO); err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.Status(http.StatusCreated)
}
