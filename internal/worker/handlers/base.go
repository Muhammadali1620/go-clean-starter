package handlers

import (
	"context"

	"github.com/hibiken/asynq"
)

type WorkerTask interface {
	TaskType() string
	ProcessTask(context.Context, *asynq.Task) error
}
