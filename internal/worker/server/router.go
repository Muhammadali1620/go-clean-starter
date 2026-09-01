package server

import (
	"log"
	"reflect"

	"new_project/internal/core/container"
	"new_project/internal/worker/handlers"

	"github.com/hibiken/asynq"
)

func RegisterAllHandlers(mux *asynq.ServeMux, workerHandlers container.WorkerHandlers) {
	v := reflect.ValueOf(workerHandlers)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.IsNil() {
			continue
		}

		if taskHandler, ok := field.Interface().(handlers.WorkerTask); ok {
			mux.Handle(taskHandler.TaskType(), taskHandler)
			log.Printf("Automatically registered worker handler for task: %s", taskHandler.TaskType())
		} else {
			log.Printf("Warning: Field %s does not implement WorkerTask interface", v.Type().Field(i).Name)
		}
	}
}
