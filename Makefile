all: bin/list-constraints

bin/list-constraints: cmd/list-constraints/main.go
	go build -o bin/list-constraints cmd/list-constraints/main.go

clean:
	rm -f bin/list-constraints

.PHONY: clean
