package main

import (
	"avito/api/controllers/bid"
	"avito/api/controllers/ping"
	"avito/api/controllers/tender"
	"avito/internal/config"
	"avito/internal/db/repos"
	"avito/internal/usecases"
	"fmt"
	"net/http"

	"context"
	"time"

	"github.com/gorilla/mux"
)

func ctxTimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.LoadEnv()

	tenderRepo, err := repos.NewTenderRepo(&cfg.DB)
	if err != nil {
		panic(fmt.Errorf("create repo: %w", err))
	}

	bidRepo, err := repos.NewBidRepo(&cfg.DB)
	if err != nil {
		panic(fmt.Errorf("create repo: %w", err))
	}

	tenderUsecase := usecases.NewTenderUsecase(tenderRepo)
	bidUsecase := usecases.NewBidUsecase(tenderRepo, bidRepo, tenderUsecase)

	pingController := ping.Controller{}
	tenderController := tender.NewTenderController(tenderUsecase)
	bidController := bid.NewBidController(bidUsecase)

	r := mux.NewRouter()
	api := r.PathPrefix("/api/").Subrouter()
	api.Use(ctxTimeoutMiddleware)
	// api.Use(dbgMiddleware)

	api.HandleFunc("/ping", pingController.Ping).Methods("GET")

	api.HandleFunc("/tenders/{tenderId}/rollback/{version}", tenderController.RollbackTender).Methods("PUT")
	api.HandleFunc("/tenders/{tenderId}/status", tenderController.GetTenderStatus).Methods("GET")
	api.HandleFunc("/tenders/{tenderId}/status", tenderController.UpdateTenderStatus).Methods("PUT")
	api.HandleFunc("/tenders/{tenderId}/edit", tenderController.PatchTender).Methods("PATCH")
	api.HandleFunc("/tenders/new", tenderController.CreateTender).Methods("POST")
	api.HandleFunc("/tenders/my", tenderController.GetMyTenders).Methods("GET")
	api.HandleFunc("/tenders", tenderController.GetTenders).Methods("GET")

	api.HandleFunc("/bids/{bidId}/rollback/{version}", bidController.RollbackBid).Methods("PUT")
	api.HandleFunc("/bids/{tenderId}/reviews", bidController.PrevRewiews).Methods("GET")
	api.HandleFunc("/bids/{bidId}/feedback", bidController.FeedbackBid).Methods("PUT")
	api.HandleFunc("/bids/{bidId}/submit_decision", bidController.SubmitDecisionBid).Methods("PUT")
	api.HandleFunc("/bids/{bidId}/edit", bidController.PatchBid).Methods("PATCH")
	api.HandleFunc("/bids/{bidId}/status", bidController.UpdateBidStatus).Methods("PUT")
	api.HandleFunc("/bids/{bidId}/status", bidController.GetBidStatus).Methods("GET")
	api.HandleFunc("/bids/{tenderId}/list", bidController.GetTenderBidsList).Methods("GET")
	api.HandleFunc("/bids/my", bidController.GetMyBids).Methods("GET")
	api.HandleFunc("/bids/new", bidController.CreateBid).Methods("POST")

	http.ListenAndServe(cfg.Server.ServerAddress, r)
}
