package mempool_test

import (
	"container/heap"
	"errors"
	"sort"
	"testing"

	"github.com/chadit/mempool_transaction/internal/mempool"
	"github.com/stretchr/testify/assert"
)

func TestCalculateFee(t *testing.T) {
	tests := []struct {
		gas      float64
		gasFee   float64
		expected float64
		err      error
	}{
		{gas: -1.0, gasFee: 0.0, expected: 0.0, err: nil},
		{gas: 0.0, gasFee: 0.0, expected: 0.0, err: nil},
		{gas: 10.0, gasFee: 0.0, expected: 0.0, err: nil},
		{gas: -10.0, gasFee: 10.0, expected: 0.0, err: errors.New("cannot be less than zero")},
		{gas: 10.0, gasFee: 10.0, expected: 100.0, err: nil},
		{gas: 20.0, gasFee: 10.0, expected: 200.0, err: nil},
	}

	mp := mempool.New()
	for _, test := range tests {
		fees, err := mp.CalculateFee(test.gas, test.gasFee)
		if test.err == nil {
			assert.Nil(t, err, "expected nil error")
		} else {
			assert.Equal(t, test.err.Error(), err.Error(), "expected error should be equal")
		}
		assert.Equal(t, test.expected, fees, "should equal")
	}
}

func TestParseAndCalculate(t *testing.T) {
	const (
		tnxFile   = "../transactions.txt"
		poolLimit = 5000
	)

	tests := []struct {
		tnxFile   string
		poolLimit int
		err       error
	}{
		{tnxFile: "../transactions.txt", poolLimit: 7000, err: nil},
		{tnxFile: "../transactions.txt", poolLimit: 5000, err: nil},
		{tnxFile: "../transactions.txt", poolLimit: 500, err: nil},
	}

	for _, test := range tests {
		mp := mempool.New()
		err := mp.Parse(test.tnxFile, test.poolLimit)
		assert.Nil(t, err)

		arr := make([]float64, 0)
		for len(*mp) > 0 {
			arr = append(arr, heap.Pop(mp).(*mempool.Transaction).Fee)
		}

		countArr := len(arr)
		assert.Equal(t, test.poolLimit, countArr)

		sort.Float64s(arr)

		expectedHighest := arr[0]
		for i := 1; i < countArr; i++ {
			assert.Less(t, expectedHighest, arr[i], "first item in the slice should be the highest")
		}
	}
}
