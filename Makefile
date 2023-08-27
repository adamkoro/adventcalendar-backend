SOURCE_DIR = ./src

tidy:
	cd $(SOURCE_DIR) && go mod tidy

run:
	cd $(SOURCE_DIR) && go run main.go