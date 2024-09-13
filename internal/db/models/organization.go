package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrganizationType string

const (
	IE  OrganizationType = "IE"
	LLC OrganizationType = "LLC"
	JSC OrganizationType = "JSC"
)

var orgTypeList = []OrganizationType{IE, LLC, JSC}

func (o *OrganizationType) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("not string org_type value")
	}

	*o = OrganizationType(strValue)
	return nil
}

func (o OrganizationType) Value() (driver.Value, error) {
	for _, validType := range orgTypeList {
		if o == validType {
			return string(o), nil
		}
	}
	return nil, fmt.Errorf("invalid org_type value: %s", o)
}

const OrganizationName = "organization"

type Organization struct {
	Id          uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`
	Name        string           `gorm:"type:varchar(100);not null"`
	Description string           `gorm:"type:text"`
	Type        OrganizationType `gorm:"type:organization_type"`
	CreatedAt   time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (Organization) TableName() string {
	return OrganizationName
}

const OrganizationResponsibleName = "organization_responsible"

type OrganizationResponsible struct {
	Id uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`

	OrganizationID uuid.UUID    `gorm:"type:uuid;not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;references:Id;constraint:OnDelete:CASCADE;"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserID;references:Id;constraint:OnDelete:CASCADE;"`
}

func (OrganizationResponsible) TableName() string {
	return OrganizationResponsibleName
}
