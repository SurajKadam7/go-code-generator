package test

import (
	"context"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	H2H(ctx context.Context, clietn redis.Client) (err error)
	DealEnd(ctx context.Context, userId int64, clietn redis.Client) (err error)
	Start(ctx context.Context, tourId int64, clietn mongo.Client) (err error)
}
