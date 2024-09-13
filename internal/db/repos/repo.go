package repos

import (
	"avito/internal/db/mappers"
	"avito/internal/db/models"
	"avito/internal/entity"
	"avito/internal/utils"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Clear interface {
	GetClear() *gorm.DB
}

// struct mapping
var trnsfrm = mappers.Transform{}

func newLogger() logger.Interface {
	return logger.New(&gormLog{},
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
}

type gormLog struct {
}

func (l *gormLog) Printf(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format, args...)
}

// inited
type ctrl struct {
	mu     sync.Mutex
	inited bool
}

var repoCtrl = &ctrl{}

func (c *ctrl) initIfNeed(db *gorm.DB) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.inited {
		return nil
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	db.Exec("CREATE TYPE service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');")
	db.Exec("CREATE TYPE tender_status_type AS ENUM ('Created', 'Published', 'Closed');")

	db.Exec("CREATE TYPE author_type AS ENUM ('Organization', 'User');")
	db.Exec("CREATE TYPE bid_status_type AS ENUM ('Created', 'Published', 'Canceled', 'Approved');")

	err := db.AutoMigrate(
		&models.Tender{},
		&models.TenderVersion{},

		&models.Bid{},
		&models.BidVersion{},
		&models.BidRewiew{},
		&models.BidShip{},
	)

	c.inited = true

	return err
}

type FilterOption func(db *gorm.DB) *gorm.DB

var (
	WithWhere = func(condition string, args ...any) FilterOption {
		return FilterOption(func(db *gorm.DB) *gorm.DB {
			return db.Where(condition, args...)
		})
	}

	WithOr = func(condition string, args ...any) FilterOption {
		return FilterOption(func(db *gorm.DB) *gorm.DB {
			return db.Or(condition, args...)
		})
	}

	WithOrder = func(condition string) FilterOption {
		return FilterOption(func(db *gorm.DB) *gorm.DB {
			db = db.Order(condition)
			return db
		})
	}

	WithOrGroupFilters = func(conds []FilterOption, cdb Clear) FilterOption {
		return FilterOption(func(db *gorm.DB) *gorm.DB {
			internal := cdb.GetClear()
			for i, cnd := range conds {
				if i == 0 {
					internal = internal.Where(cnd(cdb.GetClear()))
				} else {
					internal = internal.Or(cnd(cdb.GetClear()))
				}
			}
			return db.Where(internal)
		})
	}

	WithPagination = func(pag entity.Pagination) FilterOption {
		return FilterOption(func(db *gorm.DB) *gorm.DB {
			db = db.Limit(pag.Limit)
			db = db.Offset(pag.Offset)
			return db
		})
	}
)

func createRecord[T any](ctx context.Context, db *gorm.DB, model *T, value *T, opts ...FilterOption) error {
	query := db.WithContext(ctx).Model(model)
	for _, opt := range opts {
		query = opt(query)
	}

	return query.Create(value).Error
}

func getSingleRecord[T any](ctx context.Context, db *gorm.DB, model *T, opts ...FilterOption) (*T, error) {
	var resp T

	query := db.WithContext(ctx).Model(model)
	for _, opt := range opts {
		query = opt(query)
	}

	if err := query.Take(&resp).Error; err != nil {
		return nil, err
	}

	return &resp, nil
}

func getMultiRecord[T any](ctx context.Context, db *gorm.DB, model *T, opts ...FilterOption) ([]T, error) {
	var resp []T

	query := db.WithContext(ctx).Model(model)
	for _, opt := range opts {
		query = opt(query)
	}

	if err := query.Find(&resp).Error; err != nil {
		return nil, err
	}

	return resp, nil
}

func getSingleMappedRecord[E, M any](ctx context.Context, db *gorm.DB, errNotFound error, opts ...FilterOption) (*E, error) {
	var model M
	resp, err := getSingleRecord(ctx, db, &model, opts...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errNotFound
		}
	}

	return utils.MustTransformObj[M, E](resp), nil
}

func getMultiMappedRecord[E, M any](ctx context.Context, db *gorm.DB, opts ...FilterOption) ([]E, error) {
	var model M
	resp, err := getMultiRecord(ctx, db, &model, opts...)
	if err != nil {
		return nil, err
	}

	return utils.MustTransformSlice[M, E](resp), nil
}
