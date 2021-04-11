package repository

import (
	"context"
	"pandita/model"

	"github.com/juju/errors"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	Conn *gorm.DB
}

// NewGormUserRepository ...
func NewGormUserRepository(conn *gorm.DB) UserRepository {
	migrations := []interface{}{
		&model.User{},
	}
	if err := conn.Migrator().AutoMigrate(migrations...); err != nil {
		mlog.Panicw("Unable to AutoMigrate UserRepository", "error", err)
	}
	return &gormUserRepository{Conn: conn}
}

// GetUserByID
func (g *gormUserRepository) GetUserByID(ctx context.Context, uid uint64) (user *model.User, err error) {
	scope := g.Conn.WithContext(ctx)
	scope = scope.Where("uid = ?", uid).Find(&user)
	if scope.RowsAffected == 0 {
		return nil, errors.NotFoundf("userID [%d]", uid)
	}
	return user, nil
}
