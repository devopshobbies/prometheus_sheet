

-------

در Prometheus به **Bridge** به ابزارها یا مکانیزم‌هایی گفته می‌شود که Prometheus را با سایر سیستم‌های مانیتورینگ، ذخیره‌سازی داده‌ها یا پلتفرم‌های هشداردهی (Alerting) ادغام می‌کنند. این Bridgeها به Prometheus اجازه می‌دهند تا با سیستم‌های خارجی ارتباط برقرار کند، داده‌ها را به اشتراک بگذارد یا از قابلیت‌های دیگر سیستم‌ها استفاده کند. در ادامه، به بررسی مفهوم Bridge در Prometheus و برخی از ابزارهای رایج که به عنوان Bridge عمل می‌کنند، می‌پردازیم.
### **چرا Bridge در Prometheus مهم است؟**
‏Prometheus یک سیستم مانیتورینگ قدرتمند است، اما برخی از قابلیت‌ها را به‌صورت پیش‌فرض ارائه نمی‌دهد. برای مثال:
- **ذخیره‌سازی بلندمدت داده‌ها**: Prometheus به‌صورت پیش‌فرض داده‌ها را به‌صورت محلی و برای مدت محدودی ذخیره می‌کند.
- **یکپارچه‌سازی با سیستم‌های هشداردهی**: Prometheus از سیستم‌های هشداردهی خارجی مانند PagerDuty، Slack و غیره پشتیبانی می‌کند، اما برای استفاده از آن‌ها نیاز به Bridge دارد.
- **ارتباط با سایر سیستم‌های مانیتورینگ**: Prometheus ممکن است نیاز داشته باشد با سیستم‌های دیگر مانند Grafana، Thanos یا Cortex ادغام شود.

‏Bridgeها این شکاف‌ها را پر می‌کنند و به Prometheus اجازه می‌دهند تا با سیستم‌های خارجی ارتباط برقرار کند.

![image](src/bridge.png)

---

### **انواع Bridge در Prometheus**

#### ۱. **Bridge برای ذخیره‌سازی بلندمدت داده‌ها**
‏Prometheus به‌صورت پیش‌فرض داده‌ها را به‌صورت محلی و برای مدت محدودی ذخیره می‌کند. برای ذخیره‌سازی بلندمدت داده‌ها، می‌توان از Bridgeهایی مانند **Thanos** یا **Cortex** استفاده کرد:
- ‏**Thanos**: یک ابزار متن‌باز است که به Prometheus اجازه می‌دهد داده‌ها را در ذخیره‌سازی‌های ابری (مانند S3) ذخیره کند و امکان جستجوی بلندمدت در داده‌ها را فراهم می‌کند.
- ‏**Cortex**: یک پلتفرم ذخیره‌سازی توزیع‌شده است که به Prometheus اجازه می‌دهد داده‌ها را به‌صورت مقیاس‌پذیر و بلندمدت ذخیره کند.

#### ۲. **Bridge برای هشداردهی (Alerting)**
‏Prometheus از سیستم‌های هشداردهی خارجی پشتیبانی می‌کند، اما برای ارسال هشدارها به این سیستم‌ها نیاز به Bridge دارد. برخی از ابزارهای رایج عبارتند از:
-‏ **Alertmanager**: یک جزء داخلی Prometheus است که وظیفه‌ی مدیریت و ارسال هشدارها به سیستم‌های خارجی مانند Slack، PagerDuty، Email و غیره را بر عهده دارد.
-‏ **Webhook Bridge**: با استفاده از Webhook، می‌توان هشدارهای Prometheus را به هر سیستم خارجی ارسال کرد.

#### ۳. **Bridge برای یکپارچه‌سازی با Grafana**
‏Grafana یک ابزار قدرتمند برای تجسم داده‌ها است. Prometheus به‌صورت پیش‌فرض با Grafana ادغام می‌شود و می‌تواند داده‌ها را به Grafana ارسال کند. این ادغام به‌صورت مستقیم انجام می‌شود و نیاز به Bridge خاصی ندارد، اما می‌توان از ابزارهایی مانند **Prometheus Exporter** برای جمع‌آوری داده‌های اضافی و ارسال آن‌ها به Grafana استفاده کرد.

#### ۴. **Bridge برای ارتباط با سایر سیستم‌های مانیتورینگ**
‏Prometheus ممکن است نیاز داشته باشد با سایر سیستم‌های مانیتورینگ مانند **Graphite**, **InfluxDB** یا **OpenTSDB** ارتباط برقرار کند. برای این کار، می‌توان از Bridgeهایی مانند **Prometheus Remote Write** یا **Prometheus Exporter** استفاده کرد:
- ‏**Prometheus Remote Write**: این قابلیت به Prometheus اجازه می‌دهد داده‌ها را به سیستم‌های ذخیره‌سازی خارجی ارسال کند.
- ‏**Prometheus Exporter**: این ابزارها داده‌ها را از سیستم‌های دیگر جمع‌آوری کرده و در قالب متریک‌های Prometheus ارائه می‌دهند.

