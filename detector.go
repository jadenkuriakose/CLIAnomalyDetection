package main

import (
	"math"
	"sync"
)

type detectorConfig struct {
	thresholdZScore float64 // Z-score threshold for anomaly detection
	windowSize      int     // Number of data points in sliding window
	minObservations int     // Minimum observations before anomaly detection
}

type anomalyDetector struct {
	config        detectorConfig
	metrics       []float64
	mutex         sync.Mutex
	mean          float64
	stdDev        float64
	observations  int
}

// newAnomalyDetector creates a new detector with the given configuration
func newAnomalyDetector(config detectorConfig) *anomalyDetector {
	return &anomalyDetector{
		config:  config,
		metrics: make([]float64, 0, config.windowSize),
	}
}

// checkAnomaly evaluates if the metric is an anomaly based on z-score
func (ad *anomalyDetector) checkAnomaly(metric float64) (bool, float64) {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()

	// Add metric to sliding window
	ad.metrics = append(ad.metrics, metric)
	ad.observations++
	if len(ad.metrics) > ad.config.windowSize {
		ad.metrics = ad.metrics[1:]
	}

	// Update statistics
	ad.updateStatistics()

	// Require minimum observations before detecting anomalies
	if ad.observations < ad.config.minObservations {
		return false, 0
	}

	// Calculate z-score
	zScore := (metric - ad.mean) / ad.stdDev
	if math.IsNaN(zScore) || math.IsInf(zScore, 0) {
		return false, 0
	}

	return math.Abs(zScore) > ad.config.thresholdZScore, zScore
}

// updateStatistics calculates mean and standard deviation
func (ad *anomalyDetector) updateStatistics() {
	if len(ad.metrics) == 0 {
		ad.mean = 0
		ad.stdDev = 0
		return
	}

	// Calculate mean
	sum := 0.0
	for _, val := range ad.metrics {
		sum += val
	}
	ad.mean = sum / float64(len(ad.metrics))

	// Calculate standard deviation
	sumSquaredDiff := 0.0
	for _, val := range ad.metrics {
		diff := val - ad.mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(len(ad.metrics))
	ad.stdDev = math.Sqrt(variance)
}

// getMetrics returns a copy of the current metrics window
func (ad *anomalyDetector) getMetrics() []float64 {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()
	return append([]float64{}, ad.metrics...)
}