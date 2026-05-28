// Package main is the temporal lesson `l1_determinism_and_replay` homework scaffold for Vibe Learn.
//
// Задача: почини нарушения детерминизма (workflow.Now/SideEffect/сортировка ключей) + GetVersion + replay-тест.
// Реализуй workflow и активности ниже — сигнатуры и тестовая поверхность
// фиксированы; CI (.github/workflows/ci.yml) гоняет `go vet` и `go test ./...`.
// Подробности и критерии приёмки — в README.md.
//
// SDK: go.temporal.io/sdk (worker + workflow + activity).
// Воркер подключается к Temporal по TEMPORAL_ADDRESS (дефолт localhost:7233 —
// совпадает с docker-compose.yml) и слушает task queue из TaskQueue().
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// ----- config -----

// envOr returns the env var for `key` if set, else `fallback`.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// TemporalAddress — адрес Temporal frontend. Дефолт совпадает с docker-compose.yml.
func TemporalAddress() string {
	return envOr("TEMPORAL_ADDRESS", "localhost:7233")
}

// TaskQueue — очередь задач, которую слушает воркер этого урока.
func TaskQueue() string {
	return envOr("TEMPORAL_TASK_QUEUE", "lesson-l1_determinism_and_replay-tq")
}

// ----- Workflow: OrderWorkflow -----
//
// Оркеструет активности ниже. Тело — TODO: добавь ExecuteActivity-шаги,
// ActivityOptions (StartToCloseTimeout, RetryPolicy) и обработку ошибок
// согласно README.md. Должно оставаться ДЕТЕРМИНИРОВАННЫМ (никаких
// time.Now/rand/итераций по map — используй workflow.Now/SideEffect).
func OrderWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("OrderWorkflow started", "taskQueue", TaskQueue())

	// TODO #1: вызови активность ProcessOrder через workflow.ExecuteActivity.
	// var processorderRes string
	// if err := workflow.ExecuteActivity(ctx, ProcessOrder).Get(ctx, &processorderRes); err != nil {
	// 	return err
	// }
	// TODO #2: вызови активность AuditOrder через workflow.ExecuteActivity.
	// var auditorderRes string
	// if err := workflow.ExecuteActivity(ctx, AuditOrder).Get(ctx, &auditorderRes); err != nil {
	// 	return err
	// }

	return nil
}

// ----- Activity #1: ProcessOrder -----
//
// детерминированная активность вместо time.Now/rand в коде workflow
func ProcessOrder(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("ProcessOrder: not implemented")
}

// ----- Activity #2: AuditOrder -----
//
// новый шаг, добавляемый через workflow.GetVersion без слома старых историй
func AuditOrder(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("AuditOrder: not implemented")
}

// ----- main entry: register worker + run with graceful shutdown -----

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Vibe Learn — temporal lesson %s scaffold up", "l1_determinism_and_replay")
	log.Printf("temporal address: %s  task queue: %s", TemporalAddress(), TaskQueue())
	log.Printf("Реализуй workflow и активности, затем `go test ./...`. README.md содержит задачу.")

	c, err := client.Dial(client.Options{HostPort: TemporalAddress()})
	if err != nil {
		log.Fatalf("unable to create Temporal client (is `docker compose up -d` running?): %v", err)
	}
	defer c.Close()

	w := worker.New(c, TaskQueue(), worker.Options{})
	w.RegisterWorkflow(OrderWorkflow)
	w.RegisterActivity(ProcessOrder)
	w.RegisterActivity(AuditOrder)

	// Graceful shutdown so `go run .` is interactive — worker.InterruptCh()
	// stops the worker on Ctrl-C / SIGTERM.
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("worker stopped with error: %v", err)
	}
}
