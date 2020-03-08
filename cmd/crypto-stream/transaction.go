package main

import (
	"encoding/json"
	"time"
)

type TransactionOut struct {
	Value int64 `json:"value"`
}
type TrX struct {
	Time int64            `json:"time"`
	Out  []TransactionOut `json:"out"`
}

// Transaction represents crypto transaction.
type Transaction struct {
	X TrX `json:"x"`
}

func transactionFromBytes(b []byte) (Transaction, error) {
	var tr Transaction
	err := json.Unmarshal(b, &tr)
	if err != nil {
		return Transaction{}, err
	}

	return tr, nil
}

// GetUSDValue calculates the USD value of transaction.
func (t Transaction) GetUSDValue() float64 {
	var val float64 = 0

	for _, outT := range t.X.Out {
		val += float64(outT.Value) / float64(100000000)
	}

	return val
}

// GetDate returns the date transaction was made.
func (t Transaction) GetDate() string {
	return time.Unix(t.X.Time, 0).Format(time.RFC3339)
}
