package client

import (
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
}

type redisTaskDistributor struct {
	client *asynq.Client
}

// @inject
func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &redisTaskDistributor{client: client}
}
