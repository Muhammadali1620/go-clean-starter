# 🧩 Modularity Guide: How to Remove or Keep Modules

This project is designed as a **Universal Plug-and-Play Template**.
You can use it with all features (HTTP + Telegram Bot + Redis + Workers) or strip it down to **Pure HTTP + PostgreSQL** in 1 minute.

---

## 1. How to run in different modes (No code changes needed)

You don't need to delete files if you just want to run specific parts:

* **Only HTTP API:** `go run cmd/app/main.go -mode=http`
* **Only Telegram Bot:** `go run cmd/app/main.go -mode=bot`
* **Local Development (HTTP + Bot together in 1 process):** `go run cmd/app/main.go -mode=hybrid`
* **Background Worker:** `go run cmd/app/main.go -mode=task-worker`
* **Cron Scheduler:** `go run cmd/app/main.go -mode=scheduler`

---

## 2. How to COMPLETELY remove Telegram Bot for a new project

If your new project does not need Telegram at all:

1. Delete the folder: `internal/bot/`
2. Delete the config file: `internal/core/config/bot.go`
3. In `internal/core/config/config.go`, remove the `Bot BotConfig` field.
4. Regenerate the DI container:
   ```bash
   go run cmd/tools/digen/main.go
   ```

✅ Done! Telegram Bot is 100% removed and will not be included in the compiled binary.

## 3. How to COMPLETELY remove Redis & Background Workers

If your new project only needs simple HTTP + PostgreSQL (no cache, no queues):

1. Delete folder: `internal/worker/`
2. Delete files:
   `internal/repositories/cache_repository.go`
   `internal/core/config/redis.go`
   `internal/core/database/redis.go`
3. In `internal/core/config/config.go`, remove the `Redis RedisConfig` field.
4. Regenerate DI:
    ```bash
    go run cmd/tools/digen/main.go
    ```

✅ Done! Redis and Asynq are completely removed.

## 4. How to COMPLETELY remove Firebase Cloud Messaging (FCM)

If your project does not need Push Notifications:

1. Delete folder: `internal/providers/firebase/`
2. Regenerate DI:
   ```bash
   go run cmd/tools/digen/main.go
   ```

✅ Done! Firebase SDK will not be included in the container.

## 5. How to add a new Feature (Example: Product API)

1. Create model in `internal/models/product_model.go`
2. Create repository in `internal/repositories/product_repository.go` (add `// @inject` above constructor)
3. Create service in `internal/services/product_service.go` (add `// @inject` above constructor)
4. Create handler in `internal/handlers/product_handler.go` (add `// @inject` above constructor)
5. Run code generator:
   ```bash
   go run cmd/tools/digen/main.go
   ```

6. Add route in `internal/core/server/routes.go`:
   ``` go
    v1.GET("/products", c.Handlers.Product.GetAll)
   ```
---