---

### **مثال‌هایی از Bridge در Prometheus**

#### ۱. **استفاده از Thanos برای ذخیره‌سازی بلندمدت**
```yaml
# پیکربندی Prometheus برای ارسال داده‌ها به Thanos
remote_write:
  - url: "http://thanos-receive:10908/api/v1/receive"
```

در این مثال، Prometheus داده‌ها را به Thanos ارسال می‌کند تا در ذخیره‌سازی‌های ابری (مانند S3) ذخیره شود.

#### ۲. **استفاده از Alertmanager برای ارسال هشدار به Slack**
```yaml
# پیکربندی Alertmanager برای ارسال هشدار به Slack
route:
  receiver: 'slack-notifications'
receivers:
  - name: 'slack-notifications'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/...'
        channel: '#alerts'
```

در این مثال، Alertmanager هشدارهای Prometheus را به یک کانال Slack ارسال می‌کند.

#### ۳. **استفاده از Prometheus Remote Write برای ارسال داده‌ها به InfluxDB**
```yaml
# پیکربندی Prometheus برای ارسال داده‌ها به InfluxDB
remote_write:
  - url: "http://influxdb:8086/api/v1/prom/write?db=prometheus"
```

در این مثال، Prometheus داده‌ها را به InfluxDB ارسال می‌کند.

---

### **جمع‌بندی**
‏Bridgeها در Prometheus نقش مهمی در ادغام این سیستم با سایر پلتفرم‌ها و ابزارها ایفا می‌کنند. این Bridgeها به Prometheus اجازه می‌دهند تا:
- داده‌ها را به‌صورت بلندمدت ذخیره کند.
- هشدارها را به سیستم‌های خارجی ارسال کند.
- با سایر سیستم‌های مانیتورینگ ارتباط برقرار کند.
- داده‌ها را به ابزارهای تجسمی مانند Grafana ارسال کند.

با استفاده از Bridgeها، می‌توان قابلیت‌های Prometheus را گسترش داد و آن را به یک سیستم مانیتورینگ کامل و یکپارچه تبدیل کرد.

----------

حتماً! در این بخش، به‌صورت عمیق‌تر و با جزئیات بیشتر به مفهوم **Bridge** در Prometheus می‌پردازیم. Bridgeها در Prometheus نقش حیاتی در ادغام این سیستم با سایر پلتفرم‌ها و ابزارها ایفا می‌کنند. این ابزارها به Prometheus اجازه می‌دهند تا با سیستم‌های خارجی ارتباط برقرار کند، داده‌ها را به اشتراک بگذارد یا از قابلیت‌های دیگر سیستم‌ها استفاده کند. در ادامه، به بررسی دقیق‌تر Bridgeها در Prometheus و نحوه‌ی استفاده از آن‌ها می‌پردازیم.

---

### **چرا Bridgeها در Prometheus مهم هستند؟**
‏Prometheus یک سیستم مانیتورینگ قدرتمند است، اما برخی از قابلیت‌ها را به‌صورت پیش‌فرض ارائه نمی‌دهد. برای مثال:
- **ذخیره‌سازی بلندمدت داده‌ها**: Prometheus به‌صورت پیش‌فرض داده‌ها را به‌صورت محلی و برای مدت محدودی ذخیره می‌کند.
- **یکپارچه‌سازی با سیستم‌های هشداردهی**: Prometheus از سیستم‌های هشداردهی خارجی مانند PagerDuty، Slack و غیره پشتیبانی می‌کند، اما برای استفاده از آن‌ها نیاز به Bridge دارد.
- **ارتباط با سایر سیستم‌های مانیتورینگ**: Prometheus ممکن است نیاز داشته باشد با سیستم‌های دیگر مانند Grafana، Thanos یا Cortex ادغام شود.

‏Bridgeها این شکاف‌ها را پر می‌کنند و به Prometheus اجازه می‌دهند تا با سیستم‌های خارجی ارتباط برقرار کند.

---

### **انواع Bridge در Prometheus**

