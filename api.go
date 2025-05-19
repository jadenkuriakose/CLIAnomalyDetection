package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type anomalyResponse struct {
	Metric    float64 `json:"metric"`
	IsAnomaly bool    `json:"isAnomaly"`
	ZScore    float64 `json:"zScore"`
}

type metricsResponse struct {
	Metrics []float64 `json:"metrics"`
	Mean    float64   `json:"mean"`
	StdDev  float64   `json:"stdDev"`
}

func startApiServer(addr string, detector *anomalyDetector) {
	router := mux.NewRouter()

	// Endpoint to submit a new metric
	router.HandleFunc("/api/metric", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Value float64 `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logMessage("API Error: Invalid request body")
			return
		}

		isAnomaly, zScore := detector.checkAnomaly(req.Value)
		response := anomalyResponse{
			Metric:    req.Value,
			IsAnomaly: isAnomaly,
			ZScore:    zScore,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			logMessage("API Error: Failed to encode response")
		}
	}).Methods("POST")

	// Endpoint to get current metrics
	router.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := detector.getMetrics()
		response := metricsResponse{
			Metrics: metrics,
			Mean:    detector.mean,
			StdDev:  detector.stdDev,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			logMessage("API Error: Failed to encode response")
		}
	}).Methods("GET")

	logMessage("Starting API server on " + addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logMessage("API Server Error: " + err.Error())
	}
}