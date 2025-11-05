test-notes:
	go test ./internal/blocknote/mongo/notes -v
test-tags:
	go test ./internal/blocknote/mongo/tags -v
test-jwt:
	go test ./internal/auth/jwt -v
test-psql:
	go test ./internal/auth/psql -v

test-all:
	go test ./internal/blocknote/mongo/notes ./internal/blocknote/mongo/tags  ./internal/auth/psql ./internal/auth/jwt