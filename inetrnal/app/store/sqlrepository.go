package store

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"project/inetrnal/app/model"
)

type SqlRepository struct {
	db *sqlx.DB
}

func (s *SqlRepository) GetSomeDataFromDB(ctx context.Context, parameter1 string, parameter2 int) (*model.DbModel, error) {
	query := `SELECT first_field, second_field FROM the_table WHERE parameter1 = $1 AND parameter2 = $2`

	result := &model.DbModel{}
	err := s.db.QueryRowxContext(ctx, query, parameter1, parameter2).StructScan(result)
	if err != nil {
		return nil, fmt.Errorf("GetSomeDataFromDB: %v", err)
	}

	return result, nil
}
