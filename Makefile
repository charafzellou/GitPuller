IMG-NAME=gitpuller

help:
	@echo "make init          :   gets the necessary godotenv Go library"
	@echo "make build         :   builds the binary in Go"
	@echo "make run           :   runs the source code in Go"
	@echo "make docker-build  :   builds the docker image"
	@echo "make docker-run    :   runs the docker image"
	@echo "make docker-clean  :   stops and deletes both the container and image"

init:
	@cd src/ && go get github.com/joho/godotenv

build:
	@cd src/ && go build -o ../bin/${IMG-NAME} .

run:
	@cd src/ && go run .

docker-build:
	@docker build -t ${IMG-NAME} .

docker-run:
	@docker run -p 8080:8080 -v "$(pwd)":/go/src/app ${IMG-NAME}

docker-clean:
	@docker stop $$(docker ps -q --filter "ancestor=${IMG-NAME}")
	@docker rm $$(docker ps -a -q --filter "ancestor=${IMG-NAME}")
	@docker rmi ${IMG-NAME}