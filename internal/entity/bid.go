package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BidDecisionType string

const (
	Approved BidDecisionType = "Approved"
	Rejected BidDecisionType = "Rejected"
)

var BidDecisionTypeList = []BidDecisionType{Approved, Rejected}

type BidAuthorType string

const (
	AuthorOrganization BidAuthorType = "Organization"
	AuthorUser         BidAuthorType = "User"
)

var BidAuthorTypeList = []BidAuthorType{AuthorOrganization, AuthorUser}

type BidStatusType string

const (
	BCreated   BidStatusType = "Created"
	BPublished BidStatusType = "Published"
	BCanceled  BidStatusType = "Canceled"
	BApproved  BidStatusType = "Approved"
)

var BidStatusTypeList = []BidStatusType{BCreated, BPublished, BCanceled, BApproved}

type Bid struct {
	Id          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      BidStatusType `json:"status"`
	TenderID    uuid.UUID     `json:"tenderId"`
	AuthorType  BidAuthorType `json:"authorType"`
	AuthorID    uuid.UUID     `json:"authorId"`
	Version     int           `json:"version"`
	CreatedAt   time.Time     `json:"createdAt"`
	ShipsCount  int           `json:"-"`
	Kvorum      int           `json:"-"`
}

func (b Bid) MarshalJSON() ([]byte, error) {
	type Alias Bid
	return json.Marshal(
		struct {
			*Alias
			CreatedAt string `json:"createdAt"`
		}{
			Alias:     (*Alias)(&b),
			CreatedAt: b.CreatedAt.Format(time.RFC3339),
		},
	)
}

type BidRewiew struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`

	BidID uuid.UUID `json:"-"`
}

func (b BidRewiew) MarshalJSON() ([]byte, error) {
	type Alias BidRewiew
	return json.Marshal(
		struct {
			*Alias
			CreatedAt string `json:"createdAt"`
		}{
			Alias:     (*Alias)(&b),
			CreatedAt: b.CreatedAt.Format(time.RFC3339),
		},
	)
}
