package usecases

import (
	db "avito/internal/db/repos"
	"avito/internal/entity"
	"avito/internal/usecases/repos"
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type BidUsecase struct {
	tenderRepo    repos.TenderRepo
	bidRepo       repos.BidRepo
	tenderUsecase *TenderUsecase
}

func NewBidUsecase(
	tenderRepo repos.TenderRepo,
	bidRepo repos.BidRepo,
	tenderUsecase *TenderUsecase,
) *BidUsecase {
	return &BidUsecase{
		tenderRepo:    tenderRepo,
		bidRepo:       bidRepo,
		tenderUsecase: tenderUsecase,
	}
}

func (u *BidUsecase) CreateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error) {
	tender, err := u.tenderRepo.GetTenderByID(ctx, bid.TenderID)
	if err != nil {
		return nil, fmt.Errorf("get tender by id: %w", err)
	}

	if tender.Status != entity.Published {
		return nil, entity.ErrCreateBidTender
	}

	users, err := u.tenderRepo.GetOrgUsersIDsByID(ctx, tender.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("get org users: %w", err)
	}
	bid.Kvorum = min(3, len(users))

	if bid.AuthorType == entity.AuthorUser {
		_, err := u.tenderRepo.GetUserByID(ctx, bid.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("get user by id: %w", err)
		}

	} else {
		_, err := u.tenderRepo.GetOrgByID(ctx, bid.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("get org by id: %w", err)
		}
	}

	bid.Version = 1
	bid.Status = entity.BCreated

	bid, err = u.bidRepo.CreateBid(ctx, bid)
	if err != nil {
		return nil, fmt.Errorf("bid create: %w", err)
	}

	return bid, nil
}

func (u *BidUsecase) GetMyBids(ctx context.Context, username string, pag *entity.Pagination) ([]entity.Bid, error) {
	user, orgsIDs, err := u.tenderUsecase.getUserAndUserOrgsIDs(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user orgs ids: %w", err)
	}

	bids, err := u.bidRepo.GetBidsByFilter(ctx,
		db.WithWhere("author_id = ?", user.Id),
		db.WithOr("author_id IN ?", orgsIDs),
		db.WithOrder("name asc"),
		db.WithPagination(*pag),
	)
	if err != nil {
		return nil, fmt.Errorf("get bids: %w", err)
	}

	return bids, nil
}

func (u *BidUsecase) GetTenderBidsList(ctx context.Context, username string, tenderID uuid.UUID, pag *entity.Pagination) ([]entity.Bid, error) {
	_, err := u.tenderRepo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return nil, fmt.Errorf("get tender by id: %w", err)
	}

	user, orgsIDs, err := u.tenderUsecase.getUserAndUserOrgsIDs(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user orgs ids: %w", err)
	}

	opts := []db.FilterOption{}

	opts = append(opts, db.WithOr("author_id = ?", user.Id))
	opts = append(opts, db.WithOr("author_id IN ?", orgsIDs))

	ok, err := u.tenderUsecase.checkPermissionForTender(ctx, username, tenderID)
	if err != nil {
		return nil, fmt.Errorf("check permissions: %w", err)
	}
	if ok {
		opts = append(opts, db.WithOr("status = ?", entity.BPublished))
	}

	bids, err := u.bidRepo.GetBidsByFilter(ctx,
		db.WithWhere("tender_id = ?", tenderID),
		db.WithOrGroupFilters(opts, u.bidRepo),
		db.WithOrder("name asc"),
		db.WithPagination(*pag),
	)
	if err != nil {
		return nil, fmt.Errorf("get tender bids: %w", err)
	}
	if len(bids) == 0 {
		return nil, entity.ErrUserPermissionBidsTender
	}

	return bids, nil
}

