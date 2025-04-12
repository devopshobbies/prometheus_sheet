### **اکسپورتِر Grok**

 
همه اپلیکیشن‌ها متریک‌ها (metrics) را به شکلی تولید نمی‌کنند که بتوان با استفاده از یک اکسپورتر (exporter) آن‌ها را به چیزی تبدیل کرد که پرومتئوس (Prometheus) بفهمد. اما چنین اپلیکیشن‌هایی ممکن است لاگ (logs) تولید کنند، و اکسپورتر Grok می‌تواند برای تبدیل آن‌ها به متریک‌ها استفاده شود. Grok روشی برای پارس کردن (parse) لاگ‌های بدون ساختار است که معمولاً با Logstash استفاده می‌شود. اکسپورتر Grok از همان زبان الگو (pattern language) استفاده مجدد می‌کند، که به شما امکان می‌دهد الگوهایی را که از قبل دارید، دوباره به کار ببرید.

### بررسی موارد مشابه

- Grok
- **Fluentd + Prometheus Plugin**
- **Promtail (Loki's Log Agent)**
- **mtail**

### **جدول مقایسه**

| Tool              | Best For                                  | Language | Supports Histograms? | Push/Pull |
| ----------------- | ----------------------------------------- | -------- | -------------------- | --------- |
| **Grok Exporter** | Simple log parsing                        | Go       | ✅                    | Pull      |
| **mtail**         | Lightweight log scraping                  | Go       | ✅                    | Pull      |
| **Promtail**      | Loki users (but can export to Prometheus) | Go       | ✅                    | Pull      |
| **Vector**        | High-performance pipelines                | Rust     | ✅                    | Push/Pull |
| **Fluentd**       | Complex log processing                    | Ruby     | ✅                    | Push      |

###  مثال

فرض کنید لاگ ساده‌ای مانند زیر دارید:
```
GET /foo 6.66
GET /bar 5.5
POST /foo 7.77777
POST /pacman 22
POST /pacman 0.22 
POST /pacman 0.25
POST /pacman 0.21
```
که در فایلی به نام `example.log` قرار دارد. می‌توانید با استفاده از اکسپورتِر Grok این لاگ‌ها را به متریک تبدیل کنید. ابتدا، نسخه 0.2.8 اکسپورتر Grok برای Linux amd64 را دانلود و از حالت فشرده خارج کنید. سپس، فایلی به نام `grok.yml` با محتوای مثال ۱۰-۷ ایجاد کنید.

**مثال زیر. فایل `grok.yml` برای پارس کردن یک فایل لاگ ساده و تولید متریک‌ها**
```yaml
global:
  config_version: 2
input:
  type: file
  path: example.log
  readall: true # در محیط پروداکشن از false استفاده کنید
grok:
  additional_patterns:
    - 'METHOD [A-Z]+'
    - 'PATH [^ ]+'
    - 'NUMBER [0-9.]+'
metrics:
  - type: counter
    name: log_http_requests_total
    help: HTTP requests
    match: '%{METHOD} %{PATH:path} %{NUMBER:latency}'
    labels:
      path: '{{.path}}'
  - type: histogram
    name: log_http_request_latency_seconds_total
    help: HTTP request latency
    match: '%{METHOD} %{PATH:path} %{NUMBER:latency}'
    value: '{{.latency}}'
server:
  port: 9144
```

در نهایت، اکسپورتر Grok را اجرا کنید:
```bash
./grok_exporter -config grok.yml
```

بیایید این را جزء به جزء بررسی کنیم. ابتدا، مقداری کد قالبی/استاندارد وجود دارد:
```yaml
global:
  config_version: 2
```
سپس، باید فایلی را که قرار است خوانده شود، تعریف کنید. اینجا ما از `readall: true` استفاده می‌کنیم، بنابراین نتایج مشابهی با این مثال خواهید دید. در محیط پروداکشن، آن را روی مقدار پیش‌فرض `false` باقی می‌گذارید تا فایل دنبال شود (فقط تغییرات جدید خوانده شوند):
```yaml
input:
  type: file
  path: example.log
  readall: true # در محیط پروداکشن از false استفاده کنید
```
ءGrok با الگوهایی مبتنی بر عبارات منظم (regular expressions) کار می‌کند. ما تمام الگوهای خود را اینجا به صورت دستی تعریف کرده‌ایم تا بتوانید بهتر درک کنید چه اتفاقی می‌افتد، اما می‌توانید از الگوهایی که از قبل دارید نیز استفاده مجدد کنید:
```yaml
grok:
  additional_patterns:
    - 'METHOD [A-Z]+'
    - 'PATH [^ ]+'
    - 'NUMBER [0-9.]+'
```
ما دو متریک داریم. اولی یک شمارنده (counter) به نام `log_http_requests_total` است که دارای برچسب (label) `path` می‌باشد:
```yaml
metrics:
  - type: counter
    name: log_http_requests_total
    help: HTTP requests
    match: '%{METHOD} %{PATH:path} %{NUMBER:latency}'
    labels:
      path: '{{.path}}'
```
دومی یک هیستوگرام (histogram) به نام `log_http_request_latency_seconds_total` است که مقدار تأخیر (latency) را مشاهده می‌کند و هیچ برچسبی ندارد:
```yaml
  - type: histogram
    name: log_http_request_latency_seconds_total
    help: HTTP request latency
    match: '%{METHOD} %{PATH:path} %{NUMBER:latency}'
    value: '{{.latency}}'
```
در نهایت، مشخص می‌کنیم که اکسپورتر متریک‌های خود را در کجا ارائه دهد:
```yaml
server:
  port: 9144
```
هنگامی که به `http://localhost:9144` مراجعه می‌کنید، در میان خروجی آن متریک‌های زیر را خواهید یافت:
```
# HELP log_http_request_latency_seconds_total HTTP request latency
# TYPE log_http_request_latency_seconds_total histogram
log_http_request_latency_seconds_total_bucket{le="0.005"} 0
log_http_request_latency_seconds_total_bucket{le="0.01"} 0
log_http_request_latency_seconds_total_bucket{le="0.025"} 0
log_http_request_latency_seconds_total_bucket{le="0.05"} 0
log_http_request_latency_seconds_total_bucket{le="0.1"} 0
log_http_request_latency_seconds_total_bucket{le="0.25"} 3
log_http_request_latency_seconds_total_bucket{le="0.5"} 3
log_http_request_latency_seconds_total_bucket{le="1"} 3
log_http_request_latency_seconds_total_bucket{le="2.5"} 3
log_http_request_latency_seconds_total_bucket{le="5"} 3
log_http_request_latency_seconds_total_bucket{le="10"} 6
log_http_request_latency_seconds_total_bucket{le="+Inf"} 7
log_http_request_latency_seconds_total_sum 42.61777
log_http_request_latency_seconds_total_count 7
# HELP log_http_requests_total HTTP requests
# TYPE log_http_requests_total counter
log_http_requests_total{path="/bar"} 1
log_http_requests_total{path="/foo"} 2
log_http_requests_total{path="/pacman"} 4
```
همانطور که می‌بینید، پیکربندی اکسپورتر Grok نسبت به یک اکسپورتر معمولی پیچیده‌تر است؛ از نظر میزان تلاش، به ابزارگذاری مستقیم (direct instrumentation) نزدیک‌تر است، زیرا باید هر متریکی را که می‌خواهید ارائه دهید، به صورت جداگانه تعریف کنید. معمولاً برای هر نمونه (instance) اپلیکیشن که نیاز به مانیتورینگ دارد، یک اکسپورتر Grok اجرا می‌کنید و داده‌های آن را به روش معمول با پرومتئوس جمع‌آوری (scrape) می‌کنید، همانطور که در مثال زیر نشان داده شده است.

**مثال زیر. فایل `prometheus.yml` برای جمع‌آوری داده از یک اکسپورتر Grok محلی**
```yaml
global:
  scrape_interval: 10s
scrape_configs:
  - job_name: grok
    static_configs:
      - targets:
          - localhost:9144
```



--------













