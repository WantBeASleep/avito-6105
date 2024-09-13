package usecases

import (
	"avito/internal/entity"
	"context"

	"github.com/google/uuid"
)

type TenderUsecase interface {
	CreateTender(ctx context.Context, username string, tender *entity.Tender) (*entity.Tender, error)

	GetTenders(ctx context.Context, serviceTypes []entity.TenderServiceType, pag *entity.Pagination) ([]entity.Tender, error)
	GetMyTenders(ctx context.Context, username string, pag *entity.Pagination) ([]entity.Tender, error)
	GetTenderStatus(ctx context.Context, username string, tenderID uuid.UUID) (entity.TenderStatusType, error)

	UpdateTenderStatus(ctx context.Context, username string, tenderID uuid.UUID, status entity.TenderStatusType) (*entity.Tender, error)
	PatchTender(ctx context.Context, username string, tenderID uuid.UUID, patchTender *entity.Tender) (*entity.Tender, error)

	RollbackTender(ctx context.Context, username string, tenderID uuid.UUID, version int) (*entity.Tender, error)
}
