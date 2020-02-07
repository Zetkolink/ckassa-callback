package models

import (
	"github.com/Zetkolink/ckassa/models"
)

// Model модель для работы с данными карт.
type PaymentCallback struct {
	models.PaymentCallback
}

func (c PaymentCallback) TableName() string {
	return "cashless.payments"
}
