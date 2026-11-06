package domain

import "gorm.io/gorm"

type DBTransactionManager interface {
	Execute(fn func(tx *gorm.DB) error) error
}
