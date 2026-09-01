package app

import (
	"log"

	"github.com/hibiken/asynq"

	"new_project/internal/worker/scheduler"
)

func RunScheduler(a *App) {
	log.Println("⏰ Starting Asynq Background Task Scheduler (Cron)...")

	schedulerServer := asynq.NewScheduler(a.RedisOpt, nil)
	scheduler.RegisterAllPeriodicTasks(schedulerServer)

	if err := schedulerServer.Run(); err != nil {
		log.Fatalf("Scheduler failed: %v", err)
	}
}
