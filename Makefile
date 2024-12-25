CMD_DIR = cmd
BIN_DIR = bin

.PHONY: build
build: clean
	mkdir -p $(BIN_DIR)
	mkdir -p $(BIN_DIR)/compiler
	mkdir -p $(BIN_DIR)/runtime
	go build -o $(BIN_DIR)/compiler/compiler $(CMD_DIR)/compiler/main.go
	go build -o $(BIN_DIR)/runtime/runtime $(CMD_DIR)/runtime/main.go

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/compiler
	rm -rf $(BIN_DIR)/runtime
	rm -rf $(BIN_DIR)