#### ۱. **Bridge برای ذخیره‌سازی بلندمدت داده‌ها**
‏Prometheus به‌صورت پیش‌فرض داده‌ها را به‌صورت محلی و برای مدت محدودی ذخیره می‌کند. برای ذخیره‌سازی بلندمدت داده‌ها، می‌توان از Bridgeهایی مانند **Thanos** یا **Cortex** استفاده کرد:
-‏ **Thanos**: یک ابزار متن‌باز است که به Prometheus اجازه می‌دهد داده‌ها را در ذخیره‌سازی‌های ابری (مانند S3) ذخیره کند و امکان جستجوی بلندمدت در داده‌ها را فراهم می‌کند.
-‏ **Cortex**: یک پلتفرم ذخیره‌سازی توزیع‌شده است که به Prometheus اجازه می‌دهد داده‌ها را به‌صورت مقیاس‌پذیر و بلندمدت ذخیره کند.

#### ۲. **Bridge برای هشداردهی (Alerting)**
‏Prometheus از سیستم‌های هشداردهی خارجی پشتیبانی می‌کند، اما برای ارسال هشدارها به این سیستم‌ها نیاز به Bridge دارد. برخی از ابزارهای رایج عبارتند از:
-‏ **Alertmanager**: یک جزء داخلی Prometheus است که وظیفه‌ی مدیریت و ارسال هشدارها به سیستم‌های خارجی مانند Slack، PagerDuty، Email و غیره را بر عهده دارد.
-‏ **Webhook Bridge**: با استفاده از Webhook، می‌توان هشدارهای Prometheus را به هر سیستم خارجی ارسال کرد.

#### ۳. **Bridge برای یکپارچه‌سازی با Grafana**
‏Grafana یک ابزار قدرتمند برای تجسم داده‌ها است. Prometheus به‌صورت پیش‌فرض با Grafana ادغام می‌شود و می‌تواند داده‌ها را به Grafana ارسال کند. این ادغام به‌صورت مستقیم انجام می‌شود و نیاز به Bridge خاصی ندارد، اما می‌توان از ابزارهایی مانند **Prometheus Exporter** برای جمع‌آوری داده‌های اضافی و ارسال آن‌ها به Grafana استفاده کرد.

#### ۴. **Bridge برای ارتباط با سایر سیستم‌های مانیتورینگ**
‏Prometheus ممکن است نیاز داشته باشد با سایر سیستم‌های مانیتورینگ مانند **Graphite**, **InfluxDB** یا **OpenTSDB** ارتباط برقرار کند. برای این کار، می‌توان از Bridgeهایی مانند **Prometheus Remote Write** یا **Prometheus Exporter** استفاده کرد:
-‏ **Prometheus Remote Write**: این قابلیت به Prometheus اجازه می‌دهد داده‌ها را به سیستم‌های ذخیره‌سازی خارجی ارسال کند.
- ‏**Prometheus Exporter**: این ابزارها داده‌ها را از سیستم‌های دیگر جمع‌آوری کرده و در قالب متریک‌های Prometheus ارائه می‌دهند.

---

### **مثال‌هایی از Bridge در Prometheus**

#### ۱. **استفاده از Thanos برای ذخیره‌سازی بلندمدت**
```yaml
# پیکربندی Prometheus برای ارسال داده‌ها به Thanos
remote_write:
  - url: "http://thanos-receive:10908/api/v1/receive"
```

در این مثال، Prometheus داده‌ها را به Thanos ارسال می‌کند تا در ذخیره‌سازی‌های ابری (مانند S3) ذخیره شود.

#### ۲. **استفاده از Alertmanager برای ارسال هشدار به Slack**
```yaml
# پیکربندی Alertmanager برای ارسال هشدار به Slack
route:
  receiver: 'slack-notifications'
receivers:
  - name: 'slack-notifications'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/...'
        channel: '#alerts'
```

در این مثال، Alertmanager هشدارهای Prometheus را به یک کانال Slack ارسال می‌کند.

#### ۳. **استفاده از Prometheus Remote Write برای ارسال داده‌ها به InfluxDB**
```yaml
# پیکربندی Prometheus برای ارسال داده‌ها به InfluxDB
remote_write:
  - url: "http://influxdb:8086/api/v1/prom/write?db=prometheus"
```

در این مثال، Prometheus داده‌ها را به InfluxDB ارسال می‌کند.

---

### **نتیجه:**
‏Bridgeها در Prometheus نقش مهمی در ادغام این سیستم با سایر پلتفرم‌ها و ابزارها ایفا می‌کنند. این Bridgeها به Prometheus اجازه می‌دهند تا:
- داده‌ها را به‌صورت بلندمدت ذخیره کند.
- هشدارها را به سیستم‌های خارجی ارسال کند.
- با سایر سیستم‌های مانیتورینگ ارتباط برقرار کند.
- داده‌ها را به ابزارهای تجسمی مانند Grafana ارسال کند.

با استفاده از Bridgeها، می‌توان قابلیت‌های Prometheus را گسترش داد و آن را به یک سیستم مانیتورینگ کامل و یکپارچه تبدیل کرد.


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