package budget

import (
	"log"
	"my-budget/internal/testutils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBudget(t *testing.T) {
	// Todo: Add more test around edge cases and points of failure
	// ! Currently this test is just going down the happy path
	dt := testutils.NewDomainTest(t)
	println(dt)

	file, err := os.Open("/Users/appstack/Developer/Personal/budget/cmd/generateTestData/Chase Activity Oct 6.CSV")
	assert.NoError(t, err)

	budget, err := NewBudget(dt.Ctx, dt.Conn)
	if err != nil {
		log.Fatal(err)
	}

	err = budget.AddNewTransactionsFromDocument("Chase Activity Oct 6.CSV", file)
	assert.NoError(t, err)

	err = budget.GetIncomeVsExpense()
	assert.NoError(t, err)
}
