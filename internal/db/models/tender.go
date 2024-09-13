package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TenderServiceType string

const (
	Construction TenderServiceType = "Construction"
	Delivery     TenderServiceType = "Delivery"
	Manufacture  TenderServiceType = "Manufacture"
)

var TenderServiceTypeList = []TenderServiceType{Construction, Delivery, Manufacture}

func (s *TenderServiceType) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("not string tender_service_type value")
	}

	*s = TenderServiceType(strValue)
	return nil
}

func (s TenderServiceType) Value() (driver.Value, error) {
	for _, validType := range TenderServiceTypeList {
		if s == validType {
			return string(s), nil
		}
	}
	return nil, fmt.Errorf("invalid tender_service_type value: %s", s)
}

type TenderStatusType string

const (
	Created   TenderStatusType = "Created"
	Published TenderStatusType = "Published"
	Closed    TenderStatusType = "Closed"
)

var TenderStatusTypeList = []TenderStatusType{Created, Published, Closed}

func (s *TenderStatusType) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("not string tender_status_type value")
	}

	*s = TenderStatusType(strValue)
	return nil
}

func (s TenderStatusType) Value() (driver.Value, error) {
	for _, validType := range TenderStatusTypeList {
		if s == validType {
			return string(s), nil
		}
	}
	return nil, fmt.Errorf("invalid tender_status_type value: %s", s)
}

const TenderName = "tender"

type Tender struct {
	Id          uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`
	Name        string            `gorm:"type:varchar(100);not null"`
	Description string            `gorm:"type:varchar(500);"`
	ServiceType TenderServiceType `gorm:"type:service_type;"`
	Status      TenderStatusType  `gorm:"type:tender_status_type;"`

	OrganizationID uuid.UUID    `gorm:"type:uuid;not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;references:Id;" copier:"-"`

	Version   int       `gorm:"type:bigint"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (Tender) TableName() string {
	return TenderName
}

const TenderVersionName = "tender_backup"

type TenderVersion struct {
	Id uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;" copier:"-"`

	TenderID uuid.UUID `gorm:"type:uuid;not null"`
	Tender   Tender    `gorm:"foreignKey:TenderID;references:Id;" copier:"-"`

	Name        string            `gorm:"type:varchar(100);not null"`
	Description string            `gorm:"type:varchar(500);"`
	ServiceType TenderServiceType `gorm:"type:service_type;"`
	Status      TenderStatusType  `gorm:"type:tender_status_type;"`

	OrganizationID uuid.UUID    `gorm:"type:uuid;not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;references:Id;" copier:"-"`

	Version   int       `gorm:"type:bigint"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (TenderVersion) TableName() string {
	return TenderVersionName
}
