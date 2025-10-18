https://github.com/antonputra/tutorials/tree/main/lessons/141/prometheus-nginx-exporter

https://prometheus.io/docs/instrumenting/writing_exporters/
گاهی اوقات شما قادر نخواهید بود ابزار دقیق مستقیم (direct instrumentation) به یک برنامه اضافه کنید،
و یا یک اکسپورتر (exporter) موجود که آن را پوشش دهد پیدا کنید. این شما را با این وضعیت مواجه می‌کند که مجبور شوید خودتان یک اکسپورتر بنویسید. خبر خوب این است که نوشتن اکسپورترها نسبتاً آسان است.
بخش دشوار، فهمیدن معنای متریک‌هایی است که توسط برنامه‌ها افشا (expose) می‌شوند. واحدها (Units) اغلب نامشخص هستند و مستندات، اگر اصلاً وجود داشته باشند، می‌توانند مبهم باشند. در این قسمت شما یاد خواهید گرفت که چگونه اکسپورتر بنویسید.

**تله‌متری Consul**
ما قصد داریم یک اکسپورتر کوچک برای Consul بنویسیم تا این فرآیند را نشان دهیم. ما
قبلاً Consul و Consul Exporter را در بخش «Consul» در صفحه ۱۸۷ دیده‌ایم، بنابراین بیایید یک
اکسپورت ساده با متریک‌هایی از API تله‌متری آن ایجاد کنیم.¹

درحالی‌که شما می‌توانید اکسپورترها را به هر زبان برنامه‌نویسی بنویسید، اکثر آن‌ها به زبان
Go نوشته شده‌اند، و این زبانی است که ما در اینجا استفاده خواهیم کرد. با این حال، تعداد کمی
اکسپورت نوشته شده به زبان پایتون (Python) و تعداد بسیار کمتری به زبان جاوا (Java) پیدا خواهید کرد.
اگر Consul شما در حال اجرا نیست، دوباره آن را طبق دستورالعمل‌های مثال ۸-۸ اجرا کنید.
اگر به آدرس http://localhost:8500/v1/agent/metrics مراجعه کنید، خروجی JSON را که
با آن کار خواهید کرد مشاهده می‌کنید، که شبیه مثال ۱۲-۱ است. به راحتی، Consul
یک کتابخانه Go ارائه می‌دهد که می‌توانید از آن استفاده کنید، بنابراین لازم نیست نگران تجزیه (parsing)
JSON به صورت دستی باشید.

مثال ۱۲-۱. یک نمونه خروجی خلاصه شده از خروجی متریک‌های یک ایجنت Consul
{
"Timestamp": "2018-01-31 14:42:10 +0000 UTC",
"Gauges": [
{
"Name": "consul.autopilot.failure_tolerance",
"Value": 0,
"Labels": {}
}
],
"Points": [],
"Counters": [
{
"Name": "consul.raft.apply",
"Count": 1,
"Sum": 2, "Min": 1, "Max": 1, "Mean": 1, "Stddev": 0,
"Labels": {}
}
],
"Samples": [
{
"Name": "consul.fsm.coordinate.batch-update",
"Count": 1,
"Sum": 0.13156799972057343,
"Min": 0.13156799972057343, "Max": 0.13156799972057343,
"Mean": 0.13156799972057343, "Stddev": 0,
"Labels": {}
}
]
}

شما خوش‌شانس هستید که Consul شمارنده‌ها (counters) و گیج‌ها (gauges) را برای شما جدا کرده است.² به نظر می‌رسد
Samples (نمونه‌ها) را نیز می‌توانید با استفاده از Count و Sum در یک متریک خلاصه (summary) استفاده کنید. با نگاهی دوباره
به تمام Samples، شک داریم که آن‌ها تأخیر (latency) را ردیابی می‌کنند. با بررسی
مستندات تأیید می‌شود که آن‌ها تایمر (timer) هستند، که به معنی یک خلاصه (summary) پرومتئوس است
(به «خلاصه» در صفحه ۵۰ مراجعه کنید). تایمرها همچنین همگی بر حسب میلی‌ثانیه هستند، بنابراین
می‌توانیم آن‌ها را به ثانیه تبدیل کنیم.³ در حالی که JSON فیلدی برای برچسب‌ها (labels) دارد، هیچ‌کدام استفاده نشده‌اند،
بنابراین می‌توانید آن را نادیده بگیرید. گذشته از آن، تنها کار دیگری که باید انجام دهید این است که اطمینان حاصل کنید
هر کاراکتر نامعتبر در نام‌های متریک پاکسازی (sanitize) شده است.
اکنون منطقی را که باید روی متریک‌هایی که Consul افشا می‌کند اعمال کنید، می‌دانید، بنابراین
می‌توانید اکسپورتر خود را مانند مثال ۱۲-۲ بنویسید.

