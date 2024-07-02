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

type MetricsStats struct {
	sync.Mutex
	runtime.MemStats
	PollCount   int64
	RandomValue float64
}

func CollectMetrics(metrics *MetricsStats) {
	for {
		metrics.Lock()
		metrics.PollCount += 1
		runtime.ReadMemStats(&metrics.MemStats)

		metrics.Unlock()

		//fmt.Println("Collect metrics")
		time.Sleep(pollInterval * time.Second)
	}
}

// При неудачной отправке сообщения можно добавить возвращение ошибки из функции
func Push(addres, action, typemetric, namemetric, valuemetric string) {
	url := addres + "/" + action + "/" + typemetric + "/" + namemetric + "/" + valuemetric
	resp, err := http.Post(url, "text/plain", nil)

	if err != nil {
		fmt.Printf("Error with post: %s \n", url)
		fmt.Println(err)
		fmt.Print('\n')
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error with post: %s \n", url)
		fmt.Printf("received non-200 response status: %d \n\n", resp.StatusCode)
	}
}

func PushMetrics(addres, action string, metrics *MetricsStats) {
	for {
		metrics.Lock()
		typemetricgauge := "gauge"
		Push(addres, action, typemetricgauge, "alloc", strconv.FormatUint(metrics.Alloc, 10))
		Push(addres, action, typemetricgauge, "buckhashsys", strconv.FormatUint(metrics.BuckHashSys, 10))
		Push(addres, action, typemetricgauge, "formatunit", strconv.FormatUint(metrics.Frees, 10))
		Push(addres, action, typemetricgauge, "gccpufraction", strconv.FormatFloat(metrics.GCCPUFraction, 'f', 6, 64))
		Push(addres, action, typemetricgauge, "gcsys", strconv.FormatUint(metrics.GCSys, 10))
		Push(addres, action, typemetricgauge, "heapalloc", strconv.FormatUint(metrics.HeapAlloc, 10))
		Push(addres, action, typemetricgauge, "heapidle", strconv.FormatUint(metrics.HeapIdle, 10))
		Push(addres, action, typemetricgauge, "heapinuse", strconv.FormatUint(metrics.HeapInuse, 10))
		Push(addres, action, typemetricgauge, "heapobjects", strconv.FormatUint(metrics.HeapObjects, 10))
		Push(addres, action, typemetricgauge, "heapreleased", strconv.FormatUint(metrics.HeapReleased, 10))
		Push(addres, action, typemetricgauge, "heapsys", strconv.FormatUint(metrics.HeapSys, 10))
		Push(addres, action, typemetricgauge, "lastgc", strconv.FormatUint(metrics.LastGC, 10))
		Push(addres, action, typemetricgauge, "lookups", strconv.FormatUint(metrics.Lookups, 10))
		Push(addres, action, typemetricgauge, "mcacheinuse", strconv.FormatUint(metrics.MCacheInuse, 10))
		Push(addres, action, typemetricgauge, "mcachesys", strconv.FormatUint(metrics.MCacheSys, 10))
		Push(addres, action, typemetricgauge, "mspaninuse", strconv.FormatUint(metrics.MSpanInuse, 10))
		Push(addres, action, typemetricgauge, "mspansys", strconv.FormatUint(metrics.MSpanSys, 10))
		Push(addres, action, typemetricgauge, "mallocs", strconv.FormatUint(metrics.Mallocs, 10))
		Push(addres, action, typemetricgauge, "nextgc", strconv.FormatUint(metrics.NextGC, 10))
		Push(addres, action, typemetricgauge, "numforcedgc", strconv.FormatUint(uint64(metrics.NumForcedGC), 10))
		Push(addres, action, typemetricgauge, "numgc", strconv.FormatUint(uint64(metrics.NumGC), 10))
		Push(addres, action, typemetricgauge, "othersys", strconv.FormatUint(metrics.OtherSys, 10))
		Push(addres, action, typemetricgauge, "pausetotalns", strconv.FormatUint(metrics.PauseTotalNs, 10))
		Push(addres, action, typemetricgauge, "stackinuse", strconv.FormatUint(metrics.StackInuse, 10))
		Push(addres, action, typemetricgauge, "stacksys", strconv.FormatUint(metrics.StackSys, 10))
		Push(addres, action, typemetricgauge, "sys", strconv.FormatUint(metrics.Sys, 10))
		Push(addres, action, typemetricgauge, "totalalloc", strconv.FormatUint(metrics.TotalAlloc, 10))

		//fmt.Println("Push metrics")

		metrics.Unlock()
		time.Sleep(reportInterval * time.Second)
	}
}
