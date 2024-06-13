package dataaccess

import (
	"1dv027/aad/internal/dto"
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DogShelterQueries struct {
	dogShelterQuery string
	totalCountQuery string
	filterValues    []any
}

type DogSheltersDataAccess struct {
	dbPool *pgxpool.Pool
}

func NewDogSheltersDataAccess(dbPool *pgxpool.Pool) DogSheltersDataAccess {
	return DogSheltersDataAccess{
		dbPool: dbPool,
	}
}

func (d DogSheltersDataAccess) GetDogShelterByUsername(ctx context.Context, username string) (model.DogShelter, error) {
	emptyModel := model.DogShelter{}
	var dogShelter model.DogShelter
	query := `SELECT * FROM DogShelters WHERE username = $1;`
	err := d.dbPool.QueryRow(ctx, query, username).Scan(
		&dogShelter.Id,
		&dogShelter.Name,
		&dogShelter.Website,
		&dogShelter.Country,
		&dogShelter.City,
		&dogShelter.Address,
		&dogShelter.Username,
		&dogShelter.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyModel, &customerrors.DogShelterNotFoundError{}
		}
		return emptyModel, err
	}
	return dogShelter, nil
}

func (d DogSheltersDataAccess) GetDogShelters(ctx context.Context, queryParams dto.QueryParams) (dogshelterdto.GetDogSheltersQueryResponseDTO, error) {
	emptyDto := dogshelterdto.GetDogSheltersQueryResponseDTO{}
	queries := d.createQuery(queryParams)
	rows, err := d.dbPool.Query(ctx, queries.dogShelterQuery, queries.filterValues...)
	if err != nil {
		return emptyDto, &customerrors.DatabaseError{}
	}

	shelters, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.DogShelter, error) {
		shelter, err := d.dogShelterScanner(row)
		return shelter, err
	})

	if err != nil {
		return emptyDto, &customerrors.DatabaseError{Message: "getDogShelters had a database error"}
	}

	var totalCount int
	err = d.dbPool.QueryRow(ctx, queries.totalCountQuery, queries.filterValues...).Scan(&totalCount)
	if err != nil {
		return emptyDto, err
	}

	dogQueryResult := dogshelterdto.GetDogSheltersQueryResponseDTO{
		DogShelters:          shelters,
		TotalAmountAvailable: totalCount,
	}

	return dogQueryResult, nil
}

func (d DogSheltersDataAccess) GetDogShelterById(ctx context.Context, shelterId int) (model.DogShelter, error) {
	emptyModel := model.DogShelter{}
	query := `SELECT * FROM DogShelters WHERE id = $1`
	row, err := d.dbPool.Query(ctx, query, shelterId)
	if err != nil {
		return model.DogShelter{}, &customerrors.DatabaseError{}
	}
	dogShelter, err := pgx.CollectRows(row, func(row pgx.CollectableRow) (model.DogShelter, error) {
		shelter, err := d.dogShelterScanner(row)
		return shelter, err
	})

	if err != nil {
		return emptyModel, &customerrors.DatabaseError{}
	}

	if len(dogShelter) == 0 {
		return emptyModel, &customerrors.DogShelterNotFoundError{}
	}
	return dogShelter[0], nil
}

func (d DogSheltersDataAccess) DeleteDogShelter(ctx context.Context, shelterId int) error {
	query := `DELETE FROM DogShelters WHERE id = $1`
	result, err := d.dbPool.Exec(ctx, query, shelterId)
	if err != nil {
		return &customerrors.DatabaseError{}
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return &customerrors.DogShelterNotFoundError{}
	} else {
		return nil
	}
}

func (d DogSheltersDataAccess) UpdateDogShelter(ctx context.Context,
	shelterId int, dogShelterData dogshelterdto.UpdateDogShelterDTO) error {
	query := d.createUpdateDogShelterQuery(shelterId, dogShelterData)
	result, err := d.dbPool.Exec(ctx, query)
	if err != nil {
		return &customerrors.DatabaseError{}
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return &customerrors.DogShelterNotFoundError{}
	}
	return nil
}

func (d DogSheltersDataAccess) CreateDogShelter(ctx context.Context, newShelter dogshelterdto.NewDogShelterDTO) (int, error) {
	query := `INSERT INTO DogShelters 
	(name, website, country, city, address, username, password)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var id int
	err := d.dbPool.QueryRow(ctx, query,
		*newShelter.Name,
		*newShelter.Website,
		*newShelter.Country,
		*newShelter.City,
		*newShelter.Address,
		*newShelter.Username,
		*newShelter.Password,
	).Scan(&id)
	if err != nil {
		return 0, &customerrors.DatabaseError{}
	}
	return id, nil
}

func (d DogSheltersDataAccess) createQuery(queryParams dto.QueryParams) DogShelterQueries {
	qb := NewQueryBuilder("SELECT * FROM DogShelters")

	dogShelterFilters := queryParams.DogShelterFilter

	if dogShelterFilters.Country != nil {
		qb.withFilterParam("country", *dogShelterFilters.Country)
	}
	if dogShelterFilters.City != nil {
		qb.withFilterParam("city", *dogShelterFilters.City)
	}
	if dogShelterFilters.Name != nil {
		qb.withFilterParam("name", *dogShelterFilters.Name)
	}

	paginationParams := queryParams.Pagination
	if paginationParams != nil {
		if paginationParams.Limit != nil {
			qb.withLimit(*paginationParams.Limit)
		}
		if paginationParams.Page != nil {
			qb.withPage(*paginationParams.Page)
		}
	}

	selectQuery, filterValues := qb.build()
	countQuery, _ := qb.buildForTotalCount()
	queries := DogShelterQueries{
		dogShelterQuery: selectQuery,
		totalCountQuery: countQuery,
		filterValues:    filterValues,
	}
	return queries
}

func (d DogSheltersDataAccess) dogShelterScanner(row pgx.CollectableRow) (model.DogShelter, error) {
	var dogShelter model.DogShelter
	err := row.Scan(
		&dogShelter.Id,
		&dogShelter.Name,
		&dogShelter.Website,
		&dogShelter.Country,
		&dogShelter.City,
		&dogShelter.Address,
		&dogShelter.Username,
		&dogShelter.Password,
	)
	return dogShelter, err
}

func (d DogSheltersDataAccess) createUpdateDogShelterQuery(dogShelterId int, dogShelter dogshelterdto.UpdateDogShelterDTO) string {
	query := `UPDATE DogShelters SET `
	updates := []string{}

	addUpdate := func(field string, value interface{}) {
		switch v := value.(type) {
		case string:
			updates = append(updates, fmt.Sprintf("%s = '%s'", field, v))
		case int, int64:
			updates = append(updates, fmt.Sprintf("%s = %d", field, v))
		}
	}

	if dogShelter.Name != nil {
		addUpdate("name", *dogShelter.Name)
	}
	if dogShelter.Website != nil {
		addUpdate("website", *dogShelter.Website)
	}
	if dogShelter.Country != nil {
		addUpdate("country", dogShelter.Country)
	}
	if dogShelter.City != nil {
		addUpdate("city", *dogShelter.City)
	}
	if dogShelter.Address != nil {
		addUpdate("address", *dogShelter.Address)
	}

	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = %d;", dogShelterId)

	return query
}
