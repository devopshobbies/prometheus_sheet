## منابع مفید:

https://training.promlabs.com/training/relabeling/training-overview/introduction/
https://relabeler.promlabs.com


‏**Relabeling** در Prometheus یک قابلیت قدرتمند است که به شما امکان می‌دهد برچسب‌ها (Labels) را قبل از ذخیره‌سازی داده‌ها در پایگاه داده‌ی Prometheus تغییر دهید، حذف کنید یا اضافه کنید. این فرآیند معمولاً در مراحل مختلف جمع‌آوری داده‌ها (مانند scrape) انجام می‌شود و به شما کمک می‌کند تا داده‌های خود را بهتر سازماندهی کنید، فیلتر کنید یا تغییر دهید.

### کاربردهای Relabeling
1. **اضافه کردن برچسب‌های جدید**: می‌توانید برچسب‌های جدیدی به متریک‌ها اضافه کنید تا اطلاعات بیشتری را در اختیار قرار دهید.
2. **حذف برچسب‌های ناخواسته**: برچسب‌هایی که نیاز ندارید را می‌توانید حذف کنید.
3. **تغییر مقدار برچسب‌ها**: مقادیر برچسب‌ها را می‌توانید تغییر دهید یا اصلاح کنید.
4. **فیلتر کردن هدف‌ها (Targets)**: می‌توانید تعیین کنید که کدام هدف‌ها (Targets) scrape شوند و کدام‌ها نادیده گرفته شوند.

### مراحل انجام Relabeling
‏Relabeling در چند مرحله‌ی مختلف در Prometheus انجام می‌شود:
1. **در هنگام کشف سرویس‌ها (Service Discovery)**: قبل از scrape کردن هدف‌ها، Prometheus می‌تواند برچسب‌ها را تغییر دهد.
2. **در هنگام scrape کردن**: بعد از جمع‌آوری داده‌ها از یک هدف، Prometheus می‌تواند برچسب‌ها را تغییر دهد.
3. **در هنگام ذخیره‌سازی داده‌ها**: قبل از ذخیره‌سازی داده‌ها در پایگاه داده، Prometheus می‌تواند برچسب‌ها را تغییر دهد.

### تنظیمات Relabeling
‏Relabeling با استفاده از تنظیمات `relabel_configs` و `metric_relabel_configs` در فایل پیکربندی Prometheus انجام می‌شود:
‏- **`relabel_configs`**: برای تغییر برچسب‌ها قبل از scrape کردن هدف‌ها استفاده می‌شود.
-‏ **`metric_relabel_configs`**: برای تغییر برچسب‌ها بعد از scrape کردن و قبل از ذخیره‌سازی داده‌ها استفاده می‌شود.

### مثال‌هایی از Relabeling

#### ۱. اضافه کردن یک برچسب جدید
```yaml
scrape_configs:
  - job_name: 'example'
    static_configs:
      - targets: ['localhost:9090']
    relabel_configs:
      - source_labels: [__address__]
        target_label: 'instance'
        replacement: 'my_instance'
```
در این مثال، برچسب `instance` با مقدار `my_instance` به همه‌ی هدف‌ها اضافه می‌شود.

#### ۲. حذف یک برچسب
```yaml
scrape_configs:
  - job_name: 'example'
    static_configs:
      - targets: ['localhost:9090']
    metric_relabel_configs:
      - action: 'labeldrop'
        regex: 'unwanted_label'
```
در این مثال، برچسب‌هایی که نام آن‌ها با `unwanted_label` مطابقت دارد، حذف می‌شوند.

#### ۳. تغییر مقدار یک برچسب
```yaml
scrape_configs:
  - job_name: 'example'
    static_configs:
      - targets: ['localhost:9090']
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_name]
        target_label: 'pod_name'
```
در این مثال، مقدار برچسب `__meta_kubernetes_pod_name` به `pod_name` تغییر می‌کند.

#### ۴. فیلتر کردن هدف‌ها
```yaml
scrape_configs:
  - job_name: 'example'
    kubernetes_sd_configs:
      - role: 'pod'
    relabel_configs:
      - source_labels: [__meta_kubernetes_namespace]
        action: 'keep'
        regex: 'production'
```
در این مثال، فقط هدف‌هایی که در namespaceی با نام `production` قرار دارند، scrape می‌شوند.

