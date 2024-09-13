package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BidAuthorType string

const (
	AuthorOrganization BidAuthorType = "Organization"
	AuthorUser         BidAuthorType = "User"
)

var BidAuthorTypeList = []BidAuthorType{AuthorOrganization, AuthorUser}

func (o *BidAuthorType) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("not string author_type value")
	}

	*o = BidAuthorType(strValue)
	return nil
}

func (o BidAuthorType) Value() (driver.Value, error) {
	for _, validType := range BidAuthorTypeList {
		if o == validType {
			return string(o), nil
		}
	}
	return nil, fmt.Errorf("invalid author_type value: %s", o)
}

type BidStatusType string

const (
	BCreated   BidStatusType = "Created"
	BPublished BidStatusType = "Published"
	BCanceled  BidStatusType = "Canceled"
)

var BidStatusTypeList = []BidStatusType{BCreated, BPublished, BCanceled}

func (s *BidStatusType) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("not string bid_status_type value")
	}

	*s = BidStatusType(strValue)
	return nil
}

func (s BidStatusType) Value() (driver.Value, error) {
	for _, validType := range BidStatusTypeList {
		if s == validType {
			return string(s), nil
		}
	}
	return nil, fmt.Errorf("invalid bid_status_type value: %s", s)
}

const BidName = "bid"

type Bid struct {
	Id          uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`
	Name        string        `gorm:"type:varchar(100);not null"`
	Description string        `gorm:"type:text;not null"`
	Status      BidStatusType `gorm:"type:bid_status_type;not null"`
	TenderID    uuid.UUID     `gorm:"type:uuid;not null"`
	Tender      Tender        `gorm:"foreignKey:TenderID;references:Id;" copier:"-"`
	AuthorType  BidAuthorType `gorm:"type:author_type;not null"`
	AuthorID    uuid.UUID     `gorm:"type:uuid;not null"`
	Version     int           `gorm:"type:bigint;not null"`
	CreatedAt   time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	ShipsCount int `gorm:"type:bigint;default:0;not null"`
	Kvorum     int `gorm:"type:bigint;not null"`
}

func (Bid) TableName() string {
	return BidName
}

const BidVersionName = "bid_backup"

type BidVersion struct {
	Id          uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;" copier:"-"`
	BidID       uuid.UUID     `gorm:"type:uuid;not null"`
	Name        string        `gorm:"type:varchar(100);not null"`
	Description string        `gorm:"type:text;not null"`
	Status      BidStatusType `gorm:"type:bid_status_type;not null"`
	TenderID    uuid.UUID     `gorm:"type:uuid;not null"`
	Tender      Tender        `gorm:"foreignKey:TenderID;references:Id;" copier:"-"`
	AuthorType  BidAuthorType `gorm:"type:author_type;not null"`
	AuthorID    uuid.UUID     `gorm:"type:uuid;not null"`
	Version     int           `gorm:"type:bigint;not null"`
	CreatedAt   time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	ShipsCount int `gorm:"type:bigint;not null"`
	Kvorum     int `gorm:"type:bigint;not null"`
}

func (BidVersion) TableName() string {
	return BidVersionName
}

const BidShipsName = "bid_ship"

type BidShip struct {
	Id uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;" copier:"-"`

	UserID uuid.UUID `gorm:"type:uuid;"`
	User   User      `gorm:"foreignKey:UserID;references:Id;" copier:"-"`

	BidID uuid.UUID `gorm:"type:uuid;"`
	Bid   Bid       `gorm:"foreignKey:BidID;references:Id;" copier:"-"`
}

func (BidShip) TableName() string {
	return BidShipsName
}

const BidRewiewName = "bid_rewiew"

type BidRewiew struct {
	Id          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	BidID uuid.UUID `gorm:"type:uuid"`
	Bid   Bid       `gorm:"foreignKey:BidID;references:Id;" copier:"-"`
}

func (BidRewiew) TableName() string {
	return BidRewiewName
}
