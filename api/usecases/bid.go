package usecases

import (
	"avito/internal/entity"
	"context"

	"github.com/google/uuid"
)

type BidUsecase interface {
	CreateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error)
	GetMyBids(ctx context.Context, username string, pag *entity.Pagination) ([]entity.Bid, error)
	GetTenderBidsList(ctx context.Context, username string, tenderID uuid.UUID, pag *entity.Pagination) ([]entity.Bid, error)
	GetBidStatus(ctx context.Context, username string, bidID uuid.UUID) (entity.BidStatusType, error)
	UpdateBidStatus(ctx context.Context, username string, bidID uuid.UUID, newStatus entity.BidStatusType) (*entity.Bid, error)
	PatchBid(ctx context.Context, username string, bidID uuid.UUID, bid *entity.Bid) (*entity.Bid, error)
	SubmitDecision(ctx context.Context, username string, bidID uuid.UUID, decision entity.BidDecisionType) (*entity.Bid, error)
	FeedbackBid(ctx context.Context, username string, bidID uuid.UUID, bidFeedback string) (*entity.Bid, error)
	RollbackBid(ctx context.Context, username string, bidID uuid.UUID, version int) (*entity.Bid, error)
	CheckPrevFeedbacks(ctx context.Context, tenderID uuid.UUID, author string, requester string, pagination entity.Pagination) ([]entity.BidRewiew, error)
}
