package repos

import (
	"avito/internal/db/repos"
	"avito/internal/entity"
	"context"

	"github.com/google/uuid"
)

type TenderRepo interface {
	CreateTender(ctx context.Context, tender *entity.Tender) (*entity.Tender, error)
	GetUserByUserName(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetOrgByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error)
	GetTenderByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error)
	GetTendersByFilter(ctx context.Context, filters ...repos.FilterOption) ([]entity.Tender, error)
	GetUserOrgsUUIDs(ctx context.Context, userID uuid.UUID) (uuid.UUIDs, error)
	GetOrgUsersIDsByID(ctx context.Context, id uuid.UUID) (uuid.UUIDs, error)
	UpdateTenderStatus(ctx context.Context, tenderID uuid.UUID, newStatus entity.TenderStatusType) error
	PatchTender(ctx context.Context, tenderID uuid.UUID, patchTender *entity.Tender) (*entity.Tender, error)
	RollbackTender(ctx context.Context, tenderID uuid.UUID, version int) (*entity.Tender, error)
}
