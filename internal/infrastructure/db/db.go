package db

import (
	"context"
	"fmt"
	"sync"

	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
	"github.com/pimp13/jira-clone-backend-go/ent"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/config"
)

var (
	entClient *ent.Client
	once      sync.Once
	initErr   error
)

func NewEntClient(cfg *config.DB) (*ent.Client, error) {
	once.Do(func() {
		connectionString := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Pass,
			cfg.Name,
		)
		client, err := ent.Open(cfg.Connection, connectionString)
		if err != nil {
			initErr = err
			return
		}

		if err := client.Schema.Create(
			context.Background(),
			schema.WithDropColumn(true), // ستون‌های قدیمی حذف میشن
			schema.WithDropIndex(true),  // ایندکس‌های قدیمی حذف میشن
			schema.WithIndent(" "),
		); err != nil {
			initErr = err
			return
		}

		entClient = client
	})

	return entClient, initErr
}
