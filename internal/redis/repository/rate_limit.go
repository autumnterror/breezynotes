package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var luaFixedWindow = redis.NewScript(`
	local current = redis.call("INCR", KEYS[1])
	if current == 1 then
	  redis.call("PEXPIRE", KEYS[1], ARGV[1])
	end
	local ttl = redis.call("PTTL", KEYS[1])
	return {current, ttl}
`,
)

func (s *Client) RateLimit(ctx context.Context, key string, windowMilliseconds int64) (count int64, ttl int64, err error) {
	res, err := luaFixedWindow.Run(ctx, s.Rdb, []string{key}, windowMilliseconds).Result()
	if err != nil {
		return 0, 0, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) != 2 {
		return 0, 0, fmt.Errorf("ratelimit: unexpected lua result: %T %v", res, res)
	}

	cn, ok1 := arr[0].(int64)
	ttlms, ok2 := arr[1].(int64)
	if !ok1 || !ok2 {
		return 0, 0, fmt.Errorf("ratelimit: unexpected lua types: %T %T", arr[0], arr[1])
	}

	if ttlms < 0 {
		return cn, windowMilliseconds, nil
	}

	return cn, ttlms, nil
}