---
¹ API تله‌متری در Consul 1.0 معرفی شد.
² فقط به این دلیل که چیزی شمارنده نامیده می‌شود به این معنی نیست که یک شمارنده است. به عنوان مثال، Dropwizard دارای شمارنده‌هایی است که می‌توانند کاهش یابند، بنابراین بسته به نحوه استفاده از شمارنده در عمل، ممکن است از نظر پرومتئوس یک شمارنده، گیج یا بدون نوع (untyped) باشد.
³ اگر فقط برخی از Samples تایمر بودند، باید بین افشای آن‌ها به همان صورت یا نگهداری لیستی از اینکه کدام متریک‌ها تأخیر هستند و کدام‌ها نیستند، یکی را انتخاب می‌کردید.
---

```go
// مثال ۱۲-۲. consul_metrics.go، یک اکسپورتر برای متریک‌های Consul که به زبان Go نوشته شده است
package main
import (
"log"
"net/http"
"regexp"
"github.com/hashicorp/consul/api"
"github.com/prometheus/client_golang/prometheus"
"github.com/prometheus/client_golang/prometheus/promhttp"
)
var (
up = prometheus.NewDesc(
"consul_up", // نام متریک
"Was talking to Consul successful.", // راهنما (Help string)
nil, // نام‌های برچسب متغیر
nil, // برچسب‌های ثابت
)
invalidChars = regexp.MustCompile("[^a-zA-Z0-9:_]") // برای پاکسازی نام متریک
)
type ConsulCollector struct {
// فیلدهای مورد نیاز برای اتصال به Consul می‌توانند در اینجا قرار گیرند
}
// اینترفیس prometheus.Collector را پیاده‌سازی می‌کند.
func (c ConsulCollector) Describe(ch chan<- *prometheus.Desc) {
// توضیحات متریک‌هایی که این کالکتور می‌تواند تولید کند را ارسال می‌کند.
// این در زمان ثبت (registration) فراخوانی می‌شود.
ch <- up // حداقل یک Desc باید ارسال شود.
}
// اینترفیس prometheus.Collector را پیاده‌سازی می‌کند.
func (c ConsulCollector) Collect(ch chan<- prometheus.Metric) {
// این متد در هر بار اسکرِیپ (scrape) فراخوانی می‌شود.
// ۱. داده‌ها را از Consul واکشی کنید.
consul, err := api.NewClient(api.DefaultConfig()) // اتصال به Consul
if err != nil {
// اگر اتصال ناموفق بود، consul_up را 0 قرار دهید.
ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
log.Printf("Error creating Consul client: %v", err)
return
}
metrics, err := consul.Agent().Metrics() // دریافت متریک‌ها از Consul agent
if err != nil {
// اگر دریافت متریک‌ها ناموفق بود، consul_up را 0 قرار دهید.
ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
log.Printf("Error fetching Consul metrics: %v", err)
return
}
// اگر تا اینجا رسیدیم، اتصال و واکشی موفق بوده است.
ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)
// ۲. متریک‌ها را پردازش و به فرمت پرومتئوس تبدیل کنید.
// پردازش Gauges
for _, g := range metrics.Gauges {
name := invalidChars.ReplaceAllLiteralString(g.Name, "_") // پاکسازی نام
desc := prometheus.NewDesc("consul_"+name, // پیشوند consul_ اضافه شده برای وضوح
"Consul metric "+g.Name, // متن راهنما
nil, // بدون برچسب متغیر اضافی
g.Labels) // استفاده از برچسب‌های ارائه شده توسط Consul
ch <- prometheus.MustNewConstMetric(
desc, prometheus.GaugeValue, float64(g.Value)) // ارسال متریک گیج
}
// پردازش Counters
for _, c := range metrics.Counters {
name := invalidChars.ReplaceAllLiteralString(c.Name, "_") // پاکسازی نام
// افزودن پسوند _total طبق قرارداد پرومتئوس برای شمارنده‌ها
desc := prometheus.NewDesc("consul_"+name+"_total",
"Consul metric "+c.Name,
nil,
c.Labels)
ch <- prometheus.MustNewConstMetric(
desc, prometheus.CounterValue, float64(c.Count)) // ارسال متریک شمارنده (مقدار Count)
}
// پردازش Samples (به عنوان Summary)
for _, s := range metrics.Samples {
// همه نمونه‌ها زمان بر حسب میلی‌ثانیه هستند، آن‌ها را به ثانیه تبدیل می‌کنیم.
name := invalidChars.ReplaceAllLiteralString(s.Name, "_") +
"_seconds" // پاکسازی نام و افزودن پسوند _seconds
// ایجاد متریک _count برای Summary
countDesc := prometheus.NewDesc(
"consul_"+name+"_count", "Consul metric "+s.Name, nil, s.Labels)
ch <- prometheus.MustNewConstMetric(
countDesc, prometheus.CounterValue, float64(s.Count))
// ایجاد متریک _sum برای Summary
sumDesc := prometheus.NewDesc(
"consul_"+name+"_sum", "Consul metric "+s.Name, nil, s.Labels)
ch <- prometheus.MustNewConstMetric(
sumDesc, prometheus.CounterValue, s.Sum/1000) // تبدیل Sum به ثانیه
}
}
func main() {
// ۳. کالکتور را ثبت و سرور HTTP را راه‌اندازی کنید.
c := ConsulCollector{} // ایجاد یک نمونه از کالکتور
prometheus.MustRegister(c) // ثبت کالکتور در رجیستری پیش‌فرض
http.Handle("/metrics", promhttp.Handler()) // ایجاد هندلر برای مسیر /metrics
log.Println("Starting Consul exporter on :8000")
log.Fatal(http.ListenAndServe(":8000", nil)) // راه‌اندازی سرور HTTP
}
```

