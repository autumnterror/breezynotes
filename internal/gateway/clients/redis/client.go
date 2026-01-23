package redis

import (
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/config"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	API brzrpc.RedisServiceClient
}

func New(
	cfg *config.Config,
) (*Client, error) {
	const op = "grpc.redis.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(cfg.RetriesCount)),
		grpcretry.WithPerRetryTimeout(cfg.Timeout),
		grpcretry.WithBackoff(grpcretry.BackoffExponential(cfg.Backoff)),
	}

	cc, err := grpc.NewClient(
		cfg.AddrRedis,
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Keepalive disabled for now; re-enable if external networks start dropping idle conns.
		// grpc.WithKeepaliveParams(keepalive.ClientParameters{
		// 	Time:                5 * time.Minute,
		// 	Timeout:             11 * time.Second,
		// 	PermitWithoutStream: true,
		// }),
	)
	if err != nil {
		return nil, format.Error(op, err)
	}

	log.Info(op, "start")
	for {
		cc.Connect()
		if cc.GetState() == connectivity.Ready {
			log.Success(op, "CONNECT!!")
			break
		}
		time.Sleep(3 * time.Second)
	}

	return &Client{
		API: brzrpc.NewRedisServiceClient(cc),
	}, nil
}
