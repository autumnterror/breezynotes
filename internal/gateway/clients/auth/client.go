package auth

import (
	"context"
	"os"
	"time"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/gateway/config"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	API brzrpc.AuthServiceClient
}

func New(
	cfg *config.Config,
) (*Client, error) {
	const op = "grpc.auth.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(cfg.RetriesCount)),
		grpcretry.WithPerRetryTimeout(cfg.Timeout),
		grpcretry.WithBackoff(grpcretry.BackoffExponential(cfg.Backoff)),
	}

	cc, err := grpc.NewClient(
		cfg.AddrAuth,
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, format.Error(op, err)
	}
	ctx, done := context.WithTimeout(context.Background(), 30*time.Second)
	defer done()
	log.Info(op, "start")
	cc.Connect()
	if cc.GetState() != connectivity.Ready {
		if !cc.WaitForStateChange(ctx, connectivity.Ready) {
			log.Error(op, "cant connect in 30 sec", nil)
			os.Exit(1)
		}
		log.Success(op, "CONNECT!!")
	}

	return &Client{
		API: brzrpc.NewAuthServiceClient(cc),
	}, nil
}
