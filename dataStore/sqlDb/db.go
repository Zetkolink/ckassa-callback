package sqlDb

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Connect initialize connect to DB.
func Connect(connLine string, driverName string) (*gorm.DB, error) {
	db, err := gorm.Open(driverName, connLine)
	if err != nil {
		return db, err
	}

	return db, nil
}
