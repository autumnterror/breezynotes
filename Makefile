.PHONY: pull-all build-all build-platform build-gateway build-auth build-blocknote build-redis \
        push-all build-migrator mig-up mig-down gen \
        compose-up compose-down compose-rest compose-up-repl compose-down-repl compose-up-db compose-down-db \
        docx

# Docker image operations
pull-all:
	docker pull zitrax78/breezynotes-gateway
	docker pull zitrax78/breezynotes-auth
	docker pull zitrax78/breezynotes-blocknote
	docker pull zitrax78/breezynotes-redis

push-all:
	docker push zitrax78/breezynotes-gateway
	docker push zitrax78/breezynotes-auth
	docker push zitrax78/breezynotes-blocknote
	docker push zitrax78/breezynotes-redis

build-all:
	docker build -t zitrax78/breezynotes-gateway --file ./build/docker/gateway/dockerfile .
	docker build -t zitrax78/breezynotes-auth --file ./build/docker/auth/dockerfile .
	docker build -t zitrax78/breezynotes-blocknote --file ./build/docker/blocknote/dockerfile .
	docker build -t zitrax78/breezynotes-redis --file ./build/docker/redis/dockerfile .

build-all-mac:
	docker build --platform linux/amd64 -t zitrax78/breezynotes-gateway --file ./build/docker/gateway/dockerfile .
	docker build --platform linux/amd64 -t zitrax78/breezynotes-auth --file ./build/docker/auth/dockerfile .
	docker build --platform linux/amd64 -t zitrax78/breezynotes-blocknote --file ./build/docker/blocknote/dockerfile .
	docker build --platform linux/amd64 -t zitrax78/breezynotes-redis --file ./build/docker/redis/dockerfile .

build-gateway:
	docker build -t zitrax78/breezynotes-gateway --file ./build/docker/gateway/dockerfile .

build-auth:
	docker build -t zitrax78/breezynotes-auth --file ./build/docker/auth/dockerfile .

build-blocknote:
	docker build -t zitrax78/breezynotes-blocknote --file ./build/docker/blocknote/dockerfile .

build-redis:
	docker build -t zitrax78/breezynotes-redis --file ./build/docker/redis/dockerfile .

# Migration operations
build-migrator:
	docker build -t zitrax78/breezynotes-migrator --file ./build/docker/migrator/dockerfile .
#go build -o ./build/breezynotes/bin/migrator.exe ./cmd/migrator
mig-up: build-migrator
	./build/breezynotes/bin/migrator.exe --type up --path ./build/breezynotes/migrations

mig-down: build-migrator
	./build/breezynotes/bin/migrator.exe --type down --path ./build/breezynotes/migrations

build-migrator-mac:
	go build -o ./build/breezynotes/bin/migrator ./cmd/migrator

mig-up-mac: build-migrator-mac
	./build/breezynotes/bin/migrator --type up --path ./build/breezynotes/migrations

mig-down-mac: build-migrator-mac
	./build/breezynotes/bin/migrator --type down --path ./build/breezynotes/migrations

# Protobuf generation
gen:
	protoc --proto_path api/proto/ -I proto \
		auth.proto \
		notes.proto \
		redis.proto \
		domain.proto \
		--go_out=./api/proto/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=./api/proto/gen \
		--go-grpc_opt=paths=source_relative

# Docker Compose operations
compose-up:
	docker compose -f ./build/breezynotes/docker-compose.yml up -d

compose-down:
	docker compose -f ./build/breezynotes/docker-compose.yml down

compose-rest:
	docker compose -f ./build/breezynotes/docker-compose.yml down
	docker compose -f ./build/breezynotes/docker-compose.yml up -d

compose-up-repl:
	docker compose -f ./build/breezynotes/docker-compose.repl.yml up -d

compose-down-repl:
	docker compose -f ./build/breezynotes/docker-compose.repl.yml down

compose-up-db:
	docker compose -f ./build/breezynotes/docker-compose.db.yml up -d

compose-down-db:
	docker compose -f ./build/breezynotes/docker-compose.db.yml down
compose-down-db-v:
	docker compose -f ./build/breezynotes/docker-compose.db.yml down -v

docx:
	swag init --dir ./cmd/gateway,./internal/gateway/net/,./internal/gateway/domain/,./api/proto/gen --output ./docs

test-method:
	go test -run $(METHOD) ./... -v
test-redis:
	go test ./internal/redis/redis -v
test-test:
	go test ./test -v
test-views:
	go test ./views -v
test-textblock:
	go test ./pkg/pkgs/default/textblock -v
test-textblock-bm:
	go test ./pkg/pkgs/default/textblock/benchmark -v
test-blocks:
	go test ./internal/blocknote/repository/blocks -v
test-notes:
	go test ./internal/blocknote/repository/notes -v
test-tags:
	go test ./internal/blocknote/repository/tags -v
test-jwt:
	go test ./internal/auth/jwt -v
test-psql:
	go test ./internal/auth/repository -v

test-all:
	go test \
	    ./internal/blocknote/repository/blocks \
		./internal/blocknote/repository/notes \
		./internal/blocknote/repository/tags  \
		./internal/auth/repository \
		./internal/auth/jwt \
		./internal/redis/repository \
		./pkg/pkgs/default/textblock \
		./test \
		./views
