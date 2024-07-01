package main

import (
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/agenthandlers"
)

func main() {
	var metrics agenthandlers.MetricsStats
	for {
		go agenthandlers.CollectMetrics(&metrics)
		go agenthandlers.PushMetrics("http://localhost", "update", &metrics)
	}
}
