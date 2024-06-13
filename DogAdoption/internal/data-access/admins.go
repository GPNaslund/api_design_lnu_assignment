package dataaccess

import (
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminsDataAccess struct {
	dbPool *pgxpool.Pool
}

func NewAdminsDataAccess(dbPool *pgxpool.Pool) AdminsDataAccess {
	return AdminsDataAccess{
		dbPool: dbPool,
	}
}

func (a AdminsDataAccess) GetAdminByUsername(ctx context.Context, username string) (model.Admin, error) {
	emptyModel := model.Admin{}
	var admin model.Admin
	query := `SELECT * FROM Admins WHERE username = $1;`
	err := a.dbPool.QueryRow(ctx, query, username).Scan(&admin.Id, &admin.Username, &admin.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyModel, &customerrors.AdminNotFoundError{}
		}
		return emptyModel, err
	}

	return admin, nil
}