اگر یک محیط توسعه Go فعال دارید، می‌توانید اکسپورتر را با دستورات زیر اجرا کنید:
```bash
go get -d -u github.com/hashicorp/consul/api
go get -d -u github.com/prometheus/client_golang/prometheus
go run consul_metrics.go
```
اگر به آدرس http://localhost:8000/metrics مراجعه کنید، متریک‌هایی مانند زیر را خواهید دید:
```
# HELP consul_autopilot_failure_tolerance Consul metric consul.autopilot.failure_tolerance
# TYPE consul_autopilot_failure_tolerance gauge
consul_autopilot_failure_tolerance 0
# HELP consul_raft_apply_total Consul metric consul.raft.apply
# TYPE consul_raft_apply_total counter
consul_raft_apply_total 1
# HELP consul_fsm_coordinate_batch_update_seconds_count Consul metric consul.fsm.coordinate.batch-update
# TYPE consul_fsm_coordinate_batch_update_seconds_count counter
consul_fsm_coordinate_batch_update_seconds_count 1
# HELP consul_fsm_coordinate_batch_update_seconds_sum Consul metric consul.fsm.coordinate.batch-update
# TYPE consul_fsm_coordinate_batch_update_seconds_sum counter
consul_fsm_coordinate_batch_update_seconds_sum 1.3156799972057343e-01
# HELP consul_up Was talking to Consul successful.
# TYPE consul_up gauge
consul_up 1
# ... سایر متریک‌ها ...
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
# ... متریک‌های go_ و process_ ...
```
همه چیز خوب پیش رفت، اما کد چگونه کار می‌کند؟ در بخش بعدی به شما نشان خواهیم داد.

**کالکتورهای سفارشی (Custom Collectors)**
با ابزار دقیق مستقیم، کتابخانه کلاینت رویدادهای ابزار دقیق را دریافت کرده و مقادیر متریک‌ها را در طول زمان ردیابی می‌کند. کتابخانه‌های کلاینت متریک‌های شمارنده، گیج، خلاصه و هیستوگرام (histogram) را برای این منظور ارائه می‌دهند که همگی نمونه‌هایی از کالکتورها (collectors) هستند. در زمان اسکرِیپ (scrape)، هر کالکتور در یک رجیستری (registry) جمع‌آوری (collect) می‌شود، به این معنی که متریک‌هایش از آن درخواست می‌شود. این متریک‌ها سپس توسط اسکرِیپ `/metrics` بازگردانده می‌شوند. شمارنده‌ها و سه نوع متریک استاندارد دیگر همیشه فقط یک خانواده متریک (metric family) را باز می‌گردانند.

اگر به جای استفاده از ابزار دقیق مستقیم، بخواهید متریک‌ها را از منبع دیگری ارائه دهید، از یک کالکتور سفارشی (custom collector) استفاده می‌کنید، که هر کالکتوری است که یکی از چهار نوع استاندارد نباشد. کالکتورهای سفارشی می‌توانند هر تعداد خانواده متریک را بازگردانند. جمع‌آوری (Collection) در هر بار اسکرِیپ صفحه `/metrics` اتفاق می‌افتد، که در آن هر جمع‌آوری یک تصویر لحظه‌ای (snapshot) سازگار از متریک‌های یک کالکتور است.

در Go، کالکتورهای شما باید اینترفیس `prometheus.Collector` را پیاده‌سازی کنند. به این معنی که کالکتورها باید اشیائی با متدهای `Describe` و `Collect` با یک امضای (signature) مشخص باشند.

متد `Describe` توضیحی از متریک‌هایی که تولید خواهد کرد، به ویژه نام متریک، نام‌های برچسب‌ها و رشته راهنما (help string) را باز می‌گرداند. متد `Describe` در زمان ثبت (registration) فراخوانی می‌شود و برای جلوگیری از ثبت متریک تکراری استفاده می‌شود.

دو نوع متریک وجود دارد که یک اکسپورتر می‌تواند داشته باشد: آنهایی که نام‌ها و برچسب‌هایشان از قبل مشخص است، و آنهایی که فقط در زمان اسکرِیپ مشخص می‌شوند. در این مثال، `consul_up` از قبل مشخص است، بنابراین می‌توانید `Desc` آن را یک بار با `NewDesc` ایجاد کرده و از طریق `Describe` ارائه دهید. تمام متریک‌های دیگر به صورت پویا در زمان اسکرِیپ تولید می‌شوند، بنابراین نمی‌توانند در `Describe` گنجانده شوند:

```go
var (
	up = prometheus.NewDesc(
	"consul_up",
	"Was talking to Consul successful.",
	nil, nil,
	)
	)
	// اینترفیس prometheus.Collector را پیاده‌سازی می‌کند.
	func (c ConsulCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up // ارسال Desc برای متریک consul_up
	}
```

