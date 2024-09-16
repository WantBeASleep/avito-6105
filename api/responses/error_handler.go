package responses

import (
	"avito/api/validation"
	"avito/internal/entity"
	"errors"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, err error) {
	var validateErr *validation.ValidateError

	switch {
	case errors.Is(err, entity.ErrUserNotFound):
		ErrorJSON(w, http.StatusUnauthorized, entity.ErrUserNotFound)

	case errors.Is(err, entity.ErrOrgNotFound):
		ErrorJSON(w, http.StatusNotFound, entity.ErrOrgNotFound)

	case errors.Is(err, entity.ErrTenderNotFound):
		ErrorJSON(w, http.StatusNotFound, entity.ErrTenderNotFound)

	case errors.Is(err, entity.ErrBidNotFound):
		ErrorJSON(w, http.StatusNotFound, entity.ErrBidNotFound)

	case errors.Is(err, entity.ErrTenderVersionNotFound):
		ErrorJSON(w, http.StatusNotFound, entity.ErrTenderVersionNotFound)

	case errors.Is(err, entity.ErrBidVersionNotFound):
		ErrorJSON(w, http.StatusNotFound, entity.ErrBidVersionNotFound)

	case errors.Is(err, entity.ErrUserNotSpecified):
		ErrorJSON(w, http.StatusUnauthorized, entity.ErrUserNotSpecified)

	case errors.Is(err, entity.ErrUserPermissionTender):
		ErrorJSON(w, http.StatusForbidden, entity.ErrUserPermissionTender)

	case errors.Is(err, entity.ErrCreateBidTender):
		ErrorJSON(w, http.StatusForbidden, entity.ErrCreateBidTender)

	case errors.Is(err, entity.ErrUserPermissionBid):
		ErrorJSON(w, http.StatusForbidden, entity.ErrUserPermissionBid)

	case errors.Is(err, entity.ErrUserPermissionBidsTender):
		ErrorJSON(w, http.StatusForbidden, entity.ErrUserPermissionBidsTender)

	case errors.Is(err, entity.ErrUserPermissionCreateTender):
		ErrorJSON(w, http.StatusForbidden, entity.ErrUserPermissionCreateTender)

	case errors.Is(err, entity.ErrUserPermissionShipBid):
		ErrorJSON(w, http.StatusForbidden, entity.ErrUserPermissionShipBid)

	case errors.Is(err, entity.ErrFeedbackPermission):
		ErrorJSON(w, http.StatusForbidden, entity.ErrFeedbackPermission)

	case errors.Is(err, entity.ErrUserPermissionRewiew):
		ErrorJSON(w, http.StatusForbidden, entity.ErrUserPermissionRewiew)

	case errors.Is(err, entity.ErrShipBidTender):
		ErrorJSON(w, http.StatusBadRequest, entity.ErrShipBidTender)

	case errors.Is(err, validation.ErrParsed):
		ErrorJSON(w, http.StatusBadRequest, validation.ErrParsed)

	case errors.As(err, &validateErr):
		ErrorJSON(w, http.StatusBadRequest, validateErr)

	default:
		Error(w, http.StatusInternalServerError)
	}
}
