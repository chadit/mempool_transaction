package mempool

import (
	"bufio"
	"container/heap"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Transaction represents the structure for the mempool.
type Transaction struct {
	Hash      string
	Signature string
	Gas       string
	Fpg       string
	Fee       float64
	index     int
}

// Mempool stores the transaction.
type Mempool []*Transaction

// New initialzes the mempool for transactions.
func New() *Mempool {
	mp := make(Mempool, 0)
	heap.Init(&mp)
	return &mp
}

// Parse processes the transaction input file into the mempool.
func (mp *Mempool) Parse(transactionFile string, mempoolLimit int) error {
	// open file for reading.
	file, err := os.Open(transactionFile)
	if err != nil {
		errors.Wrap(err, "opening transaction file")
	}

	// place file on scanner for processing.
	scanner := bufio.NewScanner(file)

	// initialize a mempool for parsing, this will also be used to maintian the limit
	parserMempool := New()
	//	returnMP := New()

	// scan through file
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, " ")

		// check that transaction line is not to small.
		if len(items) < 4 {
			return errors.New("invalid line in transaction file")
		}

		// parse the line item from the file.
		txHash := items[0]
		sig := items[3]
		gasString := items[1]
		feePerGasString := items[2]

		gas, err := strconv.ParseFloat(gasString[4:], 64)
		if err != nil {
			return errors.Wrap(err, "parsing gas amount")
		}
		feePerGas, err := strconv.ParseFloat(feePerGasString[10:], 64)
		if err != nil {
			return errors.Wrap(err, "parsing gas fees")
		}

		fee, err := mp.CalculateFee(gas, feePerGas)
		if err != nil {
			return errors.Wrap(err, "parsing fees for gas")
		}

		// calculate fee and populate Fee struct
		txn := &Transaction{
			Hash:      txHash,
			Signature: sig,
			Gas:       gasString,
			Fpg:       feePerGasString,
			Fee:       fee,
		}

		heap.Push(parserMempool, txn)

		parserMempool.reorder(txn, txn.Fee)

		// check if over capacity, if so remove lowest Fee item.
		if len(*parserMempool) > mempoolLimit {
			heap.Remove(parserMempool, len(*parserMempool)-1)
		}
	}

	// once processing of file is complete, return transactions with the highest Fee within the pool limits.
	for mp.Len() < mempoolLimit {
		mp.Push(heap.Pop(parserMempool))
	}
	return nil
}

// CalculateFee returns the fee based on Gas and FeePerGas
func (mp *Mempool) CalculateFee(gas float64, feePerGas float64) (float64, error) {
	if gas == 0 || feePerGas == 0 {
		return 0.0, nil
	}

	if gas < 0 || feePerGas < 0 {
		return 0.0, errors.New("cannot be less than zero")
	}

	return gas * feePerGas, nil
}

// Len needed for sort.Interface, returns number of items.
func (mp Mempool) Len() int { return len(mp) }

// Swap needed for sort.Interface, rearranged the slice.
func (mp Mempool) Swap(i, j int) {
	mp[i], mp[j] = mp[j], mp[i]
	mp[i].index = i
	mp[j].index = j
}

// Less needed for the sort.Interface, used to compare the structure to ensure the highest Fee items are at the top.
func (mp Mempool) Less(i, j int) bool {
	return mp[i].Fee > mp[j].Fee
}

// Push needed for the heap.Interface, allows items to be added to mempool.
func (mp *Mempool) Push(x interface{}) {
	n := len(*mp)
	txn := x.(*Transaction)
	txn.index = n
	*mp = append(*mp, txn)
}

// Pop needed for the heap.Interface, allows items to be removed to mempool.
func (mp *Mempool) Pop() interface{} {
	old := *mp
	n := len(old)
	txn := old[n-1]
	old[n-1] = nil
	txn.index = -1
	*mp = old[0 : n-1]
	return txn
}

// reorder resets the Fee listings in the heap
func (mp *Mempool) reorder(txn *Transaction, fee float64) {
	txn.Fee = fee
}
