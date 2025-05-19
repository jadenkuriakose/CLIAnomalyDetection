# Anomaly Detection System

This is a Go-based anomaly detection system designed to monitor system metrics (e.g., CPU usage) and detect anomalies using a z-score-based machine learning algorithm. The system includes a REST API for external interaction and logs metrics and anomalies to a file for monitoring.

## Features
- **Z-Score Anomaly Detection**: Identifies anomalies in system metrics based on a sliding window of data, flagging values that deviate significantly (z-score > 2.5) from the mean.
- **REST API**: Provides endpoints to submit metrics and retrieve current metric data.
- **Logging**: Records metrics, z-scores, and anomaly alerts to a log file (`anomaly_detection.log`).
- **Thread-Safe Design**: Uses mutexes to ensure safe concurrent access to the detector.
- **Extensible**: Modular design allows easy integration of real system metrics or additional detection algorithms.

## Project Structure
- `main.go`: Entry point, simulates metrics, and starts the API server.
- `detector.go`: Implements the z-score-based anomaly detection logic.
- `api.go`: Defines REST API endpoints using `gorilla/mux`.
- `logger.go`: Handles thread-safe logging to a file.

## Prerequisites
- **Go**: Version 1.16 or later (download from [golang.org](https://golang.org/dl/)).
- **Dependencies**: `github.com/gorilla/mux` for the API.
- **Optional**: `github.com/shirou/gopsutil/v3` for real system metrics (e.g., CPU usage).

## Setup
1. **Create a Project Directory**:
   ```bash
   mkdir anomaly_detection
   cd anomaly_detection
   ```

2. **Save the Source Files**:
   - Copy `main.go`, `detector.go`, `api.go`, and `logger.go` into the `anomaly_detection` directory.

3. **Initialize a Go Module**:
   ```bash
   go mod init anomaly_detection
   ```

4. **Install Dependencies**:
   ```bash
   go get github.com/gorilla/mux
   ```

5. **Run the Program**:
   ```bash
   go run *.go
   ```
   - The system starts monitoring simulated metrics and runs an API server on `http://localhost:8080`.
   - Logs are written to `anomaly_detection.log`.

## Usage
### Monitoring Metrics
- The system logs metrics every 2 seconds, including the metric value, z-score, and anomaly status.
- Check logs in `anomaly_detection.log`:
  ```bash
  cat anomaly_detection.log
  ```
  Example output:
  ```
  INFO: 2025/05/19 17:08:58 Metric: 68.00, Z-Score: 1.41, Anomaly: false
  INFO: 2025/05/19 17:09:00 ALERT: Anomaly detected with metric 150.00 (Z-Score: 5.12)
  ```

### API Endpoints
- **Submit a Metric** (`POST /api/metric`):
  ```bash
  curl -X POST -d '{"value": 100}' http://localhost:8080/api/metric
  ```
  Response:
  ```json
  {"metric":100,"isAnomaly":true,"zScore":3.5}
  ```

- **Get Current Metrics** (`GET /api/metrics`):
  ```bash
  curl http://localhost:8080/api/metrics
  ```
  Response:
  ```json
  {"metrics":[60,62,64,66,68],"mean":64,"stdDev":3.16}
  ```

### Testing Anomaly Detection
- **Simulate Anomalies**:
  - Modify `simulateSystemMetric` in `main.go` to test extreme values:
    ```go
    func simulateSystemMetric() float64 {
        if time.Now().Second()%20 == 0 {
            return 150.0 // Triggers anomaly
        }
        return 50.0 + float64(time.Now().Second()%40) // Normal: 50–90
    }
    ```
  - Rerun: `go run *.go`
  - Check logs for anomaly alerts.

- **Use Real Metrics**:
  - Install `gopsutil`:
    ```bash
    go get github.com/shirou/gopsutil/v3
    ```
  - Update `simulateSystemMetric` in `main.go`:
    ```go
    import "github.com/shirou/gopsutil/v3/cpu"

    func simulateSystemMetric() float64 {
        percent, err := cpu.Percent(time.Second, false)
        if err != nil {
            logMessage("Error getting CPU usage: " + err.Error())
            return 0
        }
        if len(percent) > 0 {
            return percent[0]
        }
        return 0
    }
    ```
  - Rerun and monitor actual CPU usage.

### Configuration
- Adjust detection parameters in `main.go`:
  ```go
  config := detectorConfig{
      thresholdZScore: 2.5, // Z-score threshold for anomalies
      windowSize:      10,  // Sliding window size
      minObservations: 5,   // Minimum data points before detection
  }
  ```
  - Lower `thresholdZScore` (e.g., to 1.5) for more sensitive detection.

## Troubleshooting
- **No Logs**: Ensure write permissions in the directory:
  ```bash
  chmod u+w .
  ```
- **API Not Responding**: Verify port `:8080` is free. Change to another port (e.g., `:8081`) in `api.go` if needed.
- **No Anomalies**: Test with extreme values via the API or lower `thresholdZScore`.
- **Dependency Errors**: Run `go mod tidy` and ensure `gorilla/mux` is installed.