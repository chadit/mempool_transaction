package main

import (
	"container/heap"
	"flag"
	"fmt"
	"os"

	"github.com/chadit/mempool_transaction/internal/mempool"
	"github.com/pkg/errors"
)

func main() {

	var pathToTnxFile = flag.String("intnx", "../../data/transactions.txt", "path to transaction file to parse")
	var poolLimit = flag.Int("limit", 5000, "upper limit for number of transactions to return")
	var tnxOutputFile = flag.String("outtnx", "../../data/prioritized-transactions.txt", "path to transaction file to parse")

	flag.Parse()

	mp := mempool.New()
	if err := mp.Parse(*pathToTnxFile, *poolLimit); err != nil {
		// Depending on the nature of the application, a panic hear might not be needed.
		// Panicing for now, for the sake of demo.
		panic(err)
	}

	var lines []string
	for mp.Len() > 0 {
		item := heap.Pop(mp).(*mempool.Transaction)
		lines = append(lines, fmt.Sprintf("%s %s %s %s\n", item.Hash, item.Gas, item.Fpg, item.Signature))
	}

	if len(lines) > 0 {
		if err := output(*tnxOutputFile, lines); err != nil {
			// Depending on the nature of the application, a panic hear might not be needed.
			// Panicing for now, for the sake of demo.
			panic(err)
		}
	}

}

func output(tnxOutputFile string, lines []string) error {
	f, err := os.Create(tnxOutputFile)
	if err != nil {
		return errors.Wrap(err, "create output file")
	}

	defer f.Close()

	for _, line := range lines {
		_, err := f.WriteString(line)
		if err != nil {
			return errors.Wrap(err, "writing to output file")
		}
	}

	return nil
}
