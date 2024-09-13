package entity

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrOrgNotFound    = errors.New("organization not found")
	ErrTenderNotFound = errors.New("tender not found")
	ErrBidNotFound    = errors.New("bid not found")
)

var (
	ErrUserNotSpecified = errors.New("only authorizated users have permissions to view this resource")
)

var (
	ErrTenderVersionNotFound = errors.New("tender backup version not found")
	ErrBidVersionNotFound    = errors.New("bid backup version not found")
)

var (
	ErrUserPermissionCreateTender = errors.New("user dont have permission to create tender for this org")
	ErrUserPermissionTender       = errors.New("user dont have permission to this tender")
	ErrUserPermissionBidsTender       = errors.New("user dont have permission to see bids for this tender")
	ErrCreateBidTender            = errors.New("cant create bid to not public tender")
	ErrShipBidTender              = errors.New("cant ship not public bid")
	ErrFeedbackPermission         = errors.New("cant see this feedbacks")
	ErrUserPermissionBid          = errors.New("user dont have permission to this bid")
	ErrUserPermissionShipBid      = errors.New("user dont have permission to ship this bid")
	ErrUserPermissionRewiew       = errors.New("cant create rewiew to not approved bid")
)
