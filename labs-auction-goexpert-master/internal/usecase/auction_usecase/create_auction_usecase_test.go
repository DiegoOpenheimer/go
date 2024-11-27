package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/internal_error"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type AuctionRepositoryMock struct {
	mock.Mock
}

func (m *AuctionRepositoryMock) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	m.Called(ctx, auctionEntity)
	return nil
}

func (m *AuctionRepositoryMock) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	m.Called(ctx, status, category, productName)
	return []auction_entity.Auction{}, nil
}

func (m *AuctionRepositoryMock) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	args := m.Called(ctx, id)
	return args.Get(0).(*auction_entity.Auction), nil
}

func (m *AuctionRepositoryMock) Close(ctx context.Context, auctionId string) *internal_error.InternalError {
	m.Called(ctx, auctionId)
	return nil
}

type BidEntityRepositoryMock struct {
	mock.Mock
}

func (m *BidEntityRepositoryMock) CreateBid(ctx context.Context, bidEntities []bid_entity.Bid) *internal_error.InternalError {
	m.Called(ctx, bidEntities)
	return nil
}

func (m *BidEntityRepositoryMock) FindBidByAuctionId(ctx context.Context, auctionId string) ([]bid_entity.Bid, *internal_error.InternalError) {
	m.Called(ctx, auctionId)
	return []bid_entity.Bid{}, nil
}

func (m *BidEntityRepositoryMock) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*bid_entity.Bid, *internal_error.InternalError) {
	m.Called(ctx, auctionId)
	return &bid_entity.Bid{}, nil
}

type CreateAuctionUseCaseSuite struct {
	suite.Suite
	auction_entity.AuctionRepositoryInterface
	bid_entity.BidEntityRepository
}

func (suite *CreateAuctionUseCaseSuite) SetupTest() {
	suite.AuctionRepositoryInterface = &AuctionRepositoryMock{}
	suite.BidEntityRepository = &BidEntityRepositoryMock{}

	_ = os.Setenv("AUCTION_DURATION", "2s")
}

func (suite *CreateAuctionUseCaseSuite) TestExpireAuction() {
	ctx := context.Background()
	auction := AuctionInputDTO{
		ProductName: "Product",
		Category:    "Category",
		Description: "Description",
		Condition:   ProductCondition(auction_entity.New),
	}
	var id string
	suite.AuctionRepositoryInterface.(*AuctionRepositoryMock).On("CreateAuction", ctx, mock.MatchedBy(func(auctionEntity *auction_entity.Auction) bool {
		id = auctionEntity.Id
		return auctionEntity.ProductName == auction.ProductName && auctionEntity.Category == auction.Category && auctionEntity.Description == auction.Description && auctionEntity.Condition == auction_entity.ProductCondition(auction.Condition)
	})).Return(nil)
	suite.AuctionRepositoryInterface.(*AuctionRepositoryMock).On("Close", ctx, mock.MatchedBy(func(auctionId string) bool {
		return id == auctionId
	})).Return(nil)

	useCase := NewAuctionUseCase(suite.AuctionRepositoryInterface, suite.BidEntityRepository)
	err := useCase.CreateAuction(ctx, auction)
	suite.Nil(err)
	time.Sleep(3 * time.Second)
	suite.AuctionRepositoryInterface.(*AuctionRepositoryMock).AssertExpectations(suite.T())
}

func TestCreateAuctionUseCase(t *testing.T) {
	suite.Run(t, new(CreateAuctionUseCaseSuite))
}