###  بررسی Relabeling درعمل 
‏- **`replace`**: مقدار یک برچسب را تغییر می‌دهد (پیش‌فرض).
‏- **`keep`**: فقط هدف‌هایی که با regex مطابقت دارند را نگه می‌دارد.
‏- **`drop`**: هدف‌هایی که با regex مطابقت دارند را حذف می‌کند.
‏- **`labelmap`**: نام برچسب‌ها را بر اساس regex تغییر می‌دهد.
‏- **`labeldrop`**: برچسب‌هایی که با regex مطابقت دارند را حذف می‌کند.
‏- **`labelkeep`**: فقط برچسب‌هایی که با regex مطابقت دارند را نگه می‌دارد.

‏Relabeling در Prometheus یک ابزار بسیار مفید برای مدیریت و سازماندهی برچسب‌ها است. با استفاده از آن می‌توانید داده‌های خود را بهتر کنترل کنید، برچسب‌های اضافی را حذف کنید یا برچسب‌های جدیدی اضافه کنید. این قابلیت به شما کمک می‌کند تا سیستم مانیتورینگ خود را انعطاف‌پذیرتر و کارآمدتر کنید.


![[src/relabling2.png]]

| Relabeled Object Type                          | Configuration Section                                                                                                                                                |
| ---------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Discovered **scrape targets**                  | `relabel_configs` section in a [`scrape_config`](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config) section                    |
| Discovered **Alertmanager targets**            | `relabel_configs` section in an [`alertmanager_config`](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#alertmanager_config) section       |
| Individual **samples from scrapes**            | `metric_relabel_configs` section in a [`scrape_config`](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config) section             |
| **Alerts** sent to _all_ Alertmanagers         | `alert_relabel_configs` section in the `alerting` section of the main configuration file                                                                             |
| **Alerts** sent to a _subset_ of Alertmanagers | `alert_relabel_configs` section in an [`alertmanager_config`](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#alertmanager_config) section |
| **Samples written to remote storage systems**  | `write_relabel_configs` section in the [`remote_write`](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) section              |

---

### **چرا Relabeling مهم است؟**
در محیط‌های دینامیک و پیچیده (مانند Kubernetes یا سیستم‌های مبتنی بر میکروسرویس)، ممکن است با حجم زیادی از هدف‌ها (Targets) و برچسب‌ها (Labels) مواجه شوید. Relabeling به شما کمک می‌کند:
1. **داده‌های غیرضروری را حذف کنید**: با فیلتر کردن هدف‌ها یا برچسب‌های ناخواسته، حجم داده‌های ذخیره‌شده را کاهش دهید.
2. **داده‌ها را سازماندهی کنید**: برچسب‌های معنادار اضافه کنید تا داده‌ها بهتر قابل جستجو و تحلیل باشند.
3. **اطلاعات اضافی اضافه کنید**: از طریق متادیتا (Metadata) یا اطلاعات محیطی، برچسب‌های جدید ایجاد کنید.
4. **هدف‌ها را فیلتر کنید**: فقط هدف‌های خاصی را scrape کنید و از scrape کردن هدف‌های غیرضروری جلوگیری کنید.

---

### **مراحل انجام Relabeling**
‏Relabeling در Prometheus در چند مرحله‌ی کلیدی انجام می‌شود:

1. **کشف سرویس (Service Discovery)**:
   - ‏Prometheus از مکانیزم‌های کشف سرویس (مانند Kubernetes, Consul, DNS و غیره) برای یافتن هدف‌ها (Targets) استفاده می‌کند.
   - در این مرحله، Prometheus برچسب‌های موقت (مثل `__meta_*`) را به هدف‌ها اضافه می‌کند. این برچسب‌ها حاوی اطلاعاتی مانند آدرس IP، نام سرویس، namespace و غیره هستند.

2. **قبل از Scrape (Relabeling با `relabel_configs`)**:
   - در این مرحله، Prometheus از `relabel_configs` استفاده می‌کند تا برچسب‌ها را تغییر دهد یا هدف‌ها را فیلتر کند.
   - این مرحله برای اضافه کردن برچسب‌های جدید، تغییر مقادیر برچسب‌ها یا حذف هدف‌های ناخواسته استفاده می‌شود.

3. **بعد از Scrape (Relabeling با `metric_relabel_configs`)**:
   - بعد از scrape کردن داده‌ها، Prometheus از `metric_relabel_configs` استفاده می‌کند تا برچسب‌های متریک‌ها را تغییر دهد یا حذف کند.
   - این مرحله برای پاک‌سازی داده‌ها قبل از ذخیره‌سازی در پایگاه داده استفاده می‌شود.

---

### **اجزای اصلی Relabeling**
هر role Relabeling شامل چند جزء کلیدی است:

‏-  **`source_labels`**: لیستی از label‌هایی که مقادیر آن‌ها به‌عنوان ورودی استفاده می‌شود.
-‏ **`separator`**: کاراکتری که برای ترکیب مقادیر `source_labels` استفاده می‌شود (پیش‌فرض: `;`).
-‏ **`target_label`**: target_labelای که مقدار جدید به آن اختصاص داده می‌شود.
-‏ **`regex`**: یک عبارت منظم (Regular Expression) برای تطبیق مقادیر.
-‏ **`replacement`**: مقداری که جایگزین مقدار تطبیق‌داده‌شده با `regex` می‌شود.
-‏ **`action`**: عملی که باید انجام شود (مانند `replace`, `keep`, `drop`, `labelmap` و غیره).

---

### **انواع Actions در Relabeling**
‏Prometheus از انواع مختلفی از Actions پشتیبانی می‌کند که هر کدام کاربرد خاصی دارند:

‏- **`replace`** (پیش‌فرض):
   - مقدار یک برچسب را تغییر می‌دهد.
   - مثال:
     ```yaml
     - source_labels: [__meta_kubernetes_pod_name]
       target_label: pod_name
     ```

. **`keep`**:
   - فقط هدف‌هایی که با `regex` مطابقت دارند را نگه می‌دارد.
   - مثال:
     ```yaml
     - source_labels: [__meta_kubernetes_namespace]
       action: keep
       regex: production
     ```

. **`drop`**:
   - هدف‌هایی که با `regex` مطابقت دارند را حذف می‌کند.
   - مثال:
     ```yaml
     - source_labels: [__meta_kubernetes_namespace]
       action: drop
       regex: test
     ```

 **`labelmap`**:
   - نام برچسب‌ها را بر اساس `regex` تغییر می‌دهد.
   - مثال:
     ```yaml
     - action: labelmap
       regex: __meta_kubernetes_(.+)
     ```

 **`labeldrop`**:
   - برچسب‌هایی که با `regex` مطابقت دارند را حذف می‌کند.
   - مثال:
     ```yaml
     - action: labeldrop
       regex: unwanted_label
     ```

 **`labelkeep`**:
   - فقط برچسب‌هایی که با `regex` مطابقت دارند را نگه می‌دارد.
   - مثال:
     ```yaml
     - action: labelkeep
       regex: important_label
     ```

---

### **مثال‌های پیشرفته‌تر**


![[src/relabing3.png]]
#### ۱. اضافه کردن برچسب‌های دینامیک از Kubernetes
```yaml
scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_name]
        target_label: pod_name
      - source_labels: [__meta_kubernetes_namespace]
        target_label: namespace
      - source_labels: [__meta_kubernetes_pod_ip]
        target_label: pod_ip
```

در این مثال، Prometheus از متادیتای Kubernetes استفاده می‌کند تا برچسب‌های `pod_name`, `namespace` و `pod_ip` را به متریک‌ها اضافه کند.

#### ۲. فیلتر کردن هدف‌ها بر اساس Namespace
```yaml
scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_namespace]
        action: keep
        regex: production
```

در این مثال، فقط پادهای (Pods) که در namespaceی با نام `production` قرار دارند، scrape می‌شوند.

#### ۳. تغییر نام برچسب‌ها با `labelmap`
```yaml
scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - action: labelmap
        regex: __meta_kubernetes_(.+)
```

در این مثال، همه‌ی برچسب‌هایی که با `__meta_kubernetes_` شروع می‌شوند، به‌صورت `kubernetes_` تغییر نام می‌دهند.

---

### **جمع‌بندی**
‏Relabeling در Prometheus یک ابزار بسیار قدرتمند برای مدیریت و سازماندهی داده‌ها است. با استفاده از آن می‌توانید:
- برچسب‌های جدید اضافه کنید.
- برچسب‌های ناخواسته را حذف کنید.
- هدف‌ها را فیلتر کنید.
- نام برچسب‌ها را تغییر دهید.

این قابلیت به شما کمک می‌کند تا سیستم مانیتورینگ خود را انعطاف‌پذیرتر و کارآمدتر کنید، به‌ویژه در محیط‌های دینامیک مانند Kubernetes.


---

### **چگونه Relabeling کار می‌کند؟**
‏Relabeling در Prometheus بر اساس یک سری قوانین (Rules) انجام می‌شود. هر قانون شامل چند جزء اصلی است که Prometheus از آن‌ها برای تغییر یا فیلتر کردن برچسب‌ها استفاده می‌کند. این قوانین در دو مرحله‌ی اصلی اعمال می‌شوند:

1. **قبل از Scrape کردن (`relabel_configs`)**:
   - در این مرحله، Prometheus از `relabel_configs` استفاده می‌کند تا برچسب‌های هدف‌ها (Targets) را تغییر دهد یا هدف‌ها را فیلتر کند.
   - این مرحله برای اضافه کردن برچسب‌های جدید، تغییر مقادیر برچسب‌ها یا حذف هدف‌های ناخواسته استفاده می‌شود.

2. **بعد از Scrape کردن (`metric_relabel_configs`)**:
   - در این مرحله، Prometheus از `metric_relabel_configs` استفاده می‌کند تا برچسب‌های متریک‌ها را تغییر دهد یا حذف کند.
   - این مرحله برای پاک‌سازی داده‌ها قبل از ذخیره‌سازی در پایگاه داده استفاده می‌شود.

---

### **جزئیات بیشتر درباره Actions در Relabeling**
هر قانون Relabeling شامل یک `action` است که مشخص می‌کند چه عملی باید انجام شود. در ادامه، به بررسی دقیق‌تر هر یک از این Actions می‌پردازیم:

#### ۱. **`replace`** (پیش‌فرض)
- این Action مقدار یک برچسب را تغییر می‌دهد.
- از `source_labels` برای دریافت مقادیر ورودی و از `target_label` برای مشخص کردن برچسب هدف استفاده می‌کند.
- مثال:
  ```yaml
  - source_labels: [__meta_kubernetes_pod_name]
    target_label: pod_name
  ```

#### ۲. **`keep`**
- فقط هدف‌هایی که با `regex` مطابقت دارند را نگه می‌دارد.
- اگر هدفی با `regex` مطابقت نداشته باشد، scrape نمی‌شود.
- مثال:
  ```yaml
  - source_labels: [__meta_kubernetes_namespace]
    action: keep
    regex: production
  ```

#### ۳. **`drop`**
- هدف‌هایی که با `regex` مطابقت دارند را حذف می‌کند.
- اگر هدفی با `regex` مطابقت داشته باشد، scrape نمی‌شود.
- مثال:
  ```yaml
  - source_labels: [__meta_kubernetes_namespace]
    action: drop
    regex: test
  ```

#### ۴. **`labelmap`**
- نام برچسب‌ها را بر اساس `regex` تغییر می‌دهد.
- این Action برای تغییر نام برچسب‌های موقت (مانند `__meta_*`) به برچسب‌های دائمی استفاده می‌شود.
- مثال:
  ```yaml
  - action: labelmap
    regex: __meta_kubernetes_(.+)
  ```

#### ۵. **`labeldrop`**
- برچسب‌هایی که با `regex` مطابقت دارند را حذف می‌کند.
- این Action برای حذف برچسب‌های ناخواسته استفاده می‌شود.
- مثال:
  ```yaml
  - action: labeldrop
    regex: unwanted_label
  ```

#### ۶. **`labelkeep`**
- فقط برچسب‌هایی که با `regex` مطابقت دارند را نگه می‌دارد.
- این Action برای حفظ برچسب‌های مهم و حذف بقیه استفاده می‌شود.
- مثال:
  ```yaml
  - action: labelkeep
    regex: important_label
  ```

---

### **مثال‌های پیشرفته‌تر**

#### ۱. **اضافه کردن برچسب‌های دینامیک از Kubernetes**
```yaml
scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_name]
        target_label: pod_name
      - source_labels: [__meta_kubernetes_namespace]
        target_label: namespace
      - source_labels: [__meta_kubernetes_pod_ip]
        target_label: pod_ip
```

در این مثال، Prometheus از متادیتای Kubernetes استفاده می‌کند تا برچسب‌های `pod_name`, `namespace` و `pod_ip` را به متریک‌ها اضافه کند.

#### ۲. **فیلتر کردن هدف‌ها بر اساس Namespace**
```yaml
scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_namespace]
        action: keep
        regex: production
```

در این مثال، فقط پادهای (Pods) که در namespaceی با نام `production` قرار دارند، scrape می‌شوند.

#### ۳. **تغییر نام برچسب‌ها با `labelmap`**
```yaml
scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - action: labelmap
        regex: __meta_kubernetes_(.+)
```

در این مثال، همه‌ی برچسب‌هایی که با `__meta_kubernetes_` شروع می‌شوند، به‌صورت `kubernetes_` تغییر نام می‌دهند.

#### ۴. **حذف برچسب‌های ناخواسته**
```yaml
scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    metric_relabel_configs:
      - action: labeldrop
        regex: 'unwanted_label'
```

در این مثال، برچسب‌هایی که نام آن‌ها با `unwanted_label` مطابقت دارد، حذف می‌شوند.

---

### **بهترین روش‌های استفاده از Relabeling**
1. **استفاده از `relabel_configs` برای فیلتر کردن هدف‌ها**:
   - از `relabel_configs` برای فیلتر کردن هدف‌ها قبل از scrape کردن استفاده کنید. این کار باعث کاهش حجم داده‌های جمع‌آوری شده و بهبود عملکرد Prometheus می‌شود.

2. **استفاده از `metric_relabel_configs` برای پاک‌سازی داده‌ها**:
   - از `metric_relabel_configs` برای حذف برچسب‌های ناخواسته یا اضافه کردن برچسب‌های جدید قبل از ذخیره‌سازی داده‌ها استفاده کنید.

3. **استفاده از `labelmap` برای تغییر نام برچسب‌ها**:
   - اگر از کشف سرویس (Service Discovery) استفاده می‌کنید، از `labelmap` برای تغییر نام برچسب‌های موقت به برچسب‌های دائمی استفاده کنید.

4. **استفاده از `keep` و `drop` برای فیلتر کردن هدف‌ها**:
   - از `keep` و `drop` برای فیلتر کردن هدف‌ها بر اساس شرایط خاص (مانند namespace یا محیط) استفاده کنید.

5. **استفاده از `regex` برای تطبیق پیشرفته**:
   - از عبارات منظم (Regex) برای تطبیق دقیق‌تر مقادیر برچسب‌ها استفاده کنید.

---

### **جمع‌بندی**
‏Relabeling در Prometheus یک ابزار بسیار قدرتمند برای مدیریت و سازماندهی داده‌ها است. با استفاده از آن می‌توانید:
- برچسب‌های جدید اضافه کنید.
- برچسب‌های ناخواسته را حذف کنید.
- هدف‌ها را فیلتر کنید.
- نام برچسب‌ها را تغییر دهید.

این قابلیت به شما کمک می‌کند تا سیستم مانیتورینگ خود را انعطاف‌پذیرتر و کارآمدتر کنید، به‌ویژه در محیط‌های دینامیک مانند Kubernetes. با تسلط بر Relabeling، می‌توانید داده‌های خود را به‌صورت دقیق‌تر و معنادارتر مدیریت کنید.


---
