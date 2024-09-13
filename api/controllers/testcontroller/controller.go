package testcontroller

import (
	"avito/api/responses"
	"avito/internal/config"
	"avito/internal/db/models"
	"fmt"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Testcontroller struct {
	db *gorm.DB
}

func NewTestcontroller(cfg *config.DB) (*Testcontroller, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.PostgresConn,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("create db gorm obj: %w", err)
	}

	return &Testcontroller{
		db: db,
	}, nil
}

func (c *Testcontroller) GetData(w http.ResponseWriter, r *http.Request) {
	usrs := []models.User{}
	orgs := []models.Organization{}
	usr_orgs := []models.OrganizationResponsible{}

	c.db.Find(&usrs)
	c.db.Find(&orgs)
	c.db.Find(&usr_orgs)

	responses.OkJSON(w, http.StatusOK, []any{usrs, orgs, usr_orgs})
}
