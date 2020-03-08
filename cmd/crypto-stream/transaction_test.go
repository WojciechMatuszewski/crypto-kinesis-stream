package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_GetUSDValue(t *testing.T) {
	t.Run("empty transaction list", func(t *testing.T) {
		tr := Transaction{TrX{
			Time: 0,
			Out:  []TransactionOut{},
		}}

		out := tr.GetUSDValue()
		assert.Equal(t, float64(0), out)
	})
}
