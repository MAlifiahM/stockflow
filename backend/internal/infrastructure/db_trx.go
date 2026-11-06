package infrastructure

import (
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"gorm.io/gorm"
)

type gormTransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) domain.DBTransactionManager {
	return &gormTransactionManager{db: db}
}

func (t *gormTransactionManager) Execute(fn func(tx *gorm.DB) error) (err error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = fn(tx)

	return err
}
