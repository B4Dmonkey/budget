package budget

import (
	"context"
	"database/sql"
	"my-budget/internal/errors"
	"my-budget/internal/testutils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBudget(t *testing.T) {
	dt := testutils.NewDomainTest(t)

	tests := map[string]struct {
		ctx     context.Context
		conn    *sql.DB
		errType error
	}{
		"It Fails when the context is nil": {
			ctx: nil, conn: dt.Conn(), errType: &errors.VerificationError{},
		},
		"It Fails when the connection is nil": {
			ctx: dt.Ctx(), conn: nil, errType: &errors.VerificationError{},
		},
		"It Succeeds when the context and connection are provided": {ctx: dt.Ctx(), conn: dt.Conn(), errType: nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewBudget(tc.ctx, tc.conn)
			dt.Assert.IsType(tc.errType, err)
		})
	}
}

func TestAddNewTransactionsFromDocument(t *testing.T) {
	dt := testutils.NewDomainTest(t)

	budget := &Budget{}

	err := budget.AddNewTransactionsFromDocument("", nil)
	dt.Assert.IsType(&errors.VerificationError{}, err)
  dt.ResetTestState()

	tests := map[string]struct {
		fileNames []string
	}{
		"It passes with one file": {fileNames: []string{"Chase Activity Oct 6.CSV"}},
		"It passes with multiple files": {fileNames: []string{
			"Chase Activity Sept 27.CSV",
			"Chase Activity Oct 6.CSV",
			"Chase Activity Oct 30.CSV",
      // "Chase9931_Activity_20240412.CSV",
		}},
	}

	base_dir := "/Users/appstack/Developer/Personal/budget/cmd/generateTestData"

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			budget, err := NewBudget(dt.Ctx(), dt.Conn())
			assert.NoError(t, err)

			for _, fileName := range tc.fileNames {
				file, err := os.Open(base_dir + "/" + fileName)
				assert.NoError(t, err)

				err = budget.AddNewTransactionsFromDocument(fileName, file)
				assert.NoError(t, err)
			}

			dt.ResetTestState()
		})
	}
}

func TestGetIncomeVsExpense(t *testing.T) {
	dt := testutils.NewDomainTest(t)
	budget := &Budget{}

	err := budget.GetIncomeVsExpense()
	dt.Assert.IsType(&errors.VerificationError{}, err)
}
