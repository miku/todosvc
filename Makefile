TARGETS = todosvc

.PHONY: all
all: $(TARGETS)

%: cmd/%/main.go
	go build -ldflags="-w -s" -o $@ $^

.PHONY: clean
clean:
	rm -f geomimg

