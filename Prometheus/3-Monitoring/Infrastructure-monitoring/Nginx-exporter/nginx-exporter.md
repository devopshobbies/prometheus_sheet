# Setting Up Nginx with Prometheus Monitoring on Your Laptop

Here's a comprehensive guide to setting up Nginx on your laptop and monitoring it with Prometheus:

![image](src/nginx_exporter_img.png)
## 1. Install Nginx

### On Linux (Ubuntu/Debian):
```bash
sudo apt update
sudo apt install nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

### On macOS (using Homebrew):
```bash
brew install nginx
brew services start nginx
```

### On Windows:
Download from [nginx.org](https://nginx.org/en/download.html) and run the executable

## 2. Enable Nginx Status Module

Edit your nginx configuration (typically `/etc/nginx/nginx.conf` or `/usr/local/etc/nginx/nginx.conf`):

```nginx
server {
    listen 8089;
    server_name localhost;
    
    location /nginx_status {
        stub_status on;
        access_log off;
        allow 127.0.0.1;
        deny all;
    }
}
```

Test and reload nginx:
```bash
sudo nginx -t
sudo nginx -s reload
```

## 3. Install Prometheus

### Download and install Prometheus:
```bash
# Download (adjust version as needed)
wget https://github.com/prometheus/prometheus/releases/download/v2.47.0/prometheus-2.47.0.linux-amd64.tar.gz
tar xvfz prometheus-*.tar.gz
cd prometheus-*
```

### Configure Prometheus (`prometheus.yml`):

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'nginx'
    static_configs:
      - targets: ['localhost:8089']
```

### Start Prometheus:
```bash
./prometheus --config.file=prometheus.yml
```

## 4. Install and Configure Nginx Exporter

The nginx status page needs to be parsed for Prometheus. Use the nginx-exporter:

```bash
# Download nginx-exporter
wget https://github.com/nginx/nginx-prometheus-exporter/releases/download/v1.4.1/nginx-prometheus-exporter_1.4.1_linux_amd64.tar.gz
tar xvfz nginx-prometheus-exporter_*.tar.gz
cd nginx-prometheus-exporter_*
```

Run the exporter:
```bash
./nginx-prometheus-exporter -nginx.scrape-uri http://localhost:8089/nginx_status
```

Update your `prometheus.yml` to scrape the exporter:
```yaml
scrape_configs:
  - job_name: 'nginx-exporter'
    static_configs:
      - targets: ['localhost:9113']
```

## 5. Access the Dashboards

- Prometheus: http://localhost:9090
- Nginx status: http://localhost:8089/nginx_status
- Nginx exporter metrics: http://localhost:9113/metrics

## 6. (Optional) Install Grafana for Visualization

```bash
# On Ubuntu/Debian
sudo apt-get install -y apt-transport-https
sudo apt-get install -y software-properties-common wget
wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -
echo "deb https://packages.grafana.com/oss/deb stable main" | sudo tee -a /etc/apt/sources.list.d/grafana.list
sudo apt-get update
sudo apt-get install grafana
sudo systemctl start grafana-server
```

Access Grafana at http://localhost:3000 (default credentials admin/admin)

Import Nginx dashboard (ID 12708) for pre-built visualizations.

## 7. Generate Load (for testing)

You can use tools like `ab` (Apache Benchmark) or `wrk` to generate load:

```bash
# Install ab (Apache Benchmark)
sudo apt install apache2-utils

# Generate load
ab -n 1000 -c 10 http://localhost/
```

This setup will give you metrics like:
- Requests per second
- Active connections
- Connection states (reading, writing, waiting)
- Response metrics

Would you like me to elaborate on any specific part of this setup?


------------

# Load Testing with k6

k6 is a modern, developer-centric load testing tool that's more powerful than `ab` (Apache Benchmark). Here's how to use k6 to send load to your Nginx server:

## Installation

### Linux (Debian/Ubuntu):
```bash
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt update
sudo apt install k6
```

### macOS (Homebrew):
```bash
brew install k6
```

### Windows (Chocolatey):
```bash
choco install k6
```

## Basic Load Test (100 seconds)

Create a file `test.js`:
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 50 },  // ramp-up
    { duration: '80s', target: 100 }, // main load
    { duration: '10s', target: 0 },   // ramp-down
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],   // <1% errors
    http_req_duration: ['p(95)<500'], // 95% of requests <500ms
  },
};

export default function () {
  let res = http.get('http://localhost:8089/');
  check(res, {
    'status was 200': (r) => r.status == 200,
    'response time OK': (r) => r.timings.duration < 1000,
  });
  sleep(0.1); // slight pause between requests
}
```

Run the test:
```bash
k6 run test.js
```

## Advanced Scenarios

### 1. Different HTTP Methods:
```javascript
import http from 'k6/http';

export default function () {
  const payload = JSON.stringify({ username: 'test', password: 'test' });
  const params = { headers: { 'Content-Type': 'application/json' } };
  
  http.get('http://localhost:8080/api');
  http.post('http://localhost:8080/login', payload, params);
}
```

### 2. Variable Load Patterns:
```javascript
export let options = {
  stages: [
    { duration: '2m', target: 100 },  // spike to 100 VUs
    { duration: '1m', target: 50 },   // scale down
    { duration: '2m', target: 200 },  // spike higher
    { duration: '1m', target: 0 },    // ramp-down
  ],
};
```

### 3. CSV Data Input:
```javascript
import { SharedArray } from 'k6/data';
const testData = new SharedArray('users', function() {
  return open('./users.csv').split('\n').slice(1);
});

export default function () {
  const [username, password] = testData[__VU % testData.length].split(',');
  // use username/password in requests
}
```

## Monitoring Integration

### 1. Output to Prometheus:
```bash
k6 run --out experimental-prometheus-rw test.js
```

### 2. Real-time Web Dashboard:
```bash
k6 run -o cloud test.js  # requires k6 cloud account
```

## Key Features of k6:

1. **JavaScript-based** - Write complex test scenarios
2. **CLI-first** - Easy to integrate into CI/CD pipelines
3. **Metrics-rich** - Detailed performance data out of the box
4. **Thresholds** - Define pass/fail criteria for your tests

Would you like me to elaborate on any specific aspect of k6 testing, such as:
- Testing authenticated endpoints?
- Handling WebSockets?
- Integrating with Grafana for visualization?
- Distributed load testing across multiple machines?