> کلاینت Go ایجاب می‌کند که حداقل یک `Desc` توسط `Describe` ارائه شود. اگر تمام متریک‌های شما پویا هستند، می‌توانید یک `Desc` ساختگی (dummy) برای رفع این مشکل ارائه دهید.

هسته یک کالکتور سفارشی، متد `Collect` است. در این متد، شما تمام داده‌های مورد نیاز خود را از نمونه برنامه‌ای که با آن کار می‌کنید واکشی (fetch) می‌کنید، در صورت نیاز آن‌ها را تغییر شکل (munge) می‌دهید و سپس متریک‌ها را به کتابخانه کلاینت باز می‌گردانید. در اینجا شما باید به Consul متصل شوید و سپس متریک‌های آن را واکشی کنید. اگر خطایی رخ دهد، `consul_up` به عنوان 0 بازگردانده می‌شود؛ در غیر این صورت، هنگامی که می‌دانیم جمع‌آوری موفقیت‌آمیز خواهد بود، به عنوان 1 بازگردانده می‌شود. بازگرداندن یک متریک فقط گاهی اوقات، کار با آن در PromQL دشوار است⁴؛ داشتن `consul_up` به شما امکان می‌دهد تا در مورد مشکلات ارتباط با Consul هشدار دهید تا بدانید که مشکلی وجود دارد.

برای بازگرداندن `consul_up`، از `prometheus.MustNewConstMetric` برای ارائه یک نمونه (sample) فقط برای این اسکرِیپ استفاده می‌شود. این تابع `Desc`، نوع و مقدار خود را می‌گیرد:

```go
// اینترفیس prometheus.Collector را پیاده‌سازی می‌کند.
func (c ConsulCollector) Collect(ch chan<- prometheus.Metric) {
consul, err := api.NewClient(api.DefaultConfig()) // اتصال به Consul
if err != nil {
// ارسال متریک consul_up با مقدار 0 در صورت خطا در اتصال
ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
return
}
metrics, err := consul.Agent().Metrics() // واکشی متریک‌ها
if err != nil {
// ارسال متریک consul_up با مقدار 0 در صورت خطا در واکشی
ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
return
}
// ارسال متریک consul_up با مقدار 1 نشان‌دهنده موفقیت
ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)
// ... پردازش سایر متریک‌ها ...
}
```

---
⁴ برای مثال، اگر فقط در صورت موفقیت `consul_up` را برگردانید، ممکن است در صورت عدم موفقیت، آخرین مقدار خوب قدیمی را در نمودارها ببینید.
---

سه مقدار ممکن وجود دارد: `GaugeValue`، `CounterValue` و `UntypedValue`.
`Gauge` و `Counter` را قبلاً می‌شناسید، و `Untyped` برای مواردی است که مطمئن نیستید یک متریک شمارنده است یا گیج. این با ابزار دقیق مستقیم ممکن نیست، اما غیرمعمول نیست که نوع متریک‌ها از سایر سیستم‌های مانیتورینگ و ابزار دقیق نامشخص و تعیین آن غیرعملی باشد.

اکنون که متریک‌ها را از Consul دارید، می‌توانید گیج‌ها را پردازش کنید. کاراکترهای نامعتبر در نام متریک، مانند نقطه و خط تیره، به زیرخط (_) تبدیل می‌شوند. یک `Desc` در لحظه ایجاد می‌شود و بلافاصله در `MustNewConstMetric` استفاده می‌شود:

```go
for _, g := range metrics.Gauges {
	// پاکسازی نام متریک Consul (مثلاً consul.autopilot.failure_tolerance -> consul_autopilot_failure_tolerance)
	name := invalidChars.ReplaceAllLiteralString(g.Name, "_")
	// ایجاد یک Desc جدید برای این متریک گیج خاص
	// ما nil را برای نام‌های برچسب متغیر ارسال می‌کنیم زیرا این گیج‌ها در این مثال برچسب متغیر ندارند.
	// g.Labels شامل برچسب‌های ثابتی است که از Consul می‌آیند.
	desc := prometheus.NewDesc("consul_"+name, "Consul metric "+g.Name, nil, g.Labels)
	// ایجاد و ارسال متریک ثابت با استفاده از Desc، نوع GaugeValue و مقدار گیج
	ch <- prometheus.MustNewConstMetric(
		desc, prometheus.GaugeValue, float64(g.Value))
		}
```

> ما `g.Labels` را به عنوان آخرین پارامتر به `prometheus.NewDesc` ارسال می‌کنیم. این مجموعه‌ای از برچسب‌هایی است که توسط Consul تنظیم شده‌اند، مانند برچسب مرکز داده (datacenter). ما باید آن‌ها را ارسال کنیم زیرا برخی از گیج‌ها کاردینالیتی (cardinality) بزرگتر از یک دارند و بدون آن برچسب‌ها، صفحه `/metrics` با خطا مواجه می‌شود.

پردازش شمارنده‌ها مشابه است، با این تفاوت که پسوند `_total` به نام متریک اضافه می‌شود:

