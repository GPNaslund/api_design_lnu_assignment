package dataaccess

import (
	"fmt"
	"strings"
)

type queryBuilder struct {
	baseQuery    string
	queryFilters []string
	filterValues []any
	limit        int
	offset       int
}

func NewQueryBuilder(baseQuery string) queryBuilder {
	return queryBuilder{
		baseQuery: baseQuery,
	}
}

func (q *queryBuilder) withFilterParam(key, value string) {
	filterIndex := len(q.queryFilters) + 1
	q.queryFilters = append(q.queryFilters, fmt.Sprintf("%s=$%d", key, filterIndex))
	q.filterValues = append(q.filterValues, value)
}

func (q *queryBuilder) withLimit(limit int) {
	q.limit = limit
}

func (q *queryBuilder) withPage(page int) {
	q.offset = (page - 1) * q.limit
}

func (q *queryBuilder) build() (string, []any) {
	finalQuery := q.baseQuery

	if len(q.queryFilters) > 0 {
		finalQuery += " WHERE " + strings.Join(q.queryFilters, " AND ")
	}

	if q.limit > 0 {
		finalQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", q.limit, q.offset)
	}

	finalQuery += ";"
	return finalQuery, q.filterValues
}

func (q *queryBuilder) buildForTotalCount() (string, []any) {
	countQuery := strings.Replace(q.baseQuery, "SELECT *", "SELECT COUNT(*)", 1)

	if len(q.queryFilters) > 0 {
		countQuery += " WHERE " + strings.Join(q.queryFilters, " AND ")
	}

	countQuery += ";"

	return countQuery, q.filterValues
}
