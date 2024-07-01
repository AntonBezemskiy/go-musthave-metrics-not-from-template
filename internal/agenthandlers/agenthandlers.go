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
	metrics.Lock()
	metrics.PollCount += 1
	runtime.ReadMemStats(&metrics.MemStats)

	metrics.Unlock()
	time.Sleep(pollInterval * time.Second)
}

func Push(addres, typemetric, namemetric, valuemetric string) {
	resp, err := http.Post(addres+"/"+typemetric+"/"+namemetric+"/"+valuemetric, "text/plain", nil)

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
		fmt.Printf("received non-200 response status: %d", resp.StatusCode)
	}

	//return nil
}

func PushMetrics(addres string, typemetric string, metrics *MetricsStats) {
	metrics.Lock()

	Push(addres, typemetric, "alloc", strconv.FormatUint(metrics.Alloc, 10))
	Push(addres, typemetric, "buckhashsys", strconv.FormatUint(metrics.BuckHashSys, 10))
	Push(addres, typemetric, "formatunit", strconv.FormatUint(metrics.Frees, 10))
	Push(addres, typemetric, "gccpufraction", strconv.FormatFloat(metrics.GCCPUFraction, 'f', 6, 64))
	Push(addres, typemetric, "gcsys", strconv.FormatUint(metrics.GCSys, 10))
	Push(addres, typemetric, "heapalloc", strconv.FormatUint(metrics.HeapAlloc, 10))
	Push(addres, typemetric, "heapidle", strconv.FormatUint(metrics.HeapIdle, 10))
	Push(addres, typemetric, "heapinuse", strconv.FormatUint(metrics.HeapInuse, 10))
	Push(addres, typemetric, "heapobjects", strconv.FormatUint(metrics.HeapObjects, 10))
	Push(addres, typemetric, "heapreleased", strconv.FormatUint(metrics.HeapReleased, 10))
	Push(addres, typemetric, "heapsys", strconv.FormatUint(metrics.HeapSys, 10))
	Push(addres, typemetric, "lastgc", strconv.FormatUint(metrics.LastGC, 10))
	Push(addres, typemetric, "lookups", strconv.FormatUint(metrics.Lookups, 10))
	Push(addres, typemetric, "mcacheinuse", strconv.FormatUint(metrics.MCacheInuse, 10))
	Push(addres, typemetric, "mcachesys", strconv.FormatUint(metrics.MCacheSys, 10))
	Push(addres, typemetric, "mspaninuse", strconv.FormatUint(metrics.MSpanInuse, 10))
	Push(addres, typemetric, "mspansys", strconv.FormatUint(metrics.MSpanSys, 10))
	Push(addres, typemetric, "mallocs", strconv.FormatUint(metrics.Mallocs, 10))
	Push(addres, typemetric, "nextgc", strconv.FormatUint(metrics.NextGC, 10))
	Push(addres, typemetric, "numforcedgc", strconv.FormatUint(uint64(metrics.NumForcedGC), 10))
	Push(addres, typemetric, "numgc", strconv.FormatUint(uint64(metrics.NumGC), 10))
	Push(addres, typemetric, "othersys", strconv.FormatUint(metrics.OtherSys, 10))
	Push(addres, typemetric, "pausetotalns", strconv.FormatUint(metrics.PauseTotalNs, 10))
	Push(addres, typemetric, "stackinuse", strconv.FormatUint(metrics.StackInuse, 10))
	Push(addres, typemetric, "stacksys", strconv.FormatUint(metrics.StackSys, 10))
	Push(addres, typemetric, "sys", strconv.FormatUint(metrics.Sys, 10))
	Push(addres, typemetric, "totalalloc", strconv.FormatUint(metrics.TotalAlloc, 10))

	metrics.Unlock()
	time.Sleep(reportInterval * time.Second)
}
