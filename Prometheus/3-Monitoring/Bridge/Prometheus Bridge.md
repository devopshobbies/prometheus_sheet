he Prometheus Bridge is a tool that allows you to expose non-Prometheus metrics in a format that Prometheus can scrape. This is useful when you have custom metrics or legacy systems that don't natively support Prometheus metrics.

Below is a full example of how to use the Prometheus Bridge to expose custom metrics and scrape them using Prometheus.


---------
Here’s a **complete example** of a Prometheus Bridge in action. We’ll create a simple application that exposes metrics in JSON format, then use a Python-based bridge to convert those metrics into Prometheus format for scraping.



![[bridge.png]]

---

### **Step 1: Create a Python App with JSON Metrics**
First, let’s build a Python app that exposes metrics in JSON format (not Prometheus-native).

#### `app.py`
```python
from flask import Flask
import random
import time

app = Flask(__name__)

# Simulate metrics
request_count = 0
error_count = 0

@app.route('/')
def home():
    global request_count, error_count
    request_count += 1
    if random.random() < 0.2:  # 20% chance of error
        error_count += 1
    return "Hello, World!"

@app.route('/stats')
def stats():
    # Return metrics as JSON
    return {
        "requests_total": request_count,
        "errors_total": error_count,
        "latency_seconds": random.uniform(0.1, 1.0)
    }

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

Run the app:
```bash
python app.py
```

The app will expose:
- `http://localhost:5000/` – The main endpoint.
- `http://localhost:5000/stats` – Metrics in JSON format (e.g., `{"requests_total": 5, "errors_total": 1, "latency_seconds": 0.4}`).

---

### **Step 2: Build the Prometheus Bridge**
Create a Python script (`prometheus_bridge.py`) that acts as a bridge. It will:
1. Fetch metrics from `http://localhost:5000/stats`.
2. Convert the JSON metrics into Prometheus format.
3. Expose the converted metrics at `/metrics` for Prometheus to scrape.

#### `prometheus_bridge.py`
```python
from flask import Flask
import requests

app = Flask(__name__)

# URL of the app's JSON metrics endpoint
TARGET_URL = "http://localhost:5000/stats"

@app.route('/metrics')
def metrics():
    try:
        # Fetch JSON metrics from the app
        response = requests.get(TARGET_URL)
        data = response.json()

        # Convert JSON to Prometheus format
        prom_metrics = []
        prom_metrics.append("# HELP app_requests_total Total HTTP requests.")
        prom_metrics.append("# TYPE app_requests_total counter")
        prom_metrics.append(f"app_requests_total {data['requests_total']}")

        prom_metrics.append("# HELP app_errors_total Total HTTP errors.")
        prom_metrics.append("# TYPE app_errors_total counter")
        prom_metrics.append(f"app_errors_total {data['errors_total']}")

        prom_metrics.append("# HELP app_latency_seconds HTTP request latency in seconds.")
        prom_metrics.append("# TYPE app_latency_seconds gauge")
        prom_metrics.append(f"app_latency_seconds {data['latency_seconds']}")

        return "\n".join(prom_metrics), 200, {'Content-Type': 'text/plain'}

    except Exception as e:
        return f"Error: {str(e)}", 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
```

Run the bridge:
```bash
python prometheus_bridge.py
```

The bridge will now expose Prometheus-compatible metrics at `http://localhost:8080/metrics`.

---

### **Step 3: Configure Prometheus to Scrape the Bridge**
Update your Prometheus configuration (`prometheus.yml`) to scrape the bridge’s `/metrics` endpoint.

#### `prometheus.yml`
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus_bridge'
    static_configs:
      - targets: ['localhost:8080']  # Bridge's metrics endpoint
```

Start Prometheus:
```bash
./prometheus --config.file=prometheus.yml
```

---

### **Step 4: Test the Setup**
4. **Generate Traffic**:
   Use `curl` or a browser to hit the app’s endpoint:
   ```bash
   curl http://localhost:5000
   ```

5. **Check the Bridge’s Metrics**:
   ```bash
   curl http://localhost:8080/metrics
   ```
   Output:
   ```
   # HELP app_requests_total Total HTTP requests.
   # TYPE app_requests_total counter
   app_requests_total 10
   # HELP app_errors_total Total HTTP errors.
   # TYPE app_errors_total counter
   app_errors_total 2
   # HELP app_latency_seconds HTTP request latency in seconds.
   # TYPE app_latency_seconds gauge
   app_latency_seconds 0.45
   ```

6. **Check Prometheus Targets**:
   Open `http://localhost:9090/targets` in your browser. Ensure the `prometheus_bridge` target is healthy.

7. **Query Metrics in Prometheus**:
   Use PromQL to query metrics:
   - `app_requests_total`
   - `rate(app_requests_total[1m])`
   - `app_errors_total`

---

### **Step 5: Visualize in Grafana (Optional)**
8. Install Grafana:
   ```bash
   docker run -d -p 3000:3000 grafana/grafana
   ```

9. Add Prometheus as a data source in Grafana (URL: `http://localhost:9090`).

10. Create a dashboard with panels for `app_requests_total`, `app_errors_total`, and `app_latency_seconds`.

---

### **How It Works**
- The app exposes metrics in JSON format.
- The bridge converts JSON metrics to Prometheus format.
- Prometheus scrapes the bridge’s `/metrics` endpoint.
- You can monitor the metrics in Prometheus or Grafana.

![Architecture](https://i.imgur.com/9W7oE5L.png)

---

### **Key Takeaways**
- A Prometheus Bridge acts as a translator for non-Prometheus metrics.
- Use cases include legacy systems, third-party apps, or custom metric formats.
- Tools like custom scripts, exporters, or OpenTelemetry Collector can serve as bridges.