package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationType string

const (
	IE  OrganizationType = "IE"
	LLC OrganizationType = "LLC"
	JSC OrganizationType = "JSC"
)

var OrganizationTypeList = []OrganizationType{IE, LLC, JSC}

type Organization struct {
	Id          uuid.UUID
	Name        string
	Description string
	Type        OrganizationType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Pagination struct {
	Limit  int
	Offset int
}

type User struct {
	Id        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
