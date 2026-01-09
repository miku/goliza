TARGETS = goliza

.PHONY: all
all: $(TARGETS)

goliza: goliza.go
	go build -o $@ $^

.PHONY: clean
clean:
	rm -rf $(TARGETS)
