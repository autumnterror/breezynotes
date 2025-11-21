
test-method:
	go test -run $(METHOD) ./... -v
test-test:
	go test ./test -v
test-views:
	go test ./views -v
test-textblock:
	go test ./pkg/pkgs/default/textblock -v
test-textblock-bm:
	go test ./pkg/pkgs/default/textblock/benchmark -v
test-blocks:
	go test ./internal/blocknote/mongo/blocks -v
test-notes:
	go test ./internal/blocknote/mongo/notes -v
test-tags:
	go test ./internal/blocknote/mongo/tags -v
test-jwt:
	go test ./internal/auth/jwt -v
test-psql:
	go test ./internal/auth/psql -v

test-all:
	go test ./internal/blocknote/mongo/blocks ./internal/blocknote/mongo/notes ./internal/blocknote/mongo/tags  ./internal/auth/psql ./internal/auth/jwt
