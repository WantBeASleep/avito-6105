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

type TenderUsecase struct {
	tenderRepo repos.TenderRepo
}

func NewTenderUsecase(tenderRepo repos.TenderRepo) *TenderUsecase {
	return &TenderUsecase{
		tenderRepo: tenderRepo,
	}
}

func (u *TenderUsecase) CreateTender(ctx context.Context, username string, tender *entity.Tender) (*entity.Tender, error) {
	ok, err := u.checkUserResponsibleOrg(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("check user permission: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionCreateTender
	}

	tender.Version = 1
	tender.Status = entity.Created

	tender, err = u.tenderRepo.CreateTender(ctx, tender)
	if err != nil {
		return nil, fmt.Errorf("create tender: %w", err)
	}

	return tender, nil
}

func (u *TenderUsecase) GetTenders(ctx context.Context, serviceTypes []entity.TenderServiceType, pag *entity.Pagination) ([]entity.Tender, error) {
	conds := []db.FilterOption{}
	if serviceTypes != nil {
		conds = append(conds, db.WithWhere("service_type IN ?", serviceTypes))
	}
	conds = append(conds,
		db.WithWhere("status = ?", entity.Published),
		db.WithOrder("name asc"),
		db.WithPagination(*pag),
	)

	tenders, err := u.tenderRepo.GetTendersByFilter(ctx, conds...)
	if err != nil {
		return nil, fmt.Errorf("get tenders: %w", err)
	}

	return tenders, nil
}

func (u *TenderUsecase) GetMyTenders(ctx context.Context, username string, pag *entity.Pagination) ([]entity.Tender, error) {
	_, userOrgsIDs, err := u.getUserAndUserOrgsIDs(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user orgs ids: %w", err)
	}

	tenders, err := u.tenderRepo.GetTendersByFilter(ctx,
		db.WithWhere("organization_id IN ?", userOrgsIDs),
		db.WithPagination(*pag),
		db.WithOrder("name asc"),
	)
	if err != nil {
		return nil, fmt.Errorf("get tenders: %w", err)
	}

	return tenders, nil
}

func (u *TenderUsecase) GetTenderStatus(ctx context.Context, username string, tenderID uuid.UUID) (entity.TenderStatusType, error) {
	tender, err := u.tenderRepo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return "", fmt.Errorf("get tender by id: %w", err)
	}

	if tender.Status == entity.Published {
		return tender.Status, nil
	}

	if username == "" {
		return "", entity.ErrUserNotSpecified
	}

	ok, err := u.checkPermissionForTender(ctx, username, tender.Id)
	if err != nil {
		return "", fmt.Errorf("check user permission: %w", err)
	}
	if !ok {
		return "", entity.ErrUserPermissionTender
	}

	return tender.Status, nil
}

func (u *TenderUsecase) UpdateTenderStatus(ctx context.Context, username string, tenderID uuid.UUID, status entity.TenderStatusType) (*entity.Tender, error) {
	ok, err := u.checkPermissionForTender(ctx, username, tenderID)
	if err != nil {
		return nil, fmt.Errorf("check user permission: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionTender
	}

	if err := u.tenderRepo.UpdateTenderStatus(ctx, tenderID, status); err != nil {
		return nil, fmt.Errorf("update tender status %w", err)
	}

	tender, err := u.tenderRepo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return nil, fmt.Errorf("get tender by id: %w", err)
	}

	return tender, nil
}

func (u *TenderUsecase) PatchTender(ctx context.Context, username string, tenderID uuid.UUID, patchTender *entity.Tender) (*entity.Tender, error) {
	ok, err := u.checkPermissionForTender(ctx, username, tenderID)
	if err != nil {
		return nil, fmt.Errorf("check user permission: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionTender
	}

	tender, err := u.tenderRepo.PatchTender(ctx, tenderID, patchTender)
	if err != nil {
		return nil, fmt.Errorf("patch tender: %w", err)
	}

	return tender, nil
}

func (u *TenderUsecase) RollbackTender(ctx context.Context, username string, tenderID uuid.UUID, version int) (*entity.Tender, error) {
	ok, err := u.checkPermissionForTender(ctx, username, tenderID)
	if err != nil {
		return nil, fmt.Errorf("check user permission: %w", err)
	}
	if !ok {
		return nil, entity.ErrUserPermissionTender
	}

	tender, err := u.tenderRepo.RollbackTender(ctx, tenderID, version)
	if err != nil {
		return nil, fmt.Errorf("rollback tender: %w", err)
	}

	return tender, nil
}

func (u *TenderUsecase) getUserAndUserOrgsIDs(ctx context.Context, username string) (*entity.User, uuid.UUIDs, error) {
	user, err := u.tenderRepo.GetUserByUserName(ctx, username)
	if err != nil {
		return nil, nil, fmt.Errorf("get user by user name: %w", err)
	}

	userOrgsUUIDs, err := u.tenderRepo.GetUserOrgsUUIDs(ctx, user.Id)
	if err != nil {
		return nil, nil, fmt.Errorf("get user organizations: %w", err)
	}

	return user, userOrgsUUIDs, nil
}

func (u *TenderUsecase) checkUserResponsibleOrg(ctx context.Context, username string, orgID uuid.UUID) (bool, error) {
	_, userResponsibleOrgsIDs, err := u.getUserAndUserOrgsIDs(ctx, username)
	if err != nil {
		return false, fmt.Errorf("get user and user orgs id: %w", err)
	}

	return slices.Contains(userResponsibleOrgsIDs, orgID), nil
}

func (u *TenderUsecase) checkPermissionForTender(ctx context.Context, username string, tenderID uuid.UUID) (bool, error) {
	tender, err := u.tenderRepo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return false, fmt.Errorf("get tender by id: %w", err)
	}

	ok, err := u.checkUserResponsibleOrg(ctx, username, tender.OrganizationID)
	if err != nil {
		return false, fmt.Errorf("check user permission: %w", err)
	}
	return ok, nil
}
