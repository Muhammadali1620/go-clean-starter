package app

import (
	"log"

	"github.com/hibiken/asynq"

	workerServer "new_project/internal/worker/server"
)

func RunWorker(a *App) {
	log.Println("🚀 Starting Asynq Background Task Worker...")

	srv := asynq.NewServer(
		a.RedisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	workerServer.RegisterAllHandlers(mux, a.Container.WorkerHandlers)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Worker server failed: %v", err)
	}
}
