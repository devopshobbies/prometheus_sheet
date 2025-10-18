
ref:
- https://prometheus.io/docs/introduction/first_steps/
- https://prometheus.io/docs/tutorials/getting_started/
```

wget https://github.com/prometheus/prometheus/releases/download/v3.0.0-rc.0/prometheus-3.0.0-rc.0.linux-amd64.tar.gz
useradd --no-create-home --shell /bin/false prometheus
mkdir /etc/prometheus
mkdir /var/lib/prometheus
chown prometheus:prometheus /etc/prometheus
chown prometheus:prometheus /var/lib/prometheus

tar -xzf prometheus-*.linux-amd64.tar.gz
cd prometheus-*.linux-amd64/
cp prometheus /usr/local/bin/
cp promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r consoles /etc/prometheus

cp -r console_libraries /etc/prometheus

chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries

```

nano /etc/prometheus/prometheus.yml

```yml
global:
scrape_interval: 10s
scrape_configs:
- job_name: Prometheus_master
static_configs:
- targets: ['localhost:9090']
```

nano /etc/systemd/system/prometheus.service

```yml
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target
[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
--config.file /etc/prometheus/prometheus.yml \
--storage.tsdb.path /var/lib/prometheus/ \
--web.console.templates=/etc/prometheus/consoles \
--web.console.libraries=/etc/prometheus/console_libraries
[Install]
WantedBy=multi-user.target
```


```
sudo 


systemctl daemon-reload
sudo 


systemctl enable --now prometheus
sudo 


systemctl status prometheus
sudo 


ufw --add-port=9090/tcp --permanent
sudo 


ufw--reload
```


### website

https://prometheus.io/docs/prometheus/latest/getting_started/

tmux ls
tmux new -s mysession
tmux at -t mysession
tmux a -t mysession
tmux kill-session -a -t mysession

### مثال کاربردی
راه اندازی سرور  prometheus و ارتباط آن با node expoter  در local 

- نقشه راه https://prometheus.io/docs/prometheus/latest/getting_started/
- مراحل نصب سرور صفحه دانلود ها https://prometheus.io/download/
- مراحل نصب اکسپورتر https://prometheus.io/docs/guides/node-exporter/
- صفحه کانفیگ ها https://prometheus.io/docs/prometheus/latest/configuration/configuration/
- معرفی سرویس کردن  https://sbcode.net/prometheus/prometheus-service/


---
## معنی aggregating در prometheus چیست؟

در **Prometheus**، مفهوم **aggregating** یا **تجمیع** به معنای جمع‌آوری و ترکیب داده‌ها برای انجام تحلیل‌های پیچیده‌تر یا ساده‌سازی نمایش داده‌ها است. این کار اغلب با استفاده از **عملگرهای aggregation** (مانند `sum`, `avg`, `max`, `min`, `count`) انجام می‌شود که در زبان **PromQL** (Prometheus Query Language) استفاده می‌شوند.

### اهداف اصلی Aggregation
۱. **خلاصه‌سازی داده‌ها:** کاهش حجم داده‌ها برای نمایش یا تحلیل.
۲. **محاسبات آماری:** استخراج میانگین، بیشینه، کمینه و سایر شاخص‌های آماری.
۳. **مقایسه و تحلیل گروهی:** انجام عملیات روی گروه‌های مختلف از داده‌ها بر اساس برچسب‌ها (labels).

---

### ساختار Aggregation در PromQL
عملگرهای aggregation در PromQL به دو صورت استفاده می‌شوند:
۱. **ساده:** بدون گروه‌بندی (`by`)
۲. **با گروه‌بندی:** استفاده از `by` یا `without` برای تعیین گروه‌بندی داده‌ها.

### مثال‌ها

#### ۱. مجموع مقدار یک متریک
```promql
sum(http_requests_total)
```
این کوئری مجموع کل درخواست‌های HTTP را بدون توجه به برچسب‌ها محاسبه می‌کند.

#### ۲. مجموع مقدار یک متریک با گروه‌بندی
```promql
sum(http_requests_total) by (method)
```
در این مثال، درخواست‌های HTTP بر اساس برچسب `method` گروه‌بندی شده و سپس مجموع هر گروه محاسبه می‌شود.

#### ۳. میانگین مقادیر یک متریک
```promql
avg(cpu_usage)
```
این کوئری میانگین استفاده از CPU را محاسبه می‌کند.

#### ۴. حذف گروه‌بندی بر اساس برچسب‌ها
```promql
sum(http_requests_total) without (instance)
```
این کوئری مجموع درخواست‌های HTTP را محاسبه کرده اما برچسب `instance` را نادیده می‌گیرد.

---

### عملگرهای Aggregation در PromQL
Prometheus چندین عملگر تجمیعی ارائه می‌دهد که برخی از آنها عبارتند از:
- `sum`: مجموع مقادیر.
- `avg`: میانگین مقادیر.
- `max`: بزرگترین مقدار.
- `min`: کوچکترین مقدار.
- `count`: تعداد داده‌ها.
- `stddev`: انحراف استاندارد مقادیر.
- `stdvar`: واریانس مقادیر.
- `topk`: انتخاب بالاترین `k` مقدار.
- `bottomk`: انتخاب پایین‌ترین `k` مقدار.

---

### نکات مهم:
- **برچسب‌ها (Labels):** برچسب‌ها نقش مهمی در گروه‌بندی داده‌ها دارند. بدون استفاده از `by` یا `without`، همه برچسب‌ها نادیده گرفته می‌شوند.
- **کارایی:** عملیات تجمیع روی تعداد زیادی داده ممکن است زمان‌بر باشد؛ بنابراین طراحی کوئری بهینه اهمیت زیادی دارد.

