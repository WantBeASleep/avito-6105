package bid

import (
	"avito/internal/entity"

	"github.com/google/uuid"
)

type CreateBid struct {
	Name        string               `json:"name" validate:"required,max=100"`
	Description string               `json:"description" validate:"required,max=500"`
	TenderID    uuid.UUID            `json:"tenderId" validate:"required,max=100,uuid4"`
	AuthorType  entity.BidAuthorType `json:"authorType" validate:"required,oneof=Organization User"`
	AuthorID    uuid.UUID            `json:"authorId" validate:"required,max=100,uuid4"`
}

type PatchBid struct {
	Name        string `json:"name" validate:"max=100"`
	Description string `json:"description" validate:"max=500"`
}