func (u *BidUsecase) GetBidStatus(ctx context.Context, username string, bidID uuid.UUID) (entity.BidStatusType, error) {
	bid, err := u.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		return "", fmt.Errorf("get bid by id: %w", err)
	}

	ok, err := u.checkUserOwnerBid(ctx, username, bidID)
	if err != nil {
		return "", fmt.Errorf("check user bid permission: %w", err)
	}
	if ok {
		return bid.Status, nil
	}

	ok, err = u.checkUserOwnerTenderByBid(ctx, username, bidID)
	if err != nil {
		return "", fmt.Errorf("check user bid permission: %w", err)
	}
	if ok && bid.Status == entity.BPublished {
		return entity.BPublished, nil
	}

	return "", entity.ErrUserPermissionBid
}

func (u *BidUsecase) UpdateBidStatus(ctx context.Context, username string, bidID uuid.UUID, newStatus entity.BidStatusType) (*entity.Bid, error) {
	ok, err := u.checkUserOwnerBid(ctx, username, bidID)
	if err != nil {
		return nil, fmt.Errorf("check user bid permission: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionBid
	}

	if err := u.bidRepo.UpdateBidStatus(ctx, bidID, newStatus); err != nil {
		return nil, fmt.Errorf("update bid status by id: %w", err)
	}

	bid, err := u.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		return nil, fmt.Errorf("get bid by id: %w", err)
	}

	return bid, nil
}

