package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectToDatabase(t *testing.T) {
	ctx := context.Background()
	db_conn, err := ConnectToDatabase(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, db_conn)
}