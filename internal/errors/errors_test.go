package errors 

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	var err error

	verify := Verifier{}

	one := 1
	verify.That(one == 1, "1 should equal 1")

	err = verify.Flush()
	assert.NoError(t, err)

	verify.That(1 == 2, "1 should equal 2")
	err = verify.Flush()
	assert.Error(t, err)
}
