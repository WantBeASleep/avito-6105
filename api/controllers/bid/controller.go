package bid

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
	bidUsecase usecases.BidUsecase
}

func NewBidController(bidUsecase usecases.BidUsecase) *Controller {
	return &Controller{
		bidUsecase: bidUsecase,
	}
}

func (c *Controller) CreateBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createBid CreateBid
	if err := json.NewDecoder(r.Body).Decode(&createBid); err != nil {
		responses.ErrorHandler(w, validation.ErrParsed)
		return
	}

	if err := validation.ValidateStruct(&createBid); err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	bid := utils.MustTransformObj[CreateBid, entity.Bid](&createBid)

	resp, err := c.bidUsecase.CreateBid(ctx, bid)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) GetMyBids(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pagination, err := parsers.ParsePagination(r)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, _ := parsers.ParseQuery(r, "username", false, parsers.ParserEmptyString)

	resp, err := c.bidUsecase.GetMyBids(ctx, username, pagination)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) GetTenderBidsList(w http.ResponseWriter, r *http.Request) {
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

	pagination, err := parsers.ParsePagination(r)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.bidUsecase.GetTenderBidsList(ctx, username, tenderID, pagination)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bidID, err := parsers.ParseVar(r, "bidId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.bidUsecase.GetBidStatus(ctx, username, bidID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bidID, err := parsers.ParseVar(r, "bidId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	status, err := parsers.ParseQuery(r, "status", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}
	if err := validation.ValidateOneOf(
		[]entity.BidStatusType{
			entity.BCreated,
			entity.BPublished,
			entity.BCanceled,
		},
		status, "status"); err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.bidUsecase.UpdateBidStatus(ctx, username, bidID, entity.BidStatusType(status))
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) PatchBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bidID, err := parsers.ParseVar(r, "bidId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	var patchBid PatchBid
	if err := json.NewDecoder(r.Body).Decode(&patchBid); err != nil {
		responses.ErrorHandler(w, validation.ErrParsed)
		return
	}

	if err := validation.ValidateStruct(patchBid); err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	patchBidEnt := utils.MustTransformObj[PatchBid, entity.Bid](&patchBid)

	resp, err := c.bidUsecase.PatchBid(ctx, username, bidID, patchBidEnt)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) SubmitDecisionBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bidID, err := parsers.ParseVar(r, "bidId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	decision, err := parsers.ParseQuery(r, "decision", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}
	if err := validation.ValidateOneOf(entity.BidDecisionTypeList, decision, "decision"); err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.bidUsecase.SubmitDecision(ctx, username, bidID, entity.BidDecisionType(decision))
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) FeedbackBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bidID, err := parsers.ParseVar(r, "bidId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	feedback, err := parsers.ParseQuery(r, "bidFeedback", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.bidUsecase.FeedbackBid(ctx, username, bidID, feedback)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) RollbackBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bidID, err := parsers.ParseVar(r, "bidId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	version, err := parsers.ParseVar(r, "version", true, parsers.ParserInt)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	username, err := parsers.ParseQuery(r, "username", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.bidUsecase.RollbackBid(ctx, username, bidID, version)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}

func (c *Controller) PrevRewiews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenderID, err := parsers.ParseVar(r, "tenderId", true, parsers.ParserUUID)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	authorUsername, err := parsers.ParseQuery(r, "authorUsername", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	requesterUsername, err := parsers.ParseQuery(r, "requesterUsername", true, parsers.ParserEmptyString)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	pagination, err := parsers.ParsePagination(r)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	resp, err := c.bidUsecase.CheckPrevFeedbacks(ctx, tenderID, authorUsername, requesterUsername, *pagination)
	if err != nil {
		responses.ErrorHandler(w, err)
		return
	}

	responses.OkJSON(w, http.StatusOK, resp)
}