```go
for _, c := range metrics.Counters {
	name := invalidChars.ReplaceAllLiteralString(c.Name, "_")
	// ایجاد Desc با نام پاکسازی شده + "_total"
	desc := prometheus.NewDesc("consul_"+name+"_total", "Consul metric "+c.Name, nil,
	c.Labels)
	// ارسال متریک ثابت با نوع CounterValue و مقدار Count شمارنده
	ch <- prometheus.MustNewConstMetric(
	desc, prometheus.CounterValue, float64(c.Count)) // توجه: اینجا باید c.Count باشد نه s.Count
	}
```

محتویات `metrics.Samples` پیچیده‌تر است. در حالی که نمونه‌ها یک خلاصه (Summary) پرومتئوس هستند، کلاینت Go در حال حاضر از آن‌ها برای `MustNewConstMetric` پشتیبانی نمی‌کند. در عوض، می‌توانید آن را با استفاده از دو شمارنده شبیه‌سازی کنید. `_seconds` به نام متریک اضافه می‌شود و مجموع (sum) بر هزار تقسیم می‌شود تا از میلی‌ثانیه به ثانیه تبدیل شود:

```go
for _, s := range metrics.Samples {
	// همه نمونه‌ها زمان بر حسب میلی‌ثانیه هستند، آن‌ها را به ثانیه تبدیل می‌کنیم.
	name := invalidChars.ReplaceAllLiteralString(s.Name, "_") + "_seconds"
	// ایجاد Desc برای بخش _count خلاصه شبیه‌سازی شده
	countDesc := prometheus.NewDesc(
	"consul_"+name+"_count", "Consul metric "+s.Name, nil, s.Labels)
	// ارسال متریک شمارنده برای _count
	ch <- prometheus.MustNewConstMetric(
	countDesc, prometheus.CounterValue, float64(s.Count))
	// ایجاد Desc برای بخش _sum خلاصه شبیه‌سازی شده
	sumDesc := prometheus.NewDesc(
	"consul_"+name+"_sum", "Consul metric "+s.Name, nil, s.Labels)
	// ارسال متریک شمارنده برای _sum، با تبدیل مقدار به ثانیه
	ch <- prometheus.MustNewConstMetric(
	sumDesc, prometheus.CounterValue, s.Sum/1000)
	}
```

> `s.Sum` در اینجا یک `float64` است، اما هنگام انجام تقسیم با اعداد صحیح باید مراقب باشید تا دقت را بی‌جهت از دست ندهید. اگر `sum` یک عدد صحیح بود، `float64(sum)/1000` ابتدا به ممیز شناور تبدیل می‌شد و سپس تقسیم می‌کرد، که همان چیزی است که می‌خواهید. از طرف دیگر، `float64(sum/1000)` ابتدا مقدار صحیح را بر هزار تقسیم می‌کند و سه رقم دقت را از دست می‌دهد.

در نهایت، شی کالکتور سفارشی نمونه‌سازی (instantiate) شده و با رجیستری پیش‌فرض ثبت می‌شود، به همان روشی که یکی از متریک‌های ابزار دقیق مستقیم را ثبت می‌کنید:

```go
c := ConsulCollector{} // ایجاد نمونه
prometheus.MustRegister(c) // ثبت با رجیستری پیش‌فرض
```

افشا (Exposition) به روش معمول انجام می‌شود، که قبلاً در «Go» در صفحه ۷۱ دیده‌اید:

```go
http.Handle("/metrics", promhttp.Handler()) // هندلر برای /metrics
log.Fatal(http.ListenAndServe(":8000", nil)) // راه‌اندازی سرور HTTP
```

این، البته، یک مثال ساده شده است. در واقعیت، شما روشی برای پیکربندی سرور Consul برای ارتباط با آن خواهید داشت، مانند یک فلگ (flag) خط فرمان، به جای وابستگی به پیش‌فرض کلاینت. همچنین کلاینت را بین اسکرِیپ‌ها مجدداً استفاده می‌کنید و اجازه می‌دهید گزینه‌های مختلف احراز هویت کلاینت مشخص شوند.

> مقادیر `min`، `max`، `mean` و `stddev` از خروجی اصلی حذف شدند زیرا چندان مفید نیستند. شما می‌توانید `mean` را با استفاده از `sum` و `count` محاسبه کنید. از طرف دیگر، `min`، `max` و `stddev` قابل تجمیع (aggregate) نیستند و شما نمی‌دانید در چه بازه زمانی اندازه‌گیری شده‌اند.

از آنجایی که از رجیستری پیش‌فرض استفاده می‌شود، متریک‌های `go_` و `process_` در نتیجه گنجانده می‌شوند. این‌ها اطلاعاتی در مورد عملکرد خود اکسپورتر به شما ارائه می‌دهند و برای تشخیص مشکلاتی مانند نشت توصیف‌گر فایل (file descriptor leaks) با استفاده از `process_open_fds` مفید هستند. این شما را از نیاز به اسکرِیپ جداگانه اکسپورتر برای این متریک‌ها نجات می‌دهد.

