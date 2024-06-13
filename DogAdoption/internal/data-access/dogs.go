package dataaccess

import (
	dto "1dv027/aad/internal/dto"
	dogdto "1dv027/aad/internal/dto/dog"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DogQueries struct {
	dogsQuery       string
	totalCountQuery string
	filterValues    []any
}

type DogsDataAccess struct {
	dbPool *pgxpool.Pool
}

func NewDogsDataAccess(dbPool *pgxpool.Pool) DogsDataAccess {
	return DogsDataAccess{
		dbPool: dbPool,
	}
}

func (d DogsDataAccess) GetDogs(ctx context.Context, queryParams dto.QueryParams) (dogdto.GetDogsQueryResponseDTO, error) {
	emptyDto := dogdto.GetDogsQueryResponseDTO{}
	queries := d.createQuery(queryParams)
	rows, err := d.dbPool.Query(ctx, queries.dogsQuery, queries.filterValues...)
	if err != nil {
		return emptyDto, &customerrors.DatabaseError{}
	}

	dogs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Dog, error) {
		dog, err := d.dogScanner(row)
		return dog, err
	})

	if err != nil {
		return emptyDto, &customerrors.DatabaseError{}
	}

	var totalCount int
	err = d.dbPool.QueryRow(ctx, queries.totalCountQuery, queries.filterValues...).Scan(&totalCount)
	if err != nil {
		return emptyDto, err
	}

	dogQueryResult := dogdto.GetDogsQueryResponseDTO{
		Dogs:                 dogs,
		TotalAmountAvailable: totalCount,
	}

	return dogQueryResult, nil
}

func (d DogsDataAccess) GetDogById(ctx context.Context, dogId int) (model.Dog, error) {
	emptyModel := model.Dog{}
	query := `SELECT * FROM Dogs WHERE id = $1`
	row, err := d.dbPool.Query(ctx, query, dogId)
	if err != nil {
		return emptyModel, &customerrors.DatabaseError{}
	}
	dogData, err := pgx.CollectRows(row, func(row pgx.CollectableRow) (model.Dog, error) {
		dog, err := d.dogScanner(row)
		return dog, err
	})

	if err != nil {
		return emptyModel, &customerrors.DatabaseError{}
	}

	if len(dogData) == 0 {
		return emptyModel, &customerrors.DogNotFoundError{}
	}
	return dogData[0], nil
}

func (d DogsDataAccess) DeleteDog(ctx context.Context, dogId int) error {
	query := `DELETE FROM Dogs WHERE id = $1`
	result, err := d.dbPool.Exec(ctx, query, dogId)
	if err != nil {
		return &customerrors.DatabaseError{}
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return &customerrors.DogNotFoundError{}
	} else {
		return nil
	}
}

func (d DogsDataAccess) UpdateDog(ctx context.Context, dogId int, dogData dogdto.UpdateDogDTO) error {
	query := d.createUpdateDogQuery(dogId, dogData)
	result, err := d.dbPool.Exec(ctx, query)
	if err != nil {
		return &customerrors.DatabaseError{}
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return &customerrors.DogNotFoundError{}
	}
	return nil
}

func (d DogsDataAccess) createUpdateDogQuery(dogId int, dogData dogdto.UpdateDogDTO) string {
	query := `UPDATE Dogs SET `
	updates := []string{}

	addUpdate := func(field string, value interface{}) {
		switch v := value.(type) {
		case string:
			updates = append(updates, fmt.Sprintf("%s = '%s'", field, v))
		case int, int64:
			updates = append(updates, fmt.Sprintf("%s = %d", field, v))
		case bool:
			updates = append(updates, fmt.Sprintf("%s = %t", field, v))
		case *time.Time:
			updates = append(updates, fmt.Sprintf("%s = '%s'", field, v.Format("2006-01-02 15:04:05")))
		}
	}

	if dogData.Name != nil {
		addUpdate("name", *dogData.Name)
	}
	if dogData.Description != nil {
		addUpdate("description", *dogData.Description)
	}
	if dogData.BirthDate != nil {
		addUpdate("birth_date", dogData.BirthDate.Format("2006-01-02"))
	}
	if dogData.Breed != nil {
		addUpdate("breed", *dogData.Breed)
	}
	if dogData.IsNeutered != nil {
		addUpdate("is_neutered", *dogData.IsNeutered)
	}
	if dogData.ImageUrl != nil {
		addUpdate("image_url", *dogData.ImageUrl)
	}
	if dogData.AdoptionFee != nil {
		addUpdate("adoption_fee", *dogData.AdoptionFee)
	}
	if dogData.IsAdopted != nil {
		addUpdate("is_adopted", *dogData.IsAdopted)
	}
	if dogData.FriendlyWith != nil {
		addUpdate("friendly_with", *dogData.FriendlyWith)
	}
	if dogData.Gender != nil {
		addUpdate("gender", *dogData.Gender)
	}
	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = %d;", dogId)

	return query
}

func (d DogsDataAccess) CreateDog(ctx context.Context, dogData dogdto.NewDogDTO) (int, error) {
	query := `INSERT INTO Dogs 
	(name, description, birth_date, breed, is_neutered, shelter_id, image_url, adoption_fee, is_adopted, friendly_with, gender)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	var id int
	err := d.dbPool.QueryRow(ctx, query,
		*dogData.Name,
		*dogData.Description,
		*dogData.BirthDate,
		*dogData.Breed,
		*dogData.IsNeutered,
		*dogData.ShelterId,
		*dogData.ImageUrl,
		*dogData.AdoptionFee,
		*dogData.IsAdopted,
		*dogData.FriendlyWith,
		*dogData.Gender,
	).Scan(&id)
	if err != nil {
		return 0, &customerrors.DatabaseError{}
	}
	return id, nil
}

func (d DogsDataAccess) createQuery(queryParams dto.QueryParams) DogQueries {
	qb := NewQueryBuilder("SELECT * FROM Dogs")

	dogFilters := queryParams.DogsFilter
	if dogFilters.Breed != nil {
		qb.withFilterParam("breed", *dogFilters.Breed)
	}
	if dogFilters.Gender != nil {
		qb.withFilterParam("gender", *dogFilters.Gender)
	}
	if dogFilters.IsAdopted != nil {
		qb.withFilterParam("is_adopted", *dogFilters.IsAdopted)
	}

	if dogFilters.IsNeutered != nil {
		qb.withFilterParam("is_neutered", *dogFilters.IsNeutered)
	}
	if dogFilters.ShelterId != nil {
		qb.withFilterParam("shelter_id", fmt.Sprintf("%d", *dogFilters.ShelterId))
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
	queries := DogQueries{
		dogsQuery:       selectQuery,
		totalCountQuery: countQuery,
		filterValues:    filterValues,
	}
	return queries
}

func (d *DogsDataAccess) dogScanner(row pgx.CollectableRow) (model.Dog, error) {
	var dog model.Dog
	err := row.Scan(
		&dog.Id,
		&dog.Name,
		&dog.Description,
		&dog.BirthDate,
		&dog.Breed,
		&dog.IsNeutered,
		&dog.ShelterId,
		&dog.ImageUrl,
		&dog.AdoptionFee,
		&dog.IsAdopted,
		&dog.FriendlyWith,
		&dog.Gender,
	)
	return dog, err
}
