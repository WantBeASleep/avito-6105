package repos

import (
	"avito/internal/config"

	"avito/internal/db/models"
	"avito/internal/entity"
	"avito/internal/utils"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TenderRepo struct {
	db *gorm.DB
}

func (r *TenderRepo) GetClear() *gorm.DB { return r.db }

func NewTenderRepo(cfg *config.DB) (*TenderRepo, error) {
	log := newLogger()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.PostgresConn,
	}), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		return nil, fmt.Errorf("create db gorm obj: %w", err)
	}

	repoCtrl.initIfNeed(db)

	return &TenderRepo{
		db: db,
	}, nil
}

func (r *TenderRepo) createBackup(ctx context.Context, tender *models.Tender) error {
	backup := trnsfrm.TenderToTenderVersion(tender)

	err := createRecord(ctx, r.db, &models.TenderVersion{}, backup)

	return err
}

func (r *TenderRepo) CreateTender(ctx context.Context, tender *entity.Tender) (*entity.Tender, error) {
	tenderDB := utils.MustTransformObj[entity.Tender, models.Tender](tender)

	if err := createRecord(ctx, r.db, &models.Tender{}, tenderDB); err != nil {
		return nil, fmt.Errorf("create tender: %w", err)
	}

	if err := r.createBackup(ctx, tenderDB); err != nil {
		return nil, fmt.Errorf("create tender backup: %w", err)
	}

	return utils.MustTransformObj[models.Tender, entity.Tender](tenderDB), nil
}

func (r *TenderRepo) GetUserByUserName(ctx context.Context, username string) (*entity.User, error) {
	return getSingleMappedRecord[entity.User, models.User](ctx, r.db, entity.ErrUserNotFound, WithWhere("username = ?", username))
}

func (r *TenderRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return getSingleMappedRecord[entity.User, models.User](ctx, r.db, entity.ErrUserNotFound, WithWhere("id = ?", id))
}

func (r *TenderRepo) GetOrgByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error) {
	return getSingleMappedRecord[entity.Organization, models.Organization](ctx, r.db, entity.ErrOrgNotFound, WithWhere("id = ?", id))
}

func (r *TenderRepo) GetOrgUsersIDsByID(ctx context.Context, id uuid.UUID) (uuid.UUIDs, error) {
	orgsUsers, err := getMultiRecord(ctx, r.db, &models.OrganizationResponsible{},
		WithWhere("organization_id = ?", id),
	)
	if err != nil {
		return nil, err
	}

	return trnsfrm.OrgRespToUserUUIDSlice(orgsUsers), nil
}

func (r *TenderRepo) GetTenderByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error) {
	return getSingleMappedRecord[entity.Tender, models.Tender](ctx, r.db, entity.ErrTenderNotFound, WithWhere("id = ?", tenderID))
}

func (r *TenderRepo) GetTendersByFilter(ctx context.Context, filters ...FilterOption) ([]entity.Tender, error) {
	return getMultiMappedRecord[entity.Tender, models.Tender](ctx, r.db, filters...)
}

func (r *TenderRepo) GetUserOrgsUUIDs(ctx context.Context, userID uuid.UUID) (uuid.UUIDs, error) {
	resp, err := getMultiRecord(ctx, r.db, &models.OrganizationResponsible{}, WithWhere("user_id = ?", userID))
	if err != nil {
		return nil, err
	}

	return trnsfrm.OrgRespToOrgUUIDSlice(resp), nil
}

func (r *TenderRepo) UpdateTenderStatus(ctx context.Context, tenderID uuid.UUID, newStatus entity.TenderStatusType) error {
	queryRes := r.db.WithContext(ctx).
		Model(&models.Tender{}).
		Where("id = ?", tenderID).
		Update("status", newStatus)

	if queryRes.Error != nil {
		return queryRes.Error
	}

	if queryRes.RowsAffected == 0 {
		return entity.ErrTenderNotFound
	}

	return nil
}

func (r *TenderRepo) PatchTender(ctx context.Context, tenderID uuid.UUID, patchTender *entity.Tender) (*entity.Tender, error) {
	patchTenderDB := utils.MustTransformObj[entity.Tender, models.Tender](patchTender)

	if err := r.db.WithContext(ctx).
		Model(&models.Tender{}).
		Where("id = ?", tenderID).
		Updates(patchTenderDB).
		Error; err != nil {
		return nil, err
	}

	tenderDB, err := getSingleRecord(ctx, r.db, &models.Tender{}, WithWhere("id = ?", tenderID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrTenderNotFound
		}
		return nil, err
	}

	tenderDB.Version += 1

	if err := r.db.WithContext(ctx).
		Model(&models.Tender{}).
		Where("id = ?", tenderID).
		Update("version", tenderDB.Version).
		Error; err != nil {
		return nil, err
	}

	if err := r.createBackup(ctx, tenderDB); err != nil {
		return nil, fmt.Errorf("create backup: %w", err)
	}

	return utils.MustTransformObj[models.Tender, entity.Tender](tenderDB), nil
}

func (r *TenderRepo) RollbackTender(ctx context.Context, tenderID uuid.UUID, version int) (*entity.Tender, error) {
	currTender, err := getSingleRecord(ctx, r.db, &models.Tender{}, WithWhere("id = ?", tenderID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrTenderNotFound
		}
		return nil, err
	}

	newStatus := currTender.Status
	newVersion := currTender.Version + 1

	backupTenderDB, err := getSingleRecord(ctx, r.db, &models.TenderVersion{},
		WithWhere("tender_id = ?", tenderID),
		WithWhere("version = ?", version),
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrTenderVersionNotFound
		}
		return nil, err
	}

	rollbackTender := trnsfrm.TenderVersionToTender(backupTenderDB)
	rollbackTender.Status = newStatus
	rollbackTender.Version = newVersion

	if err := r.db.WithContext(ctx).
		Model(&models.Tender{}).
		Where("id = ?", tenderID).
		Updates(rollbackTender).
		Error; err != nil {
		return nil, err
	}

	if err := r.createBackup(ctx, rollbackTender); err != nil {
		return nil, fmt.Errorf("create backup: %w", err)
	}

	return utils.MustTransformObj[models.Tender, entity.Tender](rollbackTender), nil
}