تنها زمانی که ممکن است از رجیستری پیش‌فرض برای یک اکسپورتر استفاده نکنید، هنگام نوشتن یک اکسپورتر به سبک Blackbox/SNMP است، جایی که نیاز به تفسیر پارامترهای URL وجود دارد زیرا کالکتورها به پارامترهای URL برای یک اسکرِیپ دسترسی ندارند. در آن صورت، شما همچنین `/metrics` خود اکسپورتر را اسکرِیپ می‌کنید تا خود اکسپورتر را مانیتور کنید.

برای مقایسه، اکسپورتر معادل نوشته شده با پایتون ۳ در مثال ۱۲-۳ نشان داده شده است. این تا حد زیادی شبیه به اکسپورتری است که به زبان Go نوشته شده است؛ تنها تفاوت قابل توجه این است که یک `SummaryMetricFamily` برای نمایش خلاصه در دسترس است، به جای شبیه‌سازی آن با دو شمارنده جداگانه. کلاینت پایتون به اندازه کلاینت Go بررسی‌های سلامت (sanity checks) ندارد، بنابراین باید کمی با آن محتاط‌تر باشید.

مثال ۱۲-۳. consul_metrics.py، یک اکسپورتر برای متریک‌های Consul که به زبان پایتون ۳ نوشته شده است

```python
import json
import re
import time
from urllib.request import urlopen
from prometheus_client.core import GaugeMetricFamily, CounterMetricFamily
from prometheus_client.core import SummaryMetricFamily, REGISTRY
from prometheus_client import start_http_server

# تابعی برای پاکسازی نام‌های متریک (جایگزینی کاراکترهای نامعتبر با _)
def sanitize_name(s):
    return re.sub(r"[^a-zA-Z0-9:_]", "_", s)

class ConsulCollector(object):
    def collect(self):
        # این متد در هر بار اسکرِیپ فراخوانی می‌شود
        # ۱. واکشی داده‌ها از Consul
        try:
            out = urlopen("http://localhost:8500/v1/agent/metrics").read()
            metrics = json.loads(out.decode("utf-8"))
            # متریک consul_up برای نشان دادن موفقیت واکشی
            yield GaugeMetricFamily("consul_up", "Was fetching Consul metrics successful", value=1)
        except Exception as e:
            # در صورت خطا، consul_up را 0 برگردانید
            yield GaugeMetricFamily("consul_up", "Was fetching Consul metrics successful", value=0)
            print(f"Error fetching Consul metrics: {e}")
            return # از ادامه پردازش جلوگیری کنید

        # ۲. پردازش و تولید متریک‌ها
        # پردازش Gauges
        for g in metrics.get("Gauges", []):
            # ایجاد یک خانواده متریک گیج
            # توجه: برچسب‌ها در اینجا به صورت مستقیم پشتیبانی نمی‌شوند،
            # اگر برچسب وجود داشت، باید خانواده متریک را با labels=[] تعریف می‌کردید
            # و از .add_metric(label_values, value) استفاده می‌کردید.
            yield GaugeMetricFamily("consul_" + sanitize_name(g["Name"]),
                                    "Consul metric " + g["Name"],
                                    value=g["Value"])
        # پردازش Counters
        for c in metrics.get("Counters", []):
            # ایجاد یک خانواده متریک شمارنده با پسوند _total
            yield CounterMetricFamily("consul_" + sanitize_name(c["Name"]) + "_total",
                                      "Consul metric " + c["Name"],
                                      value=c["Count"])
        # پردازش Samples (به عنوان Summary)
        for s in metrics.get("Samples", []):
            # ایجاد یک خانواده متریک خلاصه
            # نام شامل _seconds و مقدار sum به ثانیه تبدیل می‌شود
            yield SummaryMetricFamily("consul_" + sanitize_name(s["Name"]) + "_seconds",
                                      "Consul metric " + s["Name"],
                                      count_value=s["Count"], # اینجا باید از s استفاده شود نه c
                                      sum_value=s["Sum"] / 1000)

if __name__ == '__main__':
    # ۳. ثبت کالکتور و راه‌اندازی سرور HTTP
    REGISTRY.register(ConsulCollector()) # ثبت کالکتور
    start_http_server(8000) # راه‌اندازی سرور HTTP در پورت 8000
    print("Starting Consul exporter (Python) on :8000")
    while True:
        time.sleep(1) # سرور را فعال نگه دارید
```

### برچسب‌ها (Labels)

در مثال قبلی، شما فقط متریک‌های بدون برچسب را دیدید. برای ارائه برچسب‌ها، باید نام‌های برچسب را در `Desc` (در Go) یا هنگام ایجاد `MetricFamily` (در پایتون) مشخص کنید و سپس مقادیر را در `MustNewConstMetric` (در Go) یا با `add_metric` (در پایتون) ارائه دهید.

برای افشای یک متریک با سری‌های زمانی `example_gauge{foo="bar", baz="small"}` و `example_gauge{foo="quu", baz="far"}`, با استفاده از کتابخانه کلاینت پرومتئوس Go می‌توانید این کار را انجام دهید:

