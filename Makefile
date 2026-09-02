BINARY=roguelike

all: build

build:
	go build -o $(BINARY) ./cmd/roguelike

run: build
	./$(BINARY)

test:
	go test -v ./...

clean:
	rm -f $(BINARY)
