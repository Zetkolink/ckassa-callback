package cardCallback

import (
	"../../models"
	"context"
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Store struct {
	*gorm.DB
}

func NewCardCallbacksStore(db *gorm.DB) *Store {
	return &Store{db}
}

// Insert запись уведомления о добавлени карты.
func (s Store) InsertCard(ctx context.Context, card models.CardCallback) error {
	tx := s.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&card).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
