package models

import (
	"github.com/Zetkolink/ckassa/models"
)

// Model модель для работы с данными карт.
type CardCallback struct {
	models.CardCallback
}

func (c CardCallback) TableName() string {
	return "cashless.card_callbacks"
}
