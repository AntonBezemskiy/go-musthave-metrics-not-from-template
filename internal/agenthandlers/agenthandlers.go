package agenthandlers

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	pollInterval   = 2
	reportInterval = 10
)

// MetricsStats структура для хранения метрик
type MetricsStats struct {
	sync.Mutex
	runtime.MemStats
	PollCount   int64
	RandomValue float64
}

// CollectMetrics собирает метрики
func CollectMetrics(metrics *MetricsStats) {
	metrics.Lock()
	defer metrics.Unlock()
	metrics.PollCount++
	runtime.ReadMemStats(&metrics.MemStats)
}

// CollectMetricsTimer запускает сбор метрик с интервалом
func CollectMetricsTimer(metrics *MetricsStats) {
	for {
		CollectMetrics(metrics)
		time.Sleep(pollInterval * time.Second)
	}
}

// Push отправляет метрику на сервер и возвращает ошибку при неудаче
func Push(address, action, typemetric, namemetric, valuemetric string) error {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", address, action, typemetric, namemetric, valuemetric)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("error with post: %s, %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response status: %d for url: %s", resp.StatusCode, url)
	}
	return nil
}

// PushMetrics отправляет все метрики
func PushMetrics(address, action string, metrics *MetricsStats) {
	metrics.Lock()
	defer metrics.Unlock()

	typemetricgauge := "gauge"
	metricsToSend := []struct {
		name  string
		value string
	}{
		{"alloc", strconv.FormatUint(metrics.Alloc, 10)},
		{"buckhashsys", strconv.FormatUint(metrics.BuckHashSys, 10)},
		{"formatunit", strconv.FormatUint(metrics.Frees, 10)},
		{"gccpufraction", strconv.FormatFloat(metrics.GCCPUFraction, 'f', 6, 64)},
		{"gcsys", strconv.FormatUint(metrics.GCSys, 10)},
		{"heapalloc", strconv.FormatUint(metrics.HeapAlloc, 10)},
		{"heapidle", strconv.FormatUint(metrics.HeapIdle, 10)},
		{"heapinuse", strconv.FormatUint(metrics.HeapInuse, 10)},
		{"heapobjects", strconv.FormatUint(metrics.HeapObjects, 10)},
		{"heapreleased", strconv.FormatUint(metrics.HeapReleased, 10)},
		{"heapsys", strconv.FormatUint(metrics.HeapSys, 10)},
		{"lastgc", strconv.FormatUint(metrics.LastGC, 10)},
		{"lookups", strconv.FormatUint(metrics.Lookups, 10)},
		{"mcacheinuse", strconv.FormatUint(metrics.MCacheInuse, 10)},
		{"mcachesys", strconv.FormatUint(metrics.MCacheSys, 10)},
		{"mspaninuse", strconv.FormatUint(metrics.MSpanInuse, 10)},
		{"mspansys", strconv.FormatUint(metrics.MSpanSys, 10)},
		{"mallocs", strconv.FormatUint(metrics.Mallocs, 10)},
		{"nextgc", strconv.FormatUint(metrics.NextGC, 10)},
		{"numforcedgc", strconv.FormatUint(uint64(metrics.NumForcedGC), 10)},
		{"numgc", strconv.FormatUint(uint64(metrics.NumGC), 10)},
		{"othersys", strconv.FormatUint(metrics.OtherSys, 10)},
		{"pausetotalns", strconv.FormatUint(metrics.PauseTotalNs, 10)},
		{"stackinuse", strconv.FormatUint(metrics.StackInuse, 10)},
		{"stacksys", strconv.FormatUint(metrics.StackSys, 10)},
		{"sys", strconv.FormatUint(metrics.Sys, 10)},
		{"totalalloc", strconv.FormatUint(metrics.TotalAlloc, 10)},
	}

	for _, metric := range metricsToSend {
		err := Push(address, action, typemetricgauge, metric.name, metric.value)
		if err != nil {
			fmt.Printf("Failed to push metric %s: %v\n", metric.name, err)
		}
	}
}

// PushMetricsTimer запускает отправку метрик с интервалом
func PushMetricsTimer(address, action string, metrics *MetricsStats) {
	for {
		PushMetrics(address, action, metrics)
		time.Sleep(reportInterval * time.Second)
	}
}