```go
func (c MyCollector) Collect(ch chan<- prometheus.Metric) {
	// تعریف Desc با نام متریک، راهنما و نام‌های برچسب
	desc := prometheus.NewDesc(
	"example_gauge", // نام متریک
	"A help string.", // راهنما
	[]string{"foo", "baz"}, // نام‌های برچسب
	nil, // برچسب‌های ثابت (در اینجا استفاده نشده)
	)
	// ارسال اولین سری زمانی با مقادیر برچسب "bar" و "small"
	ch <- prometheus.MustNewConstMetric(
	desc, prometheus.GaugeValue, 1, "bar", "small")
	// ارسال دومین سری زمانی با مقادیر برچسب "quu" و "far"
	ch <- prometheus.MustNewConstMetric(
	desc, prometheus.GaugeValue, 2, "quu", "far")
	}
```

ابتدا، شما می‌توانید هر سری زمانی را به صورت جداگانه ارائه دهید. رجیستری مسئول ترکیب تمام سری‌های زمانی متعلق به یک خانواده متریک در خروجی `/metrics` خواهد بود.

> رشته‌های راهنمای تمام متریک‌ها با نام یکسان باید یکسان باشند. ارائه `Desc`های متفاوت باعث شکست اسکرِیپ خواهد شد.

کلاینت پایتون کمی متفاوت عمل می‌کند؛ شما خانواده متریک را مونتاژ کرده و سپس آن را باز می‌گردانید (`yield`). اگرچه ممکن است به نظر تلاش بیشتری نیاز داشته باشد، اما معمولاً در عمل به همان سطح تلاش منجر می‌شود:

```python
from prometheus_client.core import GaugeMetricFamily

class MyCollector(object):
    def collect(self):
        # ایجاد یک خانواده متریک گیج با نام، راهنما و لیست نام‌های برچسب
        mf = GaugeMetricFamily("example_gauge", "A help string.",
                               labels=["foo", "baz"])
        # اضافه کردن اولین سری زمانی با مقادیر برچسب و مقدار متریک
        mf.add_metric(["bar", "small"], 1)
        # اضافه کردن دومین سری زمانی
        mf.add_metric(["quu", "far"], 2)
        # بازگرداندن (yield) خانواده متریک مونتاژ شده
        yield mf
```

## دستورالعمل‌ها (Guidelines)

درحالی‌که ابزار دقیق مستقیم معمولاً نسبتاً ساده است، نوشتن اکسپورترها معمولاً مبهم‌تر است و شامل بده‌بستان‌های مهندسی (engineering trade-offs) می‌شود. آیا می‌خواهید تلاش مداوم زیادی برای تولید متریک‌های عالی صرف کنید، یا کاری انجام دهید که به اندازه کافی خوب باشد و نیازی به نگهداری نداشته باشد؟ نوشتن اکسپورترها بیشتر یک هنر است تا یک علم.

شما باید سعی کنید از شیوه‌های نام‌گذاری متریک پیروی کنید، به ویژه، از پسوندهای `_count`، `_sum`، `_total`، `_bucket` و `_info` اجتناب کنید مگر اینکه سری زمانی بخشی از متریکی باشد که قرار است شامل چنین سری زمانی باشد.

اغلب ممکن یا عملی نیست که تعیین کنید آیا مجموعه‌ای از متریک‌ها گیج، شمارنده یا ترکیبی از هر دو هستند. در مواردی که ترکیبی وجود دارد، باید آن‌ها را به عنوان بدون نوع (untyped) علامت‌گذاری کنید به جای استفاده از گیج یا شمارنده که نادرست خواهد بود. اگر یک متریک شمارنده است، فراموش نکنید که پسوند `_total` را اضافه کنید.

در صورت امکان عملی، باید سعی کنید واحدهایی برای متریک‌های خود ارائه دهید و حداقل سعی کنید اطمینان حاصل کنید که واحدها در نام متریک وجود دارند. مجبور بودن به تعیین واحدها از روی متریک‌ها، مانند مثال ۱۲-۱، برای هیچ‌کس جالب نیست، بنابراین باید سعی کنید این بار را از دوش کاربران اکسپورتر خود بردارید. ثانیه و بایت همیشه ترجیح داده می‌شوند.

از نظر استفاده از برچسب‌ها در اکسپورترها، چند نکته مهم وجود دارد که باید مراقب آن‌ها بود. همانند ابزار دقیق مستقیم، کاردینالیتی (cardinality) نیز به همان دلایلی که در «کاردینالیتی» در صفحه ۹۹ بحث شد، برای اکسپورترها نگران‌کننده است. باید از متریک‌هایی با تغییرات زیاد (high churn) در برچسب‌هایشان اجتناب شود.

برچسب‌ها باید یک افراز (partition) در سراسر یک متریک ایجاد کنند، و اگر شما یک مجموع (sum) یا میانگین (average) در سراسر یک متریک بگیرید، باید معنی‌دار باشد، همانطور که در «چه زمانی از برچسب‌ها استفاده کنیم» در صفحه ۹۸ بحث شد. به طور خاص، باید مراقب هر سری زمانی باشید که فقط مجموع تمام مقادیر دیگر در یک متریک است و آن‌ها را حذف کنید. اگر هنگام نوشتن یک اکسپورتر در مورد اینکه آیا یک برچسب معنی‌دار است شک دارید، امن‌ترین کار این است که از آن استفاده نکنید، اگرچه بحث «استثنای جدول» در صفحه ۹۹ را در نظر داشته باشید. همانند ابزار دقیق مستقیم، نباید برچسبی مانند `env="prod"` را به تمام متریک‌هایی که از اکسپورتر شما می‌آیند اعمال کنید، زیرا این کاری است که برچسب‌های هدف (target labels) برای آن هستند، همانطور که در «برچسب‌های هدف» در صفحه ۱۵۳ بحث شد.

