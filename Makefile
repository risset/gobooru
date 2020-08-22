GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BIN=gobooru

all: test build
build: 
	$(GOBUILD) -o $(BIN) -v

test: 
	$(GOTEST) -v ./backend

clean: 
	$(GOCLEAN)
	rm -f $(BIN)

run:
	$(GOBUILD) -o $(BIN) -v ./...
	./$(BIN)

deps:
	$(GOGET) github.com/spf13/cobra
	$(GOGET) github.com/spf13/viper
	$(GOGET) github.com/mitchellh/go-homedir
