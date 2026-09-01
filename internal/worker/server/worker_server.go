package server

import (
	"github.com/hibiken/asynq"
)

type WorkerServer struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

// func NewWorkerServer(redisOpt asynq.RedisClientOpt, pushHandler *handlers.PushTaskHandler) *WorkerServer {
// 	srv := asynq.NewServer(
// 		redisOpt,
// 		asynq.Config{
// 			Concurrency: 10,
// 			Queues: map[string]int{
// 				"critical": 6,
// 				"default":  3,
// 				"low":      1,
// 			},
// 		},
// 	)

// 	mux := asynq.NewServeMux()

// 	mux.Handle(tasks.TypeSendPush, pushHandler)

// 	return &WorkerServer{
// 		server: srv,
// 		mux:    mux,
// 	}
// }

// func (ws *WorkerServer) Start() error {
// 	log.Println("Worker Server is starting and listening to Redis queues...")
// 	return ws.server.Run(ws.mux)
// }
