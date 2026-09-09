package handler

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/errs"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/store"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/validation"
	"github.com/google/uuid"
	"github.com/starfederation/datastar-go/datastar"
)

type TradeHandler struct {
	service service.TransactionService
}

func NewTradeHandler(service service.TransactionService) *TradeHandler {
	return &TradeHandler{service}
}

func (h *TradeHandler) CreateTrade(w http.ResponseWriter, r *http.Request) {
	// parse datastar signals
	var payload struct {
		CreateTradeSignals store.CreateTradeSignals `json:"create_trade_signals"`
	}
	if err := datastar.ReadSignals(r, &payload); err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}
	form := payload.CreateTradeSignals
	// logger.Debug("VIBE CHECK form: ", form)

	// // string conversions
	leaseIdUuid, err := uuid.Parse(form.LeaseId)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	// validation
	if err := validation.CheckCreateTradeForm(form); err != nil {
		// TODO: use datastar to show errors in the front end
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	_, err = h.service.CollectRent(r.Context(), service.CollectRentParams{
		LeaseId: leaseIdUuid,
	})
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)

	sse.Redirect("/dashboard")
}
