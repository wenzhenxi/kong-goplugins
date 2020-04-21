so-build: go-hello.so go-log.so go-exit.so go-token.so

# Docker
start:
	docker-compose -f ./docker-compose.yml up -d

stop:
	docker-compose -f ./docker-compose.yml down

build:
	 docker build -t kong-goplugins .


# So
go-hello.so:
	go build -o go-so/go-hello.so -buildmode=plugin ./app/go_hello.go

go-log.so:
	go build -o go-so/go-log.so -buildmode=plugin ./app/go_log.go

go-exit.so:
	go build -o go-so/go-exit.so -buildmode=plugin ./app/go_exit.go

go-token.so:
	go build -o go-so/go-token.so -buildmode=plugin ./app/go_token.go
