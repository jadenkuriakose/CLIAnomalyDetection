package main

import (
	"fmt"
	"time"
)

func main() {
	// Initialize logger
	initLogger("anomalyDetection.log")

	// Initialize anomaly detector
	config := detectorConfig{
		thresholdZScore: 2.5,
		windowSize:      10,
		minObservations: 5,
	}
	detector := newAnomalyDetector(config)

	// Start API server in a goroutine
	go startApiServer(":8080", detector)

	// Simulate system monitoring
	for {
		// Simulate system metrics (e.g., CPU usage)
		metric := simulateSystemMetric()
		isAnomaly, zScore := detector.checkAnomaly(metric)

		logMessage(fmt.Sprintf("Metric: %.2f, Z-Score: %.2f, Anomaly: %v", metric, zScore, isAnomaly))
		if isAnomaly {
			logMessage(fmt.Sprintf("ALERT: Anomaly detected with metric %.2f (Z-Score: %.2f)", metric, zScore))
		}

		time.Sleep(2 * time.Second)
	}
}

// simulateSystemMetric generates synthetic system metrics for testing
func simulateSystemMetric() float64 {
	// Simulate CPU usage (replace with real metrics in production)
	return 50.0 + float64(time.Now().Second()%40) // Varies between 50-90
}