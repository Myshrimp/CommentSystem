package data

import (
	"fmt"
	"review-service/internal/conf"
	"review-service/internal/data/query"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewReviewRepo, NewDB)

// Data .
type Data struct {
	// TODO wrapped database client
	query *query.Query
	log *log.Helper
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}

	// Important! Set database for generated ORM code
	query.SetDefault(db)

	return &Data{query: query.Q, log: log.NewHelper(logger)}, cleanup, nil
}

func NewDB(c *conf.Data) (*gorm.DB, error) {
	switch strings.ToLower(c.Database.GetDriver()) {
		case "mysql":
			return gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
		case "sqlite":
			return gorm.Open(sqlite.Open(c.Database.Source), &gorm.Config{})
		default:
			return nil, fmt.Errorf("unsupported database driver: %s", c.Database.GetDriver())
	}
}
