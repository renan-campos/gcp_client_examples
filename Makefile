all: bin/list-constraints bin/find-constraint bin/has-constraint

bin/list-constraints: cmd/list-constraints/main.go
	go build -o bin/list-constraints cmd/list-constraints/main.go

bin/find-constraint: cmd/find-constraint/main.go
	go build -o bin/find-constraint cmd/find-constraint/main.go

bin/has-constraint: cmd/has-constraint/main.go
	go build -o bin/has-constraint cmd/has-constraint/main.go

clean:
	rm -f bin/*

.PHONY: clean
