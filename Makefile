BINARY_NAME=klector
 
all: build test
 
build:
	go build -o ${BINARY_NAME}
 
test:
	go test -v ./...
 
run:
	go build -o ${BINARY_NAME}
	./${BINARY_NAME} run
 
clean:
	go clean
	rm ${BINARY_NAME}
