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

type BidRepo struct {
	db *gorm.DB
}

func (r *BidRepo) GetClear() *gorm.DB { return r.db }

func NewBidRepo(cfg *config.DB) (*BidRepo, error) {
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

	return &BidRepo{
		db: db,
	}, nil
}

func (r *BidRepo) createBackup(ctx context.Context, bid *models.Bid) error {
	backup := trnsfrm.BidToBidVersion(bid)

	err := createRecord(ctx, r.db, &models.BidVersion{}, backup)

	return err
}

func (r *BidRepo) CreateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error) {
	bidDB := utils.MustTransformObj[entity.Bid, models.Bid](bid)

	if err := createRecord(ctx, r.db, &models.Bid{}, bidDB); err != nil {
		return nil, fmt.Errorf("create bid: %w", err)
	}

	if err := r.createBackup(ctx, bidDB); err != nil {
		return nil, fmt.Errorf("create bid backup: %w", err)
	}

	return utils.MustTransformObj[models.Bid, entity.Bid](bidDB), nil
}

func (r *BidRepo) GetBidsByFilter(ctx context.Context, filters ...FilterOption) ([]entity.Bid, error) {
	return getMultiMappedRecord[entity.Bid, models.Bid](ctx, r.db, filters...)
}

func (r *BidRepo) GetBidByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error) {
	return getSingleMappedRecord[entity.Bid, models.Bid](ctx, r.db, entity.ErrBidNotFound, WithWhere("id = ?", bidID))
}

func (r *BidRepo) UpdateBidStatus(ctx context.Context, bidID uuid.UUID, newStatus entity.BidStatusType) error {
	queryRes := r.db.WithContext(ctx).
		Model(&models.Bid{}).
		Where("id = ?", bidID).
		Update("status", newStatus)

	if queryRes.Error != nil {
		return queryRes.Error
	}

	if queryRes.RowsAffected == 0 {
		return entity.ErrBidNotFound
	}

	return nil
}

func (r *BidRepo) PatchBid(ctx context.Context, bidID uuid.UUID, patchBid *entity.Bid) (*entity.Bid, error) {
	bidDB := utils.MustTransformObj[entity.Bid, models.Bid](patchBid)

	if err := r.db.WithContext(ctx).
		Model(&models.Bid{}).
		Where("id = ?", bidID).
		Updates(bidDB).
		Error; err != nil {
		return nil, err
	}

	bidDB, err := getSingleRecord(ctx, r.db, &models.Bid{}, WithWhere("id = ?", bidID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrBidNotFound
		}
		return nil, err
	}
	bidDB.Version += 1

	if err := r.db.WithContext(ctx).
		Model(&models.Bid{}).
		Where("id = ?", bidID).
		Update("version", bidDB.Version).
		Error; err != nil {
		return nil, err
	}

	if err := r.createBackup(ctx, bidDB); err != nil {
		return nil, fmt.Errorf("create bid backup: %w", err)
	}

	return utils.MustTransformObj[models.Bid, entity.Bid](bidDB), nil
}

func (r *BidRepo) CreateFeedback(ctx context.Context, feedback *entity.BidRewiew) (*entity.BidRewiew, error) {
	rewiewDB := utils.MustTransformObj[entity.BidRewiew, models.BidRewiew](feedback)

	if err := createRecord(ctx, r.db, &models.BidRewiew{}, rewiewDB); err != nil {
		return nil, fmt.Errorf("create rewiew: %w", err)
	}

	return utils.MustTransformObj[models.BidRewiew, entity.BidRewiew](rewiewDB), nil
}

func (r *BidRepo) RollbackBid(ctx context.Context, bidID uuid.UUID, version int) (*entity.Bid, error) {
	currBid, err := getSingleRecord(ctx, r.db, &models.Bid{}, WithWhere("id = ?", bidID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrBidNotFound
		}
		return nil, err
	}

	newStatus := currBid.Status
	newShips := currBid.ShipsCount
	newVersion := currBid.Version + 1

	backupBid, err := getSingleRecord(ctx, r.db, &models.BidVersion{},
		WithWhere("bid_id = ?", bidID),
		WithWhere("version = ?", version),
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrBidVersionNotFound
		}
		return nil, err
	}

	rollbackBid := trnsfrm.BidVersionToBid(backupBid)
	rollbackBid.Status = newStatus
	rollbackBid.ShipsCount = newShips
	rollbackBid.Version = newVersion

	if err := r.db.WithContext(ctx).
		Model(&models.Bid{}).
		Where("id = ?", bidID).
		Updates(rollbackBid).
		Error; err != nil {
		return nil, err
	}

	if err := r.createBackup(ctx, rollbackBid); err != nil {
		return nil, fmt.Errorf("create bid backup: %w", err)
	}

	return utils.MustTransformObj[models.Bid, entity.Bid](rollbackBid), nil
}

func (r *BidRepo) GetFeedbacksByFilter(ctx context.Context, filters ...FilterOption) ([]entity.BidRewiew, error) {
	return getMultiMappedRecord[entity.BidRewiew, models.BidRewiew](ctx, r.db, filters...)
}

func (r *BidRepo) ShipBid(ctx context.Context, userID uuid.UUID, bidID uuid.UUID) (bool, error) {
	_, err := getSingleRecord(ctx, r.db, &models.BidShip{},
		WithWhere("user_id = ?", userID),
		WithWhere("bid_id = ?", bidID),
	)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		if err := createRecord(ctx, r.db, &models.BidShip{}, &models.BidShip{UserID: userID, BidID: bidID}); err != nil {
			return false, err
		}

		bid, err := r.GetBidByID(ctx, bidID)
		if err != nil {
			return false, fmt.Errorf("get bid: %w", err)
		}

		bid.ShipsCount += 1
		r.db.WithContext(ctx).Model(&models.Bid{}).Where("id = ?", bid.Id).Update("ships_count", bid.ShipsCount)

		return true, nil
	}
	return false, nil
}

func (r *BidRepo) UnshipsBid(ctx context.Context, bidID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&models.BidShip{}).
		Where("bid_id = ?", bidID).
		Delete(&models.BidShip{}).
		Error; err != nil {
		return err
	}

	err := r.db.WithContext(ctx).Model(&models.Bid{}).Where("id = ?", bidID).Update("ships_count", 0).Error

	return err
}
