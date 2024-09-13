package main

import (
	"avito/internal/config"
	"avito/internal/db/models"
	"avito/internal/db/repos"
	"avito/internal/entity"
	"avito/internal/usecases"
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateTenders() {
	cfg := config.LoadEnv()

	db, _ := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DB.PostgresConn,
	}), &gorm.Config{})

	var orgs []models.Organization
	db.Find(&orgs)

	tenderRepo, _ := repos.NewTenderRepo(&cfg.DB)

	for range 5 {
		tenderRepo.CreateTender(context.Background(), &entity.Tender{
			Name:        gofakeit.Name(),
			Description: gofakeit.Paragraph(1, 3, 4, " "),
			ServiceType: entity.TenderServiceTypeList[gofakeit.IntN(len(entity.TenderServiceTypeList))],
			Status:      entity.Published,
			// Status:         entity.TenderStatusTypeList[gofakeit.IntN(len(entity.TenderServiceTypeList))],
			OrganizationID: orgs[gofakeit.IntN(len(orgs))].Id,
			Version:        1,
		})
	}
}

func CreateBids() {
	cfg := config.LoadEnv()

	db, _ := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DB.PostgresConn,
	}), &gorm.Config{})

	tenderRepo, _ := repos.NewTenderRepo(&cfg.DB)
	bidsRepo, _ := repos.NewBidRepo(&cfg.DB)

	tenderUsecase := usecases.NewTenderUsecase(tenderRepo)
	bidsUsecase := usecases.NewBidUsecase(tenderRepo, bidsRepo, tenderUsecase)

	var tenders []models.Tender
	db.Find(&tenders)

	var orgs []models.Organization
	db.Find(&orgs)

	var users []models.User
	db.Find(&users)

	for range 4 {
		bidsUsecase.CreateBid(context.Background(), &entity.Bid{
			Name:        gofakeit.Name(),
			Description: gofakeit.Paragraph(1, 3, 4, " "),
			// Status:      entity.BidStatusTypeList[gofakeit.IntN(len(entity.BidStatusTypeList))],
			TenderID:   tenders[gofakeit.IntN(len(tenders))].Id,
			AuthorType: entity.AuthorOrganization,
			AuthorID:   orgs[gofakeit.IntN(len(orgs))].Id,
			// Version:    1,
		})
	}

	for range 3 {
		bidsUsecase.CreateBid(context.Background(), &entity.Bid{
			Name:        gofakeit.Name(),
			Description: gofakeit.Paragraph(1, 3, 4, " "),
			// Status:      entity.BidStatusTypeList[gofakeit.IntN(len(entity.BidStatusTypeList))],
			TenderID:   tenders[gofakeit.IntN(len(tenders))].Id,
			AuthorType: entity.AuthorUser,
			AuthorID:   users[gofakeit.IntN(len(users))].Id,
			// Version:    1,
			// Kvorum: 1,
		})
	}
}

// func UsersHaveTenders() {
// 	cfg := config.LoadEnv()

// 	db, _ := gorm.Open(postgres.New(postgres.Config{
// 		DSN: cfg.DB.PostgresConn,
// 	}), &gorm.Config{})

// 	tenderRepo, _ := repos.NewTenderRepo(&cfg.DB)
// 	var tenders []models.Tender
// 	db.Find(&tenders)

// 	for _, t := range tenders {

// 	}
// }

func main() {
	// CreateTenders()
	// CreateBids()
}
