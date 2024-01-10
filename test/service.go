package test

import (
	"context"

	"github.com/go-redis/redis"
)

type Service interface {
	Test(ctx context.Context, clietn redis.Client) (err error)
}
