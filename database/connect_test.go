package database 

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectToDatabase(t *testing.T) {
	db_conn := Connect("file::memory:?cache=shared")
	assert.NotNil(t, db_conn)
}
