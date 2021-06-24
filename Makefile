GO=go

run:
	$(GO) run ./cmd/parser/main.go -intnx ./data/transactions.txt -limit 5000 -outtnx ./data/prioritized-transactions.txt
.PHONY: run

clean:
	rm -f ./data/prioritized-transactions.txt
.PHONY: clean
