PROTO_DIR = pkg/protos
BUILD_DIR = bin

build:
	@go build -o ${BUILD_DIR}/consenter ./cli/main.go

docker:
	@GOOS=linux go build -o ${BUILD_DIR}/consenter ./cli/main.go
	docker build -t consenter .

simulation: 
	@GOOS=linux go build -o ${BUILD_DIR}/consenter ./cli/main.go
	@docker-compose build
	@docker-compose up

proto:
	@protoc -I ${PROTO_DIR} ${PROTO_DIR}/*.proto --go_out=${PROTO_DIR}

deps:
	@dep ensure

test: 
	@go test ./... --cover


