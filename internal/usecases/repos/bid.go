package repos

import (
	"avito/internal/db/repos"
	"avito/internal/entity"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BidRepo interface {
	CreateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error)
	GetBidsByFilter(ctx context.Context, filters ...repos.FilterOption) ([]entity.Bid, error)
	GetBidByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error)
	UpdateBidStatus(ctx context.Context, bidID uuid.UUID, newStatus entity.BidStatusType) error
	PatchBid(ctx context.Context, bidID uuid.UUID, patchBid *entity.Bid) (*entity.Bid, error)
	CreateFeedback(ctx context.Context, feedback *entity.BidRewiew) (*entity.BidRewiew, error)
	ShipBid(ctx context.Context, userID uuid.UUID, bidID uuid.UUID) (bool, error)
	UnshipsBid(ctx context.Context, bidID uuid.UUID) error
	RollbackBid(ctx context.Context, bidID uuid.UUID, version int) (*entity.Bid, error)
	GetFeedbacksByFilter(ctx context.Context, filters ...repos.FilterOption) ([]entity.BidRewiew, error)

	GetClear() *gorm.DB
}