بهتر است متریک‌های خام (raw) را به پرومتئوس افشا کنید، به جای انجام محاسبات در سمت برنامه. به عنوان مثال، نیازی به افشای نرخ ۵ دقیقه‌ای (5-minute rate) نیست وقتی که یک شمارنده دارید، زیرا می‌توانید از تابع `rate` برای محاسبه نرخ در هر دوره‌ای که دوست دارید استفاده کنید. به همین ترتیب در مورد نسبت‌ها (ratios)، آن‌ها را به نفع شمارنده (numerator) و مخرج (denominator) حذف کنید. اگر درصدی (percentage) بدون شمارنده و مخرج تشکیل‌دهنده‌اش دارید، حداقل آن را به یک نسبت تبدیل کنید.⁵

فراتر از ضرب و تقسیم برای استانداردسازی واحدها، باید از ریاضیات در اکسپورترها اجتناب کنید، زیرا پردازش داده‌های خام در PromQL ترجیح داده می‌شود. شرایط رقابتی (Race conditions) بین رویدادهای ابزار دقیق متریک‌ها می‌تواند منجر به مصنوعات (artifacts) شود، به خصوص زمانی که یک متریک را از دیگری کم می‌کنید. جمع کردن متریک‌ها به منظور کاهش کاردینالیتی می‌تواند خوب باشد، اما اگر شمارنده هستند، مطمئن شوید که به دلیل ناپدید شدن برخی از آن‌ها، ریست‌های کاذب (spurious resets) وجود نخواهد داشت.

برخی متریک‌ها با توجه به نحوه استفاده مورد نظر از پرومتئوس، چندان مفید نیستند. بسیاری از برنامه‌ها متریک‌هایی مانند RAM ماشین، CPU و دیسک را افشا می‌کنند. شما نباید متریک‌های سطح ماشین (machine-level) را در اکسپورتر خود افشا کنید، زیرا این مسئولیت Node Exporter است.⁶ مقادیر حداقل (Minimums)، حداکثر (maximums) و انحراف معیار (standard deviations) را نمی‌توان به طور منطقی تجمیع کرد، بنابراین باید آن‌ها نیز حذف شوند.

شما باید برنامه‌ریزی کنید که یک اکسپورتر به ازای هر نمونه برنامه (application instance) اجرا کنید⁷، و متریک‌ها را به صورت همگام (synchronously) برای هر اسکرِیپ بدون هیچ‌گونه کَش کردن (caching) واکشی کنید. این کار مسئولیت‌های کشف سرویس (service discovery) و زمان‌بندی اسکرِیپ را با پرومتئوس نگه می‌دارد. توجه داشته باشید که باید آگاه باشید که اسکرِیپ‌های همزمان (concurrent scrapes) می‌توانند اتفاق بیفتند.

درست همانطور که پرومتئوس هنگام انجام اسکرِیپ، متریک `scrape_duration_seconds` را اضافه می‌کند، شما نیز ممکن است متریک `myexporter_scrape_duration_seconds` را برای مدت زمانی که اکسپورتر شما طول می‌کشد تا داده‌ها را از برنامه‌اش بیرون بکشد، اضافه کنید. این به اشکال‌زدایی عملکرد (performance debugging) کمک می‌کند، زیرا می‌توانید ببینید که آیا برنامه یا اکسپورتر شما کند شده است. متریک‌های اضافی مانند تعداد متریک‌های پردازش شده نیز می‌توانند مفید باشند.

ممکن است منطقی باشد که ابزار دقیق مستقیم را به اکسپورترها اضافه کنید، علاوه بر کالکتورهای سفارشی که عملکرد اصلی آن‌ها را ارائه می‌دهند. به عنوان مثال، CloudWatch Exporter دارای یک شمارنده `cloudwatch_requests_total` است که تعداد تماس‌های API را که برقرار می‌کند ردیابی می‌کند، زیرا هر تماس API هزینه دارد. اما این معمولاً فقط چیزی است که با اکسپورترهای سبک Blackbox/SNMP مشاهده خواهید کرد.

اکنون که می‌دانید چگونه متریک‌ها را هم از برنامه‌های خود و هم از کدهای شخص ثالث (third-party) دریافت کنید، در فصل بعد شروع به پوشش PromQL خواهیم کرد، که به شما امکان می‌دهد با این متریک‌ها کار کنید.

---
⁵ نسبت‌ها بین ۰ و ۱ هستند و راحت‌تر می‌توان با آن‌ها کار کرد.
⁶ مگر اینکه اکسپورتر شما Node Exporter باشد!
⁷ استثنا اکسپورترهای سبک Blackbox یا SNMP هستند، که ذاتاً به صورت پروکسی عمل می‌کنند.
---
```