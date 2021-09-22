BINARY=plutus

all: build

build:
	go build -o $(BINARY)
docker:
	docker build . -t $(BINARY):dev

.PHONY: clean
clean:
	rm -rf $(BINARY)