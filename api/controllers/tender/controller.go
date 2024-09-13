package tender

import (
	"avito/api/parsers"
	"avito/api/responses"
	"avito/api/usecases"
	"avito/api/validation"
	"avito/internal/entity"
	"avito/internal/utils"
	"encoding/json"
	"net/http"
)

type Controller struct {
	tenderUsecase usecases.TenderUsecase
}

func NewTenderController(tenderUsecase usecases.TenderUsecase) *Controller {
	return &Controller{
		tenderUsecase: tenderUsecase,
	}
}

func (c *Controller) CreateTender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createTender CreateTender
	if err := json.NewDecoder(r.Body).Decode(&createTender); err != nil {
		responses.ErrorHandler(w, validation.ErrParsed)
		return
	}

	if err := validation.ValidateStruct(&createTender); err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	tender := utils.MustTransformObj[CreateTender, entity.Tender](&createTender)

	resp, err := c.tenderUsecase.CreateTender(ctx, createTender.CreatorUserName, tender)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) GetTenders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pagination, err := parsers.ParsePagination(r)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	serviceTypes, err := parsers.ParseServiceTypes(r)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.tenderUsecase.GetTenders(ctx, serviceTypes, pagination)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) GetMyTenders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pagination, err := parsers.ParsePagination(r)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	// НЕ ПОНЯТНО, КАК ЛУЧШЕ ОСТАВИТЬ???
	username, _ := parsers.ParseQuery(r, "username", false, parsers.ParserEmptyString)
	// username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	// if err != nil {
	// 	responses.ErrorHandler(w, err)
	// 	return
	// }

	resp, err := c.tenderUsecase.GetMyTenders(ctx, username, pagination)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenderID, err := parsers.ParseVar(r, "tenderId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, _ := parsers.ParseQuery(r, "username", false, parsers.ParserEmptyString)

	resp, err := c.tenderUsecase.GetTenderStatus(ctx, username, tenderID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenderID, err := parsers.ParseVar(r, "tenderId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	status, err := parsers.ParseQuery(r, "status", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}
	if err := validation.ValidateOneOf(entity.TenderStatusTypeList, status, "status"); err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.tenderUsecase.UpdateTenderStatus(ctx, username, tenderID, entity.TenderStatusType(status))
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) PatchTender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenderID, err := parsers.ParseVar(r, "tenderId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	var patchTender PatchTender
	if err := json.NewDecoder(r.Body).Decode(&patchTender); err != nil {
		responses.ErrorHandler(w, validation.ErrParsed)
		return
	}

	if err := validation.ValidateStruct(patchTender); err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	patchTenderEnt := utils.MustTransformObj[PatchTender, entity.Tender](&patchTender)

	resp, err := c.tenderUsecase.PatchTender(ctx, username, tenderID, patchTenderEnt)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) RollbackTender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenderID, err := parsers.ParseVar(r, "tenderId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	version, err := parsers.ParseVar(r, "version", true, parsers.ParserInt)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.tenderUsecase.RollbackTender(ctx, username, tenderID, version)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}
