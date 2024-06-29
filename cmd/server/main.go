package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func handlerOther(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func handlerUpdate(res http.ResponseWriter, req *http.Request, storage *MemStorage) {
	// Проверка на nil для storage
	if storage == nil {
		http.Error(res, "Storage not initialized", http.StatusInternalServerError)
		return
	}

	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.Header().Set("Content-Type", "text/plain")

	args := strings.Split(req.URL.Path, "/")

	if len(args) < 5 {
		//http.Error(res, "Request must contain name of metric", http.StatusNotFound)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := args[2]
	metricName := args[3]
	metricValue := args[4]
	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			//http.Error(res, "Request must contain valid value of metric", http.StatusBadRequest)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.gauges[metricName] = value
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			//http.Error(res, "Request must contain valid value of metric", http.StatusBadRequest)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.counters[metricName] += value
	default:
		//http.Error(res, "Request must contain valid type of metric", http.StatusBadRequest)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func main() {
	storage := NewMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerOther)
	mux.HandleFunc("/update/", func(res http.ResponseWriter, req *http.Request) {
		handlerUpdate(res, req, storage)
	})

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
