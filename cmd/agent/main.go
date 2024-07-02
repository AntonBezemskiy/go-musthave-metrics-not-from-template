package main

import (
	"time"

	"github.com/AntonBezemskiy/go-musthave-metrics/internal/agenthandlers"
)

func main() {
	var metrics agenthandlers.MetricsStats
	go agenthandlers.CollectMetrics(&metrics)
	time.Sleep(50 * time.Millisecond)
	go agenthandlers.PushMetrics("http://localhost:8080", "update", &metrics)

	// блокировка main, чтобы функции бесконечно выполнялись
	select {}
}
