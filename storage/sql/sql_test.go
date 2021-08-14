package sql

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
)

type idS struct {
	HER string `db:"slot_id"`
}

func TestS(t *testing.T) {
	ctx := context.Background()
	dsn := "host=localhost port=5432 user=danny database=rotation password=danny sslmode=disable"
	s := NewStorage("pgx", dsn)
	err := s.Connect(ctx)
	require.NoError(t, err)

	rows, err := s.db.NamedQueryContext(ctx,
		"INSERT INTO slots (slot_description) VALUES (:description) RETURNING slot_id",
		map[string]interface{}{"description": "123123123"},
	)
	require.NoError(t, err)
	for rows.Next() {
		ids := new(idS)
		err = rows.StructScan(ids)
		require.NoError(t, err)
		fmt.Printf("%#v", ids)
	}
}
