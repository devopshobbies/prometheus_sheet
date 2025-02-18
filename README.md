# Document for Learning Prometheus Monitoring

This document provides a structured guide to learning **Prometheus**, a powerful open-source monitoring and alerting toolkit, based on a video lecture series. The goal is to help you understand the core concepts, set up Prometheus, and integrate it into your systems for effective monitoring.

---

## Table of Contents

1. **Introduction to Prometheus**
   - What is Prometheus?
   - Key Features
   - Use Cases

2. **Core Concepts**
   - Metrics and Time Series Data
   - Data Model
   - PromQL (Prometheus Query Language)
   - Exporters and Instrumentation

3. **Setting Up Prometheus**
   - Installation
   - Configuration
   - Running Prometheus

4. **Monitoring with Prometheus**
   - Setting Up Targets
   - Using Exporters
   - Visualizing Data with Grafana

5. **Alerting with Prometheus**
   - Alerting Rules
   - Alertmanager Configuration
   - Notification Integrations

6. **Advanced Topics**
   - Scaling Prometheus
   - Long-Term Storage
   - High Availability

7. **Hands-On Exercises**
   - Step-by-Step Labs
   - Sample Projects

8. **Resources**
   - Video Lecture Links
   - Official Documentation
   - Additional Reading

---

## 1. Introduction to Prometheus

### What is Prometheus?
Prometheus is an open-source systems monitoring and alerting toolkit originally built at **SoundCloud**. It is now a standalone open-source project maintained by the **Cloud Native Computing Foundation (CNCF)**.

### Key Features
- **Multi-dimensional data model**: Metrics are identified by name and key/value pairs.
- **PromQL**: A flexible query language to leverage Prometheus' multi-dimensional data model.
- **Pull-based architecture**: Prometheus scrapes metrics from targets at regular intervals.
- **Alerting**: Supports robust alerting capabilities via Alertmanager.
- **Integration**: Works seamlessly with Grafana for visualization and other exporters for data collection.

### Use Cases
- Infrastructure monitoring
- Application performance monitoring (APM)
- Real-time alerting
- Capacity planning

---

## 2. Core Concepts

### Metrics and Time Series Data
- **Metrics**: Numeric measurements representing the state or performance of a system.
- **Time Series Data**: Metrics collected over time, stored as a sequence of data points.

### Data Model
- **Metric Name**: Identifies the type of metric (e.g., `http_requests_total`).
- **Labels**: Key/value pairs that provide additional dimensions (e.g., `method="GET"`, `status="200"`).

### PromQL (Prometheus Query Language)
- Used to query and aggregate time series data.
- Example: `http_requests_total{status="200"}`

### Exporters and Instrumentation
- **Exporters**: Tools that expose metrics in Prometheus format (e.g., Node Exporter for system metrics).
- **Instrumentation**: Adding Prometheus client libraries to your application code to expose custom metrics.

---

## 3. Setting Up Prometheus

### Installation
1. Download Prometheus from the [official website](https://prometheus.io/download/).
2. Extract the files and navigate to the directory.

### Configuration
Edit the `prometheus.yml` file to define:
- **Global settings**: Scrape interval, evaluation interval.
- **Scrape configurations**: Targets to monitor (e.g., `localhost:9090`).

### Running Prometheus
Start Prometheus using the following command:
```bash
./prometheus --config.file=prometheus.yml
```

---

## 4. Monitoring with Prometheus

### Setting Up Targets
- Add targets to the `prometheus.yml` file under `scrape_configs`.
- Example:
  ```yaml
  scrape_configs:
    - job_name: 'node'
      static_configs:
        - targets: ['localhost:9100']
  ```

### Using Exporters
- Install and run exporters (e.g., Node Exporter for system metrics).
- Configure Prometheus to scrape metrics from the exporter.

### Visualizing Data with Grafana
1. Install Grafana.
2. Add Prometheus as a data source.
3. Create dashboards to visualize metrics.

---

## 5. Alerting with Prometheus

### Alerting Rules
Define alerting rules in the `prometheus.yml` file or a separate rules file.
Example:
```yaml
groups:
  - name: example
    rules:
      - alert: HighRequestLatency
        expr: job:request_latency_seconds:mean5m > 0.5
        for: 10m
        labels:
          severity: critical
        annotations:
          summary: "High request latency detected"
```

### Alertmanager Configuration
- Install and configure Alertmanager to handle alerts.
- Define notification receivers (e.g., email, Slack).

### Notification Integrations
- Integrate with tools like Slack, PagerDuty, or email for alert notifications.

---

## 6. Advanced Topics

### Scaling Prometheus
- Use **federation** to aggregate metrics from multiple Prometheus servers.
- Implement **sharding** for large-scale deployments.

### Long-Term Storage
- Integrate with remote storage solutions like Thanos or Cortex for long-term metric retention.

### High Availability
- Run multiple Prometheus instances in parallel for redundancy.

---

## 7. Hands-On Exercises

### Step-by-Step Labs
1. Install Prometheus and Node Exporter.
2. Configure Prometheus to scrape metrics from Node Exporter.
3. Create a Grafana dashboard to visualize system metrics.
4. Set up an alert for high CPU usage.

### Sample Projects
- Monitor a web application using Prometheus and Grafana.
- Set up custom metrics using Prometheus client libraries.

---

## 8. Resources

### Video Lecture Links
- [Prometheus Basics](https://example.com/prometheus-basics)
- [Advanced Prometheus](https://example.com/advanced-prometheus)

### Official Documentation
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)

### Additional Reading
- "Prometheus: Up & Running" by Brian Brazil
- CNCF Prometheus Community Resources

---

## How to Use This Document with the GitHub Repo
1. Clone the repository:
   ```bash
   [git clone https://github.com/your-username/prometheus-learning.git](https://github.com/DannyRavi/prometheus_sheet.git)
   ```
2. Follow the instructions in the `README.md` file to set up the environment.
3. Refer to this document for explanations and context while watching the video lectures.

---

By following this guide, you will gain a solid understanding of Prometheus and be able to implement it effectively in your projects. Happy monitoring! 🚀
