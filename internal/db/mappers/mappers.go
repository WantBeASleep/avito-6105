package mappers

import (
	"avito/internal/db/models"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type Transform struct{}

func (Transform) TenderToTenderVersion(tender *models.Tender) *models.TenderVersion {
	backup := models.TenderVersion{}
	copier.Copy(&backup, tender)
	backup.TenderID = tender.Id

	return &backup
}

func (Transform) TenderVersionToTender(backup *models.TenderVersion) *models.Tender {
	tender := models.Tender{}
	copier.Copy(&tender, backup)
	tender.Id = backup.TenderID

	return &tender
}

func (Transform) BidToBidVersion(bid *models.Bid) *models.BidVersion {
	backup := models.BidVersion{}
	copier.Copy(&backup, bid)
	backup.BidID = bid.Id

	return &backup
}

func (Transform) BidVersionToBid(backup *models.BidVersion) *models.Bid {
	bid := models.Bid{}
	copier.Copy(&bid, backup)
	bid.Id = backup.BidID

	return &bid
}

func (Transform) OrgRespToOrgUUIDSlice(mdl []models.OrganizationResponsible) uuid.UUIDs {
	slice := uuid.UUIDs{}
	for _, m := range mdl {
		slice = append(slice, m.OrganizationID)
	}

	return slice
}

func (Transform) OrgRespToUserUUIDSlice(mdl []models.OrganizationResponsible) uuid.UUIDs {
	slice := uuid.UUIDs{}
	for _, m := range mdl {
		slice = append(slice, m.UserID)
	}

	return slice
}
