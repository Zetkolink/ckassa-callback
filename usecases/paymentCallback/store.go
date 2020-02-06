package paymentCallback

import (
	"../../models"
	"context"
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Store struct {
	*gorm.DB
}

func NewPaymentCallbacksStore(db *gorm.DB) *Store {
	return &Store{db}
}

// Insert запись уведомления о платеже.
func (s Store) InsertPayment(ctx context.Context, payment models.PaymentCallback) error {
	tx := s.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
