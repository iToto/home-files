package userDAL

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"

	"github.com/jmoiron/sqlx"
)

type DAL interface {
	GetUserByID(ctx context.Context, wl wlog.Logger, id string) (*entities.User, error)
}

type userDAL struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (DAL, error) {
	return &userDAL{
		db: db,
	}, nil
}

func (ud *userDAL) GetUserByID(
	ctx context.Context,
	wl wlog.Logger,
	id string,
) (*entities.User, error) {
	var user entities.User

	query := "SELECT id, name, created_at, updated_at, deleted_at FROM mvp_user WHERE id = $1 AND deleted_at IS NULL"
	err := ud.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user found")
		}
		return nil, err
	}
	return &user, nil
}
