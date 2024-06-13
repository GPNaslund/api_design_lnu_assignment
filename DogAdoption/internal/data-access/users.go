package dataaccess

import (
	userdto "1dv027/aad/internal/dto/user"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDataAccess struct {
	dbPool *pgxpool.Pool
}

func NewUserDataAccess(dbPool *pgxpool.Pool) UserDataAccess {
	return UserDataAccess{
		dbPool: dbPool,
	}
}

func (u UserDataAccess) CreateNewUser(ctx context.Context, newUser userdto.NewUserDTO) (model.User, error) {
	query := `INSERT INTO Users (username, password) VALUES ($1, $2) RETURNING id, username, password`
	var user model.User
	err := u.dbPool.QueryRow(ctx, query, &newUser.Username, &newUser.Password).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return model.User{}, &customerrors.DatabaseError{}
	}
	return user, nil
}

func (u UserDataAccess) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	emptyModel := model.User{}
	query := `SELECT * FROM Users WHERE username = $1`
	var user model.User
	err := u.dbPool.QueryRow(ctx, query, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyModel, &customerrors.UserNotFoundError{}
		}
		return emptyModel, &customerrors.DatabaseError{}
	}
	return user, nil
}

func (u UserDataAccess) DeleteUser(ctx context.Context, userId int) error {
	query := `DELETE FROM Users WHERE id = $1`
	result, err := u.dbPool.Exec(ctx, query, userId)
	if err != nil {
		return &customerrors.DatabaseError{}
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return &customerrors.UserNotFoundError{}
	}
	return nil
}