func (u *BidUsecase) PatchBid(ctx context.Context, username string, bidID uuid.UUID, bid *entity.Bid) (*entity.Bid, error) {
	ok, err := u.checkUserOwnerBid(ctx, username, bidID)
	if err != nil {
		return nil, fmt.Errorf("check user bid permission: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionBid
	}

	bid, err = u.bidRepo.PatchBid(ctx, bidID, bid)
	if err != nil {
		return nil, fmt.Errorf("patch bid: %w", err)
	}

	return bid, nil
}

func (u *BidUsecase) SubmitDecision(ctx context.Context, username string, bidID uuid.UUID, decision entity.BidDecisionType) (*entity.Bid, error) {
	ok, err := u.checkUserOwnerTenderByBid(ctx, username, bidID)
	if err != nil {
		return nil, fmt.Errorf("check user owner tender by bid: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionShipBid
	}

	bid, err := u.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		return nil, fmt.Errorf("get bid by id: %w", err)
	}

	if bid.Status != entity.BPublished {
		return nil, entity.ErrShipBidTender
	}

	if decision == entity.Rejected {
		bid.Status = entity.BCanceled
		bid.ShipsCount = 0

		if err := u.bidRepo.UnshipsBid(ctx, bidID); err != nil {
			return nil, fmt.Errorf("unship bid: %w", err)
		}

		err := u.bidRepo.UpdateBidStatus(ctx, bidID, entity.BCanceled)
		if err != nil {
			return nil, fmt.Errorf("update bid to approved: %w", err)
		}

		return bid, nil
	} else {
		user, _, err := u.tenderUsecase.getUserAndUserOrgsIDs(ctx, username)
		if err != nil {
			return nil, fmt.Errorf("get user")
		}

		shipped, err := u.bidRepo.ShipBid(ctx, user.Id, bidID)
		if err != nil {
			return nil, fmt.Errorf("ship bid: %w", err)
		}

		if shipped {
			bid.ShipsCount += 1
		}

		if bid.ShipsCount >= bid.Kvorum {
			// err := u.bidRepo.UpdateBidStatus(ctx, bidID, entity.BApproved)
			// if err != nil {
			// 	return nil, fmt.Errorf("update bid to approved: %w", err)
			// }
			// bid.Status = entity.BApproved

			if err := u.tenderRepo.UpdateTenderStatus(ctx, bid.TenderID, entity.Closed); err != nil {
				return nil, fmt.Errorf("update tender status by id: %w", err)
			}
		}

		return bid, nil
	}
}

func (u *BidUsecase) FeedbackBid(ctx context.Context, username string, bidID uuid.UUID, bidFeedback string) (*entity.Bid, error) {
	ok, err := u.checkUserOwnerTenderByBid(ctx, username, bidID)
	if err != nil {
		return nil, fmt.Errorf("check user owner tender by bid: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionBid
	}

	bid, err := u.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		return nil, fmt.Errorf("get bid by id: %w", err)
	}

	if bid.ShipsCount < bid.Kvorum {
		return nil, entity.ErrUserPermissionRewiew
	}
	// if bid.Status != entity.BApproved {
	// 	return nil, entity.ErrUserPermissionRewiew
	// }

	_, err = u.bidRepo.CreateFeedback(ctx, &entity.BidRewiew{
		Description: bidFeedback,
		BidID:       bidID,
	})

	if err != nil {
		return nil, fmt.Errorf("create feedback: %w", err)
	}

	return bid, nil
}

func (u *BidUsecase) RollbackBid(ctx context.Context, username string, bidID uuid.UUID, version int) (*entity.Bid, error) {
	ok, err := u.checkUserOwnerBid(ctx, username, bidID)
	if err != nil {
		return nil, fmt.Errorf("check user owner bid: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionBid
	}

	bid, err := u.bidRepo.RollbackBid(ctx, bidID, version)
	if err != nil {
		return nil, fmt.Errorf("rollback bid: %w", err)
	}

	return bid, nil
}

func (u *BidUsecase) CheckPrevFeedbacks(ctx context.Context, tenderID uuid.UUID, author string, requester string, pag entity.Pagination) ([]entity.BidRewiew, error) {
	tender, err := u.tenderRepo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return nil, fmt.Errorf("get tender by id: %w", err)
	}

	ok, err := u.tenderUsecase.checkPermissionForTender(ctx, requester, tender.Id)
	if err != nil {
		return nil, fmt.Errorf("check user permission: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionTender
	}

	authorEnt, authorOrgsIDs, err := u.tenderUsecase.getUserAndUserOrgsIDs(ctx, author)
	if err != nil {
		return nil, fmt.Errorf("get user orgs ids: %w", err)
	}

	authorBids, err := u.bidRepo.GetBidsByFilter(ctx,
		db.WithOr("author_id = ?", authorEnt.Id),
		db.WithOr("author_id IN ?", authorOrgsIDs),
	)
	if err != nil {
		return nil, fmt.Errorf("get bids: %w", err)
	}

	finded := false
	for _, bids := range authorBids {
		if bids.TenderID == tenderID {
			finded = true
		}
	}
	if !finded {
		return nil, entity.ErrFeedbackPermission
	}

	bidsIds := uuid.UUIDs{}
	for _, b := range authorBids {
		bidsIds = append(bidsIds, b.Id)
	}

	feedbacks, err := u.bidRepo.GetFeedbacksByFilter(ctx,
		db.WithWhere("bid_id IN ?", bidsIds),
		db.WithOrder("description asc"),
		db.WithPagination(pag),
	)
	if err != nil {
		return nil, fmt.Errorf("get feedbacks: %w", err)
	}

	return feedbacks, nil
}

func (u *BidUsecase) checkUserOwnerTenderByBid(ctx context.Context, username string, bidID uuid.UUID) (bool, error) {
	bid, err := u.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		return false, fmt.Errorf("get bid by id: %w", err)
	}

	tender, err := u.tenderRepo.GetTenderByID(ctx, bid.TenderID)
	if err != nil {
		return false, fmt.Errorf("get tender by id: %w", err)
	}

	_, orgsIDs, err := u.tenderUsecase.getUserAndUserOrgsIDs(ctx, username)
	if err != nil {
		return false, fmt.Errorf("get user orgs ids: %w", err)
	}

	if slices.Contains(orgsIDs, tender.OrganizationID) {
		return true, nil
	}

	return false, nil
}

func (u *BidUsecase) checkUserOwnerBid(ctx context.Context, username string, bidID uuid.UUID) (bool, error) {
	bid, err := u.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		return false, fmt.Errorf("get bid by id: %w", err)
	}

	user, orgsIDs, err := u.tenderUsecase.getUserAndUserOrgsIDs(ctx, username)
	if err != nil {
		return false, fmt.Errorf("get user orgs ids: %w", err)
	}

	if user.Id == bid.AuthorID || slices.Contains(orgsIDs, bid.AuthorID) {
		return true, nil
	}

	return false, nil
}
