SHELL = /bin/bash
TARGETS = todosvc

.PHONY: all
all: $(TARGETS) todo.db

%: cmd/%/main.go
	go build -ldflags="-w -s" -o $@ $^

.PHONY: clean
clean:
	rm -f $(TARGETS)

.PHONY: purge
purge: clean
	rm -f todo.db

todo.db:
	go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	# migrate create -ext sql -dir db/migrations -seq create_todo_table
	migrate -database sqlite3://todo.db -path db/migrations up
