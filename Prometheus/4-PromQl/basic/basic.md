


‏`PromQL` زبان `Query Language` مربوط به `Prometheus` است. با اینکه به QL ختم می‌شود، متوجه خواهید شد که زبانی شبیه `SQL` نیست، زیرا زبان‌های `SQL` معمولاً قدرت بیانی کافی برای انجام محاسباتی که روی `time series` (سری‌های زمانی) انجام می‌دهید را ندارند. `Labels` (برچسب‌ها) بخش کلیدی `PromQL` هستند و می‌توانید از آن‌ها نه تنها برای انجام `aggregations` (تجمیع‌های) دلخواه، بلکه برای پیوند دادن `metrics` مختلف به یکدیگر جهت انجام `arithmetic operations` (عملیات حسابی) روی آن‌ها استفاده کنید. طیف گسترده‌ای از توابع از `prediction` (پیش‌بینی) گرفته تا توابع `date` (تاریخ) و `math functions` (ریاضی) در دسترس شما هستند.

این درس به شما را با مفاهیم پایه `PromQL`، شامل `aggregation` (تجمیع)، `basic types` (انواع پایه) و `HTTP API` آشنا می‌کند.

‏`Prometheus` یک زبان `query` تابعی به نام `PromQL` (`Prometheus Query Language`) ارائه می‌دهد که به کاربر اجازه می‌دهد داده‌های `time series` را در زمان واقعی انتخاب و `aggregate` کند.

وقتی یک درخواست `query` به `Prometheus` ارسال می‌کنید، می‌تواند یک `instant query` باشد که در یک نقطه زمانی ارزیابی می‌شود، یا یک `range query` در گام‌های با فاصله مساوی بین زمان شروع و پایان. `PromQL` دقیقاً به همین شکل در هر دو حالت کار می‌کند؛ `range query` مانند یک `instant query` است که چندین بار در `timestamp`های مختلف اجرا می‌شود.

در `UI` `Prometheus`، تب `"Table"` برای `instant queries` و تب `"Graph"` برای `range queries` است.

برنامه‌های دیگر می‌توانند نتیجه یک `expression` `PromQL` را از طریق `HTTP API` دریافت کنند.

**مثال‌ها**
این سند یک مرجع پایه زبان `Prometheus` است. برای یادگیری، ممکن است شروع با چند مثال آسان‌تر باشد.

**انواع داده زبان `Expression`**
در زبان `expression` `Prometheus`، یک `expression` یا `sub-expression` می‌تواند به یکی از چهار نوع ارزیابی شود:

*   ‏`Instant vector` - مجموعه‌ای از `time series` حاوی یک `sample` برای هر `time series`، که همگی `timestamp` یکسانی دارند.
*   ‏`Range vector` - مجموعه‌ای از `time series` حاوی محدوده‌ای از نقاط داده در طول زمان برای هر `time series`.
*   ‏`Scalar` - یک مقدار عددی ساده `floating point`.
*   ‏`String` - یک مقدار `string` ساده؛ در حال حاضر استفاده نمی‌شود.



## Instant Vector

یک `instant vector selector` یک `instant vector` از جدیدترین `samples` قبل از زمان ارزیابی کوئری (`query evaluation time`) را برمی‌گرداند، یعنی لیستی از صفر یا چند `time series`. هر یک از این `time series`ها یک `sample` خواهند داشت و یک `sample` هم مقدار و هم `timestamp` را در بر می‌گیرد. در حالی که `instant vector` بازگشتی توسط یک `instant vector selector` دارای `timestamp` داده اصلی است، هر `instant vector` بازگشتی توسط عملیات یا توابع دیگر، `timestamp` زمان ارزیابی کوئری را برای تمام مقادیر خود خواهد داشت.

هنگامی که استفاده فعلی حافظه را درخواست می‌کنید، نمی‌خواهید `samples` از یک `instance` که روزها پیش خاموش شده است، شامل شود، مفهومی که به عنوان `staleness` (کهنگی) شناخته می‌شود. در Prometheus 1.x این کار با برگرداندن `time series`هایی انجام می‌شد که `sample` آن‌ها حداکثر ۵ دقیقه قبل از زمان ارزیابی کوئری بود. این تا حد زیادی کار می‌کرد اما معایبی مانند شمارش مضاعف داشت اگر یک `instance` با `label` `instance` جدید در آن پنجره ۵ دقیقه‌ای راه‌اندازی مجدد می‌شد.

‏Prometheus 2.x رویکرد پیچیده‌تری دارد. اگر یک `time series` از یک `scrape` به `scrape` بعدی ناپدید شود، یا اگر یک `target` دیگر از `service discovery` بازگردانده نشود، نوع خاصی از `sample` به نام `stale marker`¹⁰ (نشانگر کهنه) به `time series` اضافه می‌شود. هنگام ارزیابی یک `instant vector selector`، ابتدا تمام `time series`هایی که تمام `matchers` را برآورده می‌کنند پیدا می‌شوند و جدیدترین `sample` در ۵ دقیقه قبل از زمان ارزیابی کوئری همچنان در نظر گرفته می‌شود. اگر `sample` یک `sample` معمولی باشد، در `instant vector` بازگردانده می‌شود، اما اگر یک `stale marker` باشد، آن `time series` در آن `instant vector` گنجانده نخواهد شد. نتیجه همه اینها این است که وقتی از `instant vector selector` استفاده می‌کنید، `time series`هایی که `stale` (کهنه) شده‌اند، بازگردانده نمی‌شوند.

اگر `exporter` ای دارید که `timestamps` را نمایش می‌دهد، همانطور که در "Timestamps" در قبل توضیح داده شده است، `stale markers` و منطق `staleness` Prometheus 2.x اعمال نخواهد شد. `time series`های تحت تأثیر به جای آن با منطق قدیمی‌تر که ۵ دقیقه به عقب نگاه می‌کند، کار خواهند کرد.


```json
curl 'http://localhost:9090/api/v1/query' \
  --data 'query=http_requests_total{code="200"}' \
  --data time=1608481001
{
  "metric": {"__name__": "http_requests_total", "code": "200"},
  "value": [1608481001, "881"]
}
```

![[instantVec2.png]]



![[instantVector.png]]





### Range Vector

نوع دوم `selector` وجود دارد که قبلاً دیده‌اید و `range vector selector` نامیده می‌شود. برخلاف `instant vector selector` که یک `sample` برای هر `time series` برمی‌گرداند، یک `range vector selector` می‌تواند چندین `sample` برای هر `time series` برگرداند.¹¹ `Range vectors` همیشه با تابع `rate` استفاده می‌شوند، به عنوان مثال:

`rate(process_cpu_seconds_total[1m])`

‏`[1m]` `instant vector selector` را به `range vector selector` تبدیل می‌کند و به `PromQL` دستور می‌دهد تا برای تمام `time series`های مطابق با `selector`، تمام `samples` مربوط به دقیقه منتهی به زمان ارزیابی کوئری را برگرداند. اگر فقط `process_cpu_seconds_total[1m]` را در تب `Console` مرورگر `expression` اجرا کنید، چیزی شبیه شکل زیر خواهید دید.

![[pic13.png]]

در این مورد، هر `time series` اتفاقاً شش `sample` در دقیقه گذشته دارد. متوجه خواهید شد که در حالی که `samples` برای هر `time series` اتفاقاً دقیقاً ۱۰ ثانیه از هم فاصله دارند مطابق با فاصله `scrape` که پیکربندی کرده‌اید، `timestamps` دو `time series` با یکدیگر هم‌تراز نیستند. یک `time series` یک `sample` با `timestamp` 1517925155.087 دارد و دیگری 1517925156.245.



این به این دلیل است که `range vectors` `timestamps` واقعی `samples` را حفظ می‌کنند و `scrapes` برای `targets` مختلف به منظور توزیع بار به طور مساوی‌تر، توزیع می‌شوند. در حالی که می‌توانید فرکانس `scrapes` و ارزیابی قوانین را کنترل کنید، نمی‌توانید `phase` یا `alignment` (ترازبندی) آن‌ها را کنترل کنید. اگر فاصله `scrape` ۱۰ ثانیه‌ای و صدها `target` دارید، تمام آن `target`ها در نقاط مختلف در یک پنجره ۱۰ ثانیه‌ای معین `scrape` خواهند شد. به عبارت دیگر، `time series`های شما همگی سن‌های کمی متفاوت دارند. این به طور کلی در عمل برای شما اهمیتی نخواهد داشت، اما می‌تواند منجر به `artifacts` (مصنوعات) شود زیرا اساساً سیستم‌های مانیتورینگ مبتنی بر `metrics` مانند `Prometheus` تخمین‌های (بسیار خوبی) به جای پاسخ‌های دقیق تولید می‌کنند.

شما به ندرت به طور مستقیم به `range vectors` نگاه خواهید کرد. این فقط زمانی پیش می‌آید که برای `debugging` نیاز به دیدن `samples` خام دارید. تقریباً همیشه از `range vector` با تابعی مانند `rate` یا `avg_over_time` استفاده خواهید کرد که `range vector` را به عنوان آرگومان می‌گیرد. `Staleness` و `stale markers` تأثیری بر `range vectors` ندارند؛ شما تمام `samples` معمولی را در یک بازه معین دریافت خواهید کرد. هر `stale marker` دیگری نیز که در آن بازه باشد، توسط `range vector selector` بازگردانده نمی‌شود.

![[ReangeVec.png]]


![[RangeVector2.png]]



```json
curl 'http://localhost:9090/api/v1/query' \
  --data 'query=http_requests_total{code="200"}[30s]' \
  --data time=1608481001
{
  "metric": {"__name__": "http_requests_total", "code": "200"},
  "values": [
    [1608480978, "863"],
    [1608480986, "874"],
    [1608480094, "881"]
  ]
}
```

بسته به مورد استفاده (مثلاً هنگام ترسیم `graph` در مقابل نمایش خروجی یک `expression`)، فقط برخی از این انواع به عنوان نتیجه یک `expression` مشخص شده توسط کاربر مجاز هستند. برای `instant queries`، هر یک از انواع داده بالا به عنوان ریشه `expression` مجاز است. `Range queries` فقط از `expression`های نوع `scalar` و `instant-vector` پشتیبانی می‌کنند.

**یادداشت‌هایی درباره `native histograms` تجربی:**

*   فعال کردن `Ingesting` `native histograms` باید از طریق یک `feature flag` انجام شود.
*   هنگامی که `native histograms` در `TSDB` `ingest` شدند (و حتی پس از غیرفعال کردن دوباره `feature flag`)، هم `instant vectors` و هم `range vectors` ممکن است اکنون حاوی `sample`هایی باشند که اعداد ساده `floating point` (`float samples`) نیستند، بلکه `histograms` کامل (`histogram samples`) هستند. یک `vector` ممکن است ترکیبی از `float samples` و `histogram samples` را شامل شود. توجه داشته باشید که اصطلاح «`histogram sample`» در مستندات `PromQL` همیشه به یک `native histogram` اشاره دارد. `Classic histograms` به تعدادی `series` از `float samples` تقسیم می‌شوند. از دیدگاه `PromQL`، چیزی به نام «`classic histogram samples`» وجود ندارد.
*   مانند `float samples`، `histogram samples` می‌توانند `counters` یا `gauges` باشند که به ترتیب `counter histograms` یا `gauge histograms` نیز نامیده می‌شوند.
*  ‏ `Native histograms` می‌توانند `bucket layout`های مختلفی داشته باشند، اما به طور کلی برای اعمال عملیات `binary` و `aggregation` روی آنها، قابل تبدیل به نسخه‌های سازگار هستند. این برای همه `bucketing schemas` صادق نیست. اگر `histograms` ناسازگار در یک عملیات مواجه شوند، عنصر `vector` خروجی مربوطه از نتیجه حذف شده و با یک `annotation` سطح هشدار (warn-level) علامت‌گذاری می‌شود. جزئیات بیشتر را می‌توان در `native histogram specification` یافت.

عبارت **Native Histograms** یکی از قابلیت‌های جدید در **Prometheus** است که برای بهبود عملکرد و دقت در ذخیره‌سازی و پردازش **توزیع مقادیر (distribution of values)** به کار می‌رود.



---

## 📦 Native Histograms در Prometheus چیست؟

‏**Native Histograms** نوع جدیدی از داده‌ها در Prometheus هستند که برای ذخیره‌سازی توزیع مقادیر به‌صورت **فشرده و دقیق‌تر** طراحی شده‌اند، به‌ویژه برای داده‌هایی مثل latency، duration یا اندازه‌ها.

آن‌ها جایگزینی یا مکملی برای **classic histograms** (`*_bucket` metrics) هستند.

---

## 📉 تفاوت با Classic Histograms

### 🎯 Classic Histogram (قدیمی):

* از متریک‌هایی مثل `http_request_duration_seconds_bucket{le="0.5"}` استفاده می‌کند.
* برای هر bucket یک متریک جدا دارد.
* دقت پایین‌تر و حجم ذخیره‌سازی بیشتر.
* محاسبات روی آن‌ها پیچیده‌تر و پرهزینه‌تر است.

### ⚡ Native Histogram:

* به صورت **یک متریک واحد** ذخیره می‌شود (نه چند `*_bucket` جدا).
* دقت بالا در محدوده وسیع.
* کارآمدتر در مصرف حافظه و فضای ذخیره‌سازی.
* از لحاظ عملکرد برای query بسیار سریع‌تر است.

---

## 🔬 چطور کار می‌کند؟

* از **logarithmic bucketing** استفاده می‌کند، یعنی فواصل بین bucketها بسته به مقدار تغییر می‌کنند.
* به طور **فشرده‌شده** و **باینری** ذخیره می‌شود.
* ‏Prometheus از نسخه **v2.40** به بعد از آن پشتیبانی می‌کند.

---

## 📘 نمونه متریک:

به جای داشتن چند متریک مثل:

```
http_request_duration_seconds_bucket{le="0.1"}
http_request_duration_seconds_bucket{le="0.5"}
...
```

در Native Histogram همه‌ی این‌ها در یک **متریک compact** ذخیره می‌شود مثل:

```
http_request_duration_seconds
```

و نوع آن به صورت `histogram` در متادیتا تعریف می‌شود.

---

## ✅ مزایا:

| ویژگی                           | توضیح                                                |
| ------------------------------- | ---------------------------------------------------- |
| 🎯 دقت بالا                     | حتی در بازه‌های بسیار کوچک یا بزرگ                   |
| 📦 مصرف فضای کمتر               | چون فشرده است                                        |
| 🚀 کارایی بالا در Query         | مخصوصاً در queryهای طولانی مدت یا aggregate          |
| 📈 مناسب برای long-term storage | با TSDBهایی مثل Thanos و Cortex هم سازگارتر خواهد شد |

---

## 🔧 فعال‌سازی:

برای فعال‌سازی Native Histograms در یک exporter (مثلاً Go client):

```go
prometheus.NewHistogram(prometheus.HistogramOpts{
    Name:                   "http_request_duration_seconds",
    NativeHistogramBucketFactor: 1.1, // فعال‌سازی native histogram
})
```

> این قابلیت فعلاً **opt-in** است، یعنی باید به‌صورت دستی فعال شود.

---




**`Literals`**

در **PromQL** (زبان کوئری‌نویسی Prometheus)، **String Literals** به رشته‌هایی اطلاق می‌شود که برای فیلتر کردن، مقایسه، یا مقداردهی به **label**‌ها و برخی توابع استفاده می‌شوند.

---

## 🔤 تعریف String Literal در Prometheus:

‏**String Literal** در Prometheus یعنی **یک مقدار متنی محصور در کوتیشن (`" "` یا `' '`)** که در کوئری‌ها استفاده می‌شود — معمولاً برای **مقدار برچسب‌ها (labels)** یا ورودی برخی توابع.

---

## 📘 مثال‌های رایج:

### 1. **در فیلتر کردن Labelها**

```promql
http_requests_total{job="api-server"}
```

در اینجا `"api-server"` یک **string literal** است که برای فیلتر کردن مقدار برچسب `job` استفاده می‌شود.

---

### 2. **با regular expressions (regex)**

```promql
http_requests_total{instance=~"web-[0-9]+"}
```

در اینجا `"web-[0-9]+"` یک string literal برای یک **عبارت منظم (regex)** است.

---

### 3. **در توابعی مثل `label_replace()`**

```promql
label_replace(up, "new_label", "$1", "job", "(.*)")
```

در این تابع، `"new_label"`, `"$1"`, `"job"` و `"(.*)"` همگی string literal هستند.

---

## ✅ نکات مهم:

| ویژگی                                         | توضیح                                      |
| --------------------------------------------- | ------------------------------------------ |
| محصور در `" "` یا `' '`                       | معمولاً از `" "` استفاده می‌شود            |
| اغلب در فیلتر label استفاده می‌شود            | به شکل `label="value"` یا `label=~"regex"` |
| استفاده در توابع خاص                          | مثل `label_replace`, `label_join`, و ...   |
| فقط نوع **String** در labelها پشتیبانی می‌شود | یعنی labelها فقط رشته هستند                |

---

## ❌ غلط رایج:

```promql
up{job=api-server}   ❌ اشتباه: مقدار باید در کوتیشن باشد
```

✅ درست:

```promql
up{job="api-server"}
```

---



**`String literals`**
`‏String literals` با `single quotes`، `double quotes` یا `backticks` مشخص می‌شوند.

‏`PromQL` از همان `escaping rules` زبان `Go` پیروی می‌کند. برای `string literals` در `single` یا `double quotes`، یک `backslash` یک `escape sequence` را آغاز می‌کند، که ممکن است با `a`, `b`, `f`, `n`, `r`, `t`, `v` یا `\` دنبال شود. کاراکترهای خاص را می‌توان با استفاده از نمادگذاری `octal` (`\nnn`) یا `hexadecimal` (`\xnn`, `\unnnn` و `\Unnnnnnnn`) ارائه داد.

برعکس، کاراکترهای `escape` در `string literals` مشخص شده با `backticks` تجزیه نمی‌شوند. مهم است توجه داشته باشید که برخلاف `Go`، `Prometheus` `newlines` داخل `backticks` را حذف نمی‌کند.

مثال:

```
"this is a string"
'these are unescaped: \n \\ \t'
`these are not unescaped: \n ' " \t`
```

**`Float literals` و `time durations`**
مقادیر `Scalar` `float` را می‌توان به صورت اعداد `integer` یا `floating-point` `literal` در قالب زیر نوشت (فضای خالی فقط برای خوانایی بهتر گنجانده شده است):

```
[-+]?(
      [0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?
    | 0[xX][0-9a-fA-F]+
    | [nN][aA][nN]
    | [iI][nN][fF]
)
```

مثال‌ها:

```
23
-2.43
3.4e-9
0x8f
-Inf
NaN
```

علاوه بر این، `underscores` (_) را می‌توان بین ارقام `decimal` یا `hexadecimal` برای بهبود خوانایی استفاده کرد.

مثال‌ها:

```
1_000_000
.123_456_789
0x_53_AB_F3_82
```

‏`Float literals` همچنین برای مشخص کردن `durations` بر حسب ثانیه استفاده می‌شوند. برای راحتی، اعداد `integer` `decimal` ممکن است با `time units` زیر ترکیب شوند:

* ‏  `ms` – milliseconds
* ‏  `s` – seconds – 1s برابر است با 1000ms
* ‏  `m` – minutes – 1m برابر است با 60s (نادیده گرفتن leap seconds)
* ‏  `h` – hours – 1h برابر است با 60m
* ‏  `d` – days – 1d برابر است با 24h (نادیده گرفتن daylight saving time)
* ‏  `w` – weeks – 1w برابر است با 7d
* ‏  `y` – years – 1y برابر است با 365d (نادیده گرفتن leap days)

افزودن پسوند یکی از واحدهای بالا به یک عدد `integer` `decimal`، نمایش متفاوتی از معادل تعداد ثانیه به عنوان یک `float literal` خالی است.

مثال‌ها:

```
1s # معادل 1.
2m # معادل 120.
1ms # معادل 0.001.
-2h # معادل -7200.
```

مثال‌های زیر کار نمی‌کنند:

```
0xABm # عدم امکان پسوندگذاری اعداد hexadecimal.
1.5h # Time units را نمی‌توان با یک floating point ترکیب کرد.
+Infd # عدم امکان پسوندگذاری ±Inf یا NaN.
```

چندین واحد را می‌توان با `concatenation` اعداد صحیح پسونددار ترکیب کرد. واحدها باید از طولانی‌ترین به کوتاه‌ترین مرتب شوند. یک واحد داده شده فقط باید یک بار در هر `float literal` ظاهر شود.

مثال‌ها:

```
1h30m # معادل 5400s و بنابراین 5400.
12h34m56s # معادل 45296s و بنابراین 45296.
54s321ms # معادل 54.321.
```

**`Time series selectors`**
اینها بلوک‌های سازنده اصلی هستند که به `PromQL` دستور می‌دهند چه داده‌هایی را `fetch` کند.

### **`Instant vector selectors`**

در Prometheus و زبان کوئری PromQL، مفهوم **Instant Vector Selector** به کوئری‌هایی اشاره دارد که **سری‌های زمانی (time series) را در یک لحظه مشخص** انتخاب می‌کنند — معمول‌ترین و پرکاربردترین نوع کوئری در PromQL است.

---

## 🟢 تعریف ساده:

**Instant Vector Selector** یعنی:
🔹 «انتخاب یک یا چند سری زمانی در یک نقطه زمانی خاص (معمولاً "الان")»

---

## ✅ ویژگی‌ها:

| ویژگی                     | توضیح                                                                      |
| ------------------------- | -------------------------------------------------------------------------- |
| فقط یک مقدار در زمان مشخص | از هر سری زمانی، مقدار مربوط به "الان" (یا یک زمان خاص) را برمی‌گرداند     |
| بدون بازه زمانی `[ ]`     | اگر بازه مثل `[5m]` نداشته باشه، یعنی Instant Vector هست                   |
| پایه‌ی اکثر کوئری‌ها      | نمایش‌های گراف، alertها و بسیاری از محاسبات از Instant Vector شروع می‌شوند |

---

## 📘 مثال‌ها:

### 1. ساده‌ترین حالت:

```promql
up
```

این کوئری همه سری‌های زمانی با نام `up` را در همین لحظه (now) انتخاب می‌کند.

---

### 2. با Label Selector:

```promql
node_cpu_seconds_total{mode="idle", instance="server01:9100"}
```

این یک Instant Vector Selector است که فقط سری‌هایی با این برچسب‌ها را در لحظه جاری انتخاب می‌کند.

---

### 3. با Offset (جابه‌جایی زمانی):

```promql
http_requests_total offset 1h
```

مقدار متریک را در **دقیقاً 1 ساعت قبل** برمی‌گرداند — همچنان یک Instant Vector است، ولی نه برای "الان".

---

### 4. اشتباه رایج – Range Vector نیست:

```promql
rate(http_requests_total[5m])   ✅ این یک Range Vector است، نه Instant
```

---

## 🧠 یادآوری:

**اگر کوئری هیچ `[ ]` نداشته باشد و متریک را در لحظه‌ای مشخص (الان یا با offset) بخواند، آن را یک Instant Vector Selector می‌نامیم.**


‏`Instant vector selectors` امکان انتخاب مجموعه‌ای از `time series` و یک مقدار `sample` برای هر کدام در یک `timestamp` معین (نقطه زمانی) را فراهم می‌کنند. در ساده‌ترین شکل، فقط یک `metric name` مشخص می‌شود، که منجر به یک `instant vector` حاوی عناصری برای تمام `time series`هایی می‌شود که این `metric name` را دارند.

مقدار بازگشتی، مقدار جدیدترین `sample` در یا قبل از `timestamp` ارزیابی `query` (در مورد `instant query`) یا گام فعلی درون `query` (در مورد `range query`) خواهد بود. `@ modifier` اجازه می‌دهد `timestamp`ی که انتخاب نسبت به آن انجام می‌شود را بازنویسی کنید. `Time series`ها تنها در صورتی بازگردانده می‌شوند که جدیدترین `sample` آنها کمتر از دوره `lookback period` پیش باشد.

این مثال تمام `time series`هایی را که `metric name` `http_requests_total` دارند انتخاب می‌کند و جدیدترین `sample` را برای هر کدام بازمی‌گرداند:

```promql
http_requests_total
```

امکان فیلتر کردن بیشتر این `time series`ها با افزودن لیستی از `label matchers` جدا شده با کاما در `curly braces ({})` وجود دارد.

این مثال فقط آن `time series`هایی را با `metric name` `http_requests_total` انتخاب می‌کند که همچنین `label` `job` آنها روی `prometheus` و `label` `group` آنها روی `canary` تنظیم شده باشد:

```promql
http_requests_total{job="prometheus",group="canary"}
```

همچنین امکان تطبیق منفی یک `label value`، یا تطبیق `label values` با `regular expressions` وجود دارد. `label matching operators` زیر وجود دارند:

*   `=`: انتخاب `labels` که دقیقاً برابر با `string` ارائه شده هستند.
*   `!=`: انتخاب `labels` که برابر با `string` ارائه شده نیستند.
*   `=~`: انتخاب `labels` که با `regex` با `string` ارائه شده مطابقت دارند.
*   `!~`: انتخاب `labels` که با `regex` با `string` ارائه شده مطابقت ندارند.

تطابق‌های `Regex` کاملاً `anchored` هستند. تطابق `"env=~"foo` به عنوان `"$env=~"^foo` در نظر گرفته می‌شود.

به عنوان مثال، این مورد تمام `time series`های `http_requests_total` را برای محیط‌های `staging`، `testing` و `development` و `HTTP methods` غیر از `GET` انتخاب می‌کند.

```promql
http_requests_total{environment=~"staging|testing|development",method!="GET"}
```

‏`Label matchers` که با `label values` خالی مطابقت دارند، همچنین تمام `time series`هایی را انتخاب می‌کنند که `label` مشخص شده را اصلاً ندارند. امکان داشتن چندین `matcher` برای یک `label name` وجود دارد.

به عنوان مثال، با توجه به `dataset`:

```
http_requests_total
http_requests_total{replica="rep-a"}
http_requests_total{replica="rep-b"}
http_requests_total{environment="development"}
```

`query` `http_requests_total{environment=""}` 
مطابقت داده و بازمی‌گرداند:

```
http_requests_total
http_requests_total{replica="rep-a"}
http_requests_total{replica="rep-b"}
```

و مستثنی می‌کند:

```
http_requests_total{environment="development"}
```

می‌توان از چندین `matcher` برای یک `label name` استفاده کرد؛ همه آنها باید برقرار باشند تا نتیجه‌ای بازگردانده شود.

`query`:

```promql
http_requests_total{replica!="rep-a",replica=~"rep.*"}
```

سپس مطابقت می‌دهد:

```
http_requests_total{replica="rep-b"}
```

‏`Vector selectors` باید یا یک نام مشخص کنند یا حداقل یک `label matcher` که با `string` خالی مطابقت نداشته باشد. `expression` زیر غیرقانونی است:

```promql
{job=~".*"} # Bad!
```

در مقابل، این `expression`ها معتبر هستند زیرا هر دو دارای یک `selector` هستند که با `label values` خالی مطابقت ندارد.

```promql
{job=~".+"}              # Good!
{job=~".*",method="get"} # Good!
```

‏`Label matchers` را می‌توان با تطبیق با `label` داخلی `__name__` به `metric names` نیز اعمال کرد. به عنوان مثال، `expression` `http_requests_total` معادل `{__name__="http_requests_total"}` است. `Matchers` غیر از `=` (`!=`, `=~`, `!~`) نیز ممکن است استفاده شوند. `expression` زیر تمام `metrics` را که نامشان با `job:` شروع می‌شود انتخاب می‌کند:

```promql
{__name__=~"job:.*"}
```

‏`Metric name` نباید یکی از `keywords` `bool`, `on`, `ignoring`, `group_left` و `group_right` باشد. `expression` زیر غیرقانونی است:

```promql
on{} # Bad!
```

یک راه حل برای این محدودیت استفاده از `label` `__name__` است:

```promql
{__name__="on"} # Good!
```


در Prometheus و زبان **PromQL**، عبارت **Range Vector Selector** برای انتخاب **مجموعه‌ای از داده‌های سری زمانی در یک بازه زمانی مشخص** به کار می‌رود.

---

## 🟦 تعریف ساده:

**Range Vector Selector** یعنی:

> انتخاب سری زمانی به‌همراه تمام مقدارهای (نمونه‌های) آن در یک **بازه زمانی مشخص**، مثلاً در 5 دقیقه یا 1 ساعت گذشته.

---

## 📘 ساختار کلی:

```
<metric_name>[<duration>]
```

مثال:

```promql
http_requests_total[5m]
```

یعنی: تمام نمونه‌های (samples) سری زمانی `http_requests_total` را در **۵ دقیقه‌ی گذشته** انتخاب کن.

---

## ✅ ویژگی‌های Range Vector:

| ویژگی                                  | توضیح                                                 |
| -------------------------------------- | ----------------------------------------------------- |
| شامل چندین نقطه داده                   | برخلاف Instant Vector که فقط ۱ مقدار در زمان خاص دارد |
| بازه زمانی صریح دارد                   | مثلاً `[1m]`, `[5m]`, `[1h]`, `[30s]`                 |
| فقط با بعضی توابع کار می‌کند           | مثل `rate()`, `avg_over_time()`, `max_over_time()`    |
| پایه‌ی محاسبهٔ نرخ، میانگین، درصد و... | چون برای تحلیل روند داده‌ها استفاده می‌شود            |

---

## 🔍 مثال‌های رایج:

### 1. فقط انتخاب بازه:

```promql
node_cpu_seconds_total[2m]
```

تمام داده‌های این متریک در ۲ دقیقه گذشته — خروجی یک **Range Vector** است.

---

### 2. استفاده در تابع `rate()`:

```promql
rate(http_requests_total[1m])
```

میانگین نرخ رشد این متریک در ۱ دقیقه گذشته را حساب می‌کند.

---

### 3. استفاده در تابع `avg_over_time()`:

```promql
avg_over_time(node_load1[15m])
```

میانگین بار سیستم در ۱۵ دقیقه گذشته.

---

## ⚠️ تفاوت با Instant Vector:

| ویژگی  | Instant Vector       | Range Vector               |
| ------ | -------------------- | -------------------------- |
| زمان   | فقط یک لحظه          | یک بازه زمانی              |
| مقدار  | فقط آخرین نمونه      | چندین نمونه                |
| نماد   | بدون `[ ]`           | با `[duration]`            |
| کاربرد | مشاهده وضعیت لحظه‌ای | تحلیل آماری روی زمان گذشته |

---

## 🧠 نکته:

🔹 اگر در PromQL دیدی متریکی در کنار `[5m]` یا `[30s]` آمده، بدون هیچ تابعی، اون **Range Vector Selector** است.
🔹 اما برای اینکه معنی‌دار بشه (مثلاً rate محاسبه شه)، معمولاً باید با توابع استفاده بشه.

---




**`Range Vector Selectors`**
‏`Range vector literals` مانند `instant vector literals` کار می‌کنند، با این تفاوت که محدوده‌ای از `samples` را از لحظه فعلی به عقب انتخاب می‌کنند. از نظر نحوی، یک `float literal` در `square brackets ([])` در انتهای یک `vector selector` اضافه می‌شود تا مشخص کند مقادیر برای چند ثانیه قبل در زمان باید برای هر عنصر `range vector` نتیجه `fetch` شوند. معمولاً، `float literal` از سینتکس با یک یا چند `time units` استفاده می‌کند، مثلاً `[5m]`. محدوده یک `interval` باز از چپ و بسته از راست است، یعنی `samples` با `timestamps` منطبق بر مرز چپ محدوده از انتخاب مستثنی می‌شوند، در حالی که `samples` منطبق بر مرز راست محدوده در انتخاب گنجانده می‌شوند.

در این مثال، ما تمام مقادیر ثبت شده کمتر از 5 دقیقه پیش را برای تمام `time series`هایی که `metric name` `http_requests_total` و `label` `job` آنها روی `prometheus` تنظیم شده است، انتخاب می‌کنیم:

```promql
http_requests_total{job="prometheus"}[5m]
```

**`Offset modifier`**
‏`offset modifier` امکان تغییر `offset` زمانی برای `instant` و `range vectors` مجزا در یک `query` را فراهم می‌کند.

به عنوان مثال، `expression` زیر مقدار `http_requests_total` را 5 دقیقه قبل نسبت به زمان ارزیابی `query` فعلی بازمی‌گرداند:

```promql
http_requests_total offset 5m
```

توجه داشته باشید که `offset modifier` همیشه باید بلافاصله بعد از `selector` بیاید، یعنی عبارت زیر صحیح خواهد بود:

```promql
sum(http_requests_total{method="GET"} offset 5m) // GOOD.
```

در حالی که عبارت زیر نادرست خواهد بود:

```promql
sum(http_requests_total{method="GET"}) offset 5m // INVALID.
```

همین امر برای `range vectors` نیز صادق است. این `rate` 5 دقیقه‌ای `http_requests_total` را یک هفته پیش بازمی‌گرداند:

```promql
rate(http_requests_total[5m] offset 1w)
```

هنگام `query` برای `samples` در گذشته، یک `offset` منفی امکان مقایسه‌های زمانی به جلو در زمان را فراهم می‌کند:

```promql
rate(http_requests_total[5m] offset -1w)
```

توجه داشته باشید که این به `query` اجازه می‌دهد جلوتر از زمان ارزیابی خود را نگاه کند.

**`@ modifier`**
‏`@ modifier` امکان تغییر زمان ارزیابی برای `instant` و `range vectors` مجزا در یک `query` را فراهم می‌کند. زمان ارائه شده به `@ modifier` یک `unix timestamp` است و با یک `float literal` توصیف می‌شود.

به عنوان مثال، `expression` زیر مقدار `http_requests_total` را در `2021-01-04T07:40:00+00:00` بازمی‌گرداند:

```promql
http_requests_total @ 1609746000
```

توجه داشته باشید که `@ modifier` همیشه باید بلافاصله بعد از `selector` بیاید، یعنی عبارت زیر صحیح خواهد بود:

```promql
sum(http_requests_total{method="GET"} @ 1609746000) // GOOD.
```

در حالی که عبارت زیر نادرست خواهد بود:

```promql
sum(http_requests_total{method="GET"}) @ 1609746000 // INVALID.
```

همین امر برای `range vectors` نیز صادق است. این `rate` 5 دقیقه‌ای `http_requests_total` را در `2021-01-04T07:40:00+00:00` بازمی‌گرداند:

```promql
rate(http_requests_total[5m] @ 1609746000)
```

‏`@ modifier` از تمام نمایش‌های `numeric literals` توصیف شده در بالا پشتیبانی می‌کند. این با `offset modifier` کار می‌کند که در آن `offset` نسبت به زمان `@ modifier` اعمال می‌شود. نتایج بدون توجه به ترتیب `modifier`ها یکسان هستند.

به عنوان مثال، این دو `query` نتیجه یکسانی تولید خواهند کرد:

```promql
# offset after @
http_requests_total @ 1609746000 offset 5m
# offset before @
http_requests_total offset 5m @ 1609746000
```

علاوه بر این، `start()` و `end()` نیز می‌توانند به عنوان مقادیر ویژه برای `@ modifier` استفاده شوند.

برای یک `range query`، آنها به ترتیب به شروع و پایان `range query` تبدیل می‌شوند و برای همه گام‌ها یکسان باقی می‌مانند.

برای یک `instant query`، `start()` و `end()` هر دو به زمان ارزیابی تبدیل می‌شوند.

```promql
http_requests_total @ start()
rate(http_requests_total[5m] @ end())
```

توجه داشته باشید که `@ modifier` به `query` اجازه می‌دهد جلوتر از زمان ارزیابی خود را نگاه کند.

**`Subquery`**
‏`Subquery` به شما امکان می‌دهد یک `instant query` را برای یک `range` و `resolution` معین اجرا کنید. نتیجه یک `subquery` یک `range vector` است.

سینتکس: `<instant_query> '[' <range> ':' [<resolution>] ']' [ @ <float_literal> ] [ offset <float_literal> ]`

‏`<resolution>` اختیاری است. پیش‌فرض فاصله ارزیابی سراسری است.

**`Operators`**
‏`Prometheus` از بسیاری از `binary` و `aggregation operators` پشتیبانی می‌کند. اینها به تفصیل در صفحه `expression language operators` شرح داده شده‌اند.

**`Functions`**
‏`Prometheus` از چندین `function` برای عملیات روی داده‌ها پشتیبانی می‌کند. اینها به تفصیل در صفحه `expression language functions` شرح داده شده‌اند.

**`Comments`**
‏`PromQL` از `line comments` که با `#` شروع می‌شوند پشتیبانی می‌کند. مثال:

```promql
    # This is a comment
```

**`Regular expressions`**
تمام `regular expressions` در `Prometheus` از سینتکس `RE2` استفاده می‌کنند.

تطابق‌های `Regex` همیشه کاملاً `anchored` هستند.

**`Gotchas`**

**`Staleness`**
‏`timestamps`ی که در آنها داده‌ها در طول یک `query` نمونه‌برداری می‌شوند، مستقل از داده‌های `time series` واقعی موجود انتخاب می‌شوند. این عمدتاً برای پشتیبانی از مواردی مانند `aggregation` (`sum`, `avg` و غیره) است، جایی که چندین `time series` `aggregate` شده دقیقاً در زمان تراز نیستند. به دلیل استقلال آنها، `Prometheus` نیاز دارد در آن `timestamps` برای هر `time series` مربوطه مقداری را تخصیص دهد. این کار را با گرفتن جدیدترین `sample` که کمتر از دوره `lookback period` پیش است، انجام می‌دهد. `lookback period` به طور پیش‌فرض 5 دقیقه است، اما می‌توان آن را با پرچم `--query.lookback-delta` تنظیم کرد.

اگر یک `target scrape` یا `rule evaluation` دیگر `sample`ی را برای یک `time series` که قبلاً وجود داشت برنگرداند، این `time series` به عنوان `stale` علامت‌گذاری می‌شود. اگر یک `target` حذف شود، `time series`های قبلاً بازیابی شده به زودی پس از حذف به عنوان `stale` علامت‌گذاری می‌شوند.

اگر یک `query` در یک `sampling timestamp` پس از اینکه یک `time series` به عنوان `stale` علامت‌گذاری شد ارزیابی شود، هیچ مقداری برای آن `time series` بازگردانده نمی‌شود. اگر `sample`های جدید متعاقباً برای آن `time series` `ingest` شوند، طبق انتظار بازگردانده خواهند شد.

یک `time series` زمانی `stale` می‌شود که دیگر `export` نشود، یا `target` دیگر وجود نداشته باشد. چنین `time series`هایی در زمان آخرین `sample` جمع‌آوری شده خود از `graphs` ناپدید می‌شوند و پس از علامت‌گذاری به عنوان `stale` در `queries` بازگردانده نخواهند شد.

برخی `exporters` که `timestamp`های خود را روی `sample`ها قرار می‌دهند، رفتار متفاوتی دارند: `series`هایی که `export` شدن آنها متوقف می‌شود، آخرین مقدار را برای (به طور پیش‌فرض) 5 دقیقه قبل از ناپدید شدن حفظ می‌کنند. تنظیم `track_timestamps_staleness` می‌تواند این را تغییر دهد.

**اجتناب از `queries` کند و `overloads`**
اگر یک `query` نیاز به عملیات بر روی مقدار قابل توجهی از داده داشته باشد، ترسیم `graph` آن ممکن است با `timeout` مواجه شود یا `server` یا `browser` را `overload` کند. بنابراین، هنگام ساخت `queries` بر روی داده‌های ناشناخته، همیشه ساخت `query` را در نمای `tabular view` `expression browser` `Prometheus` شروع کنید تا زمانی که مجموعه نتایج معقول به نظر برسد (حداکثر صدها، نه هزاران، `time series`). فقط زمانی که داده‌های خود را به اندازه کافی فیلتر یا `aggregate` کردید، به حالت `graph mode` بروید. اگر `expression` هنوز برای ترسیم `graph` به صورت `ad-hoc` بیش از حد طول می‌کشد، آن را از طریق یک `recording rule` پیش‌ضبط کنید.

این به ویژه برای زبان `query` `Prometheus` مرتبط است، جایی که یک `metric name selector` خالی مانند `api_http_requests_total` می‌تواند به هزاران `time series` با `labels` مختلف گسترش یابد. همچنین، به خاطر داشته باشید که `expressions` که بر روی بسیاری از `time series`ها `aggregate` می‌شوند، حتی اگر خروجی فقط تعداد کمی `time series` باشد، روی `server` بار ایجاد می‌کنند. این شبیه به این است که چگونه جمع کردن تمام مقادیر یک ستون در یک `relational database` کند خواهد بود، حتی اگر مقدار خروجی فقط یک عدد باشد.





‏

### مبانی Aggregation

بیایید با چند کوئری `aggregation` ساده شروع کنیم. این کوئری‌ها احتمالاً بیشتر کاربردهای بالقوه شما از `PromQL` را پوشش می‌دهند. اگرچه `PromQL` تا حد امکان قدرتمند است، در بیشتر مواقع نیازهای شما نسبتاً ساده خواهند بود.

## Gauge
‏`Gauges` تصویری لحظه‌ای (`snapshot`) از وضعیت (`state`) هستند و معمولاً هنگام تجمیع (`aggregating`) آن‌ها، می‌خواهید `sum` (مجموع)، `average` (میانگین)، `minimum` (کمینه) یا `maximum` (بیشینه) آن‌ها را بگیرید.

‏`metric` با نام `node_filesystem_size_bytes` از `Node Exporter` خود را در نظر بگیرید که اندازه هر یک از فایل‌سیستم‌های `mount` شده شما را گزارش می‌دهد و دارای `labels` با نام‌های `device`، `fstype` و `mountpoint` است. شما می‌توانید اندازه کل فایل‌سیستم را روی هر ماشین با استفاده از دستور زیر محاسبه کنید:

`sum without(device, fstype, mountpoint)(node_filesystem_size_bytes)`

این مورد به درستی کار می‌کند زیرا `without` به `aggregator` تابع `sum` می‌گوید که همه چیز را با `labels` یکسان جمع کند و آن سه `label` مشخص شده را نادیده بگیرد. بنابراین اگر `time series` زیر را داشتید:
```
node_filesystem_free_bytes{device="/dev/sda1",fstype="vfat",
instance="localhost:9100",job="node",mountpoint="/boot/efi"} 70300672
node_filesystem_free_bytes{device="/dev/sda5",fstype="ext4",
instance="localhost:9100",job="node",mountpoint="/"} 30791843840
node_filesystem_free_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run"} 817094656
node_filesystem_free_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/lock"} 5238784
node_filesystem_free_bytes{device="tmpfs",fstype="tmpfs",
instance="localhost:9100",job="node",mountpoint="/run/user/1000"} 826912768
```

نتیجه به این صورت خواهد بود:

`{instance="localhost:9100",job="node"} 32511390720`

متوجه خواهید شد که `labels` با نام‌های `device`، `fstype` و `mountpoint` اکنون حذف شده‌اند. نام `metric` نیز دیگر وجود ندارد، زیرا این دیگر `node_filesystem_free_bytes` نیست چون عملیات ریاضی روی آن انجام شده است. از آنجایی که تنها یک `Node Exporter` توسط `Prometheus` در حال `scrape` شدن است، تنها یک نتیجه وجود دارد، اما اگر تعداد بیشتری را `scrape` می‌کردید، برای هر یک از `Node Exporter`ها یک نتیجه داشتید.

می‌توانید یک قدم فراتر رفته و `label` با نام `instance` را نیز حذف کنید:

`sum without(device, fstype, mountpoint, instance)(node_filesystem_size_bytes)`

این همانطور که انتظار می‌رود `label` با نام `instance` را حذف می‌کند، اما مقدار مانند عبارت قبلی باقی می‌ماند زیرا تنها یک `Node Exporter` برای تجمیع `metrics` وجود دارد:

`{job="node"} 32511390720`

می‌توانید از همین رویکرد با سایر `aggregations` استفاده کنید. `max` به شما اندازه بزرگترین فایل‌سیستم `mount` شده روی هر ماشین را می‌گوید:

`max without(device, fstype, mountpoint)(node_filesystem_size_bytes)`

‏`labels` خروجی دقیقاً مشابه زمانی است که با استفاده از `sum` تجمیع کردید:

`{instance="localhost:9100",job="node"} 30792601600`

این قابل پیش‌بینی بودن در مورد اینکه کدام `labels` بازگردانده می‌شوند، برای تطبیق `vector` (`vector matching`) با `operators` مهم است، همانطور که در گذشته بحث خواهد شد.

شما محدود به تجمیع `metrics` درباره یک نوع `job` نیستید. به عنوان مثال، برای یافتن میانگین تعداد `file descriptors` باز در تمام `job`های خود، می‌توانید از دستور زیر استفاده کنید:

`avg without(instance, job)(process_open_fds)`

## Counter
‏`Counters` تعداد یا اندازه رویدادها را ردیابی می‌کنند و مقداری که برنامه‌های شما در `/metrics` خود نمایش می‌دهند، مقدار کل از زمان شروع آن است. اما این مقدار کل به تنهایی کاربرد چندانی برای شما ندارد؛ آنچه واقعاً می‌خواهید بدانید این است که `counter` با چه سرعتی در طول زمان افزایش می‌یابد. این کار معمولاً با استفاده از تابع `rate` انجام می‌شود، اگرچه توابع `increase` و `irate` نیز روی مقادیر `counter` عمل می‌کنند.

به عنوان مثال، برای محاسبه مقدار ترافیک شبکه دریافت شده در هر ثانیه، می‌توانید از دستور زیر استفاده کنید:

`rate(node_network_receive_bytes_total[5m])`

‏`[5m]` به `rate` می‌گوید که داده‌های ۵ دقیقه را ارائه دهد، بنابراین مقدار بازگشتی میانگینی در طول ۵ دقیقه گذشته خواهد بود:

```
{device="lo",instance="localhost:9100",job="node"} 1859.389655172414
{device="wlan0",instance="localhost:9100",job="node"} 1314.5034482758622
```

مقادیر در اینجا اعداد صحیح نیستند، زیرا پنجره ۵ دقیقه‌ای که `rate` به آن نگاه می‌کند، کاملاً با `samples` (نمونه‌هایی) که `Prometheus` جمع‌آوری (`scrape`) کرده است، هماهنگ نیست. مقداری تخمین (`estimation`) برای پر کردن شکاف‌های بین نقاط داده‌ای که دارید و مرزهای بازه استفاده می‌شود.

خروجی `rate` یک `gauge` است، بنابراین همان `aggregations` که برای `gauges` اعمال می‌شوند، در اینجا نیز کاربرد دارند. `metric` با نام `node_network_receive_bytes_total` دارای `label` با نام `device` است، بنابراین اگر آن را تجمیع کنید (`aggregate it away`)، کل بایت‌های دریافت شده در هر ماشین در هر ثانیه را به دست خواهید آورد:

`sum without(device)(rate(node_network_receive_bytes_total[5m]))`

اجرای این کوئری نتیجه‌ای مانند این به شما می‌دهد:

`{instance="localhost:9100",job="node"} 3173.8931034482762`

می‌توانید `time series` مورد درخواست را فیلتر کنید، بنابراین می‌توانید فقط به `eth0` نگاه کنید و سپس آن را در تمام ماشین‌ها با حذف `label` با نام `instance` تجمیع کنید:

`sum without(instance)(rate(node_network_receive_bytes_total{device="eth0"}[5m]))`

هنگامی که این کوئری را اجرا می‌کنید، `label` با نام `instance` حذف می‌شود، اما `label` با نام `device` باقی می‌ماند زیرا شما نخواستید حذف شود:

`{device="eth0",job="node"} 3173.8931034482762`

هیچ ترتیب یا سلسله مراتبی (`hierarchy`) در `labels` وجود ندارد، که به شما امکان می‌دهد با هر تعداد `label` که دوست دارید تجمیع کنید.


در Prometheus، توابع `rate`، `irate` و `increase` همگی برای کار با **counter**ها (مقادیر تجمعی) به کار می‌روند، ولی تفاوت‌هایی در نحوه محاسبه و کاربرد آن‌ها وجود دارد. در ادامه تفاوت این توابع را به صورت خلاصه و دقیق بررسی می‌کنیم:

---

## 🔁 1. `rate()`

* **تعریف**: نرخ متوسط تغییر مقدار یک counter در یک بازه زمانی مشخص.
* **فرمولی**: ```$(آخر - اول) ÷ مدت زمان بازه$```
```
rate(http_server_requests_seconds_count{uri!~".*actuator.*"}[40s])  
..  
(valT4-valT1)/(t4-t1)  
..  
//Rate  
(28-12)/(40-10) = 0.533  
..  
// Number of request  
(28-12)/(40-10)/(1/40) = 21,32


```
* **ویژگی‌ها**:

  * برای داده‌های **مستمر و صاف** استفاده می‌شود.
  * میانگین نرخ تغییر در کل بازه زمانی را حساب می‌کند.
  * **پیشنهاد Prometheus** برای بیشتر موارد.

**مثال**:

```promql
rate(http_requests_total[5m])
```

یعنی: نرخ میانگین درخواست‌های HTTP در ۵ دقیقه گذشته.

![[rate.png]]



![[rate2.png]]


---

## ⚡ 2. `irate()`

* **تعریف**: نرخ لحظه‌ای (تقریباً آخرین نرخ تغییر) بر اساس دو آخرین نمونه در بازه زمانی.
* **ویژگی‌ها**:

  * **سریع‌ترین تغییرات** را نشان می‌دهد.
  * بسیار **حساس به نویز** است.
  * برای **نمودارهای زنده / real-time** مناسب است.
  * اگر داده‌ها sparse یا نامنظم باشند، ممکن است unreliable شود.

**مثال**:

```promql
irate(cpu_usage_seconds_total[1m])
```

![[irate.png]]


Applying the formula to all the serie this is the result:

![[irate2.png]]


Joining the points, I get a line with a considerable variation between the points.

![[irate3.png]]


یعنی: نرخ لحظه‌ای استفاده از CPU بر اساس آخرین دو نمونه در ۱ دقیقه اخیر.

مقایسه irate , rate

![[rateVsIrate.png]]

---

## 📈 3. `increase()`

* **تعریف**: مجموع افزایش مقدار یک counter در طول بازه زمانی.
* **فرمولی**:` فقط $آخر - اول$ (نه تقسیم بر زمان).`
* **ویژگی‌ها**:

  * برای محاسبه **تعداد اتفاقات** در یک بازه زمانی استفاده می‌شود.
  * مناسب برای محاسبه تعداد کل درخواست‌ها، خطاها، bytes و غیره.

**مثال**:

```promql
increase(http_requests_total[1h])
```

یعنی: تعداد کل درخواست‌های HTTP که در یک ساعت اخیر ثبت شده‌اند.

---

## 📊 جدول مقایسه

| ویژگی                  | `rate()`       | `irate()`              | `increase()`   |
| ---------------------- | -------------- | ---------------------- | -------------- |
| نوع خروجی              | نرخ (در ثانیه) | نرخ لحظه‌ای (در ثانیه) | مجموع افزایش   |
| حساسیت به نویز         | کم             | زیاد                   | کم             |
| مناسب برای نمودار زنده | نه خیلی        | بله                    | نه             |
| مناسب برای شمارش تجمعی | نه             | نه                     | بله            |
| استفاده رایج           | بله            | بله (برای لحظه‌ای)     | بله (برای جمع) |

---

### ✅ چه زمانی از کدام استفاده کنیم؟

* 🔸 `rate`: برای مشاهده نرخ تغییر یکنواخت در طول زمان (اغلب بهترین انتخاب پیش‌فرض).
* 🔸 `irate`: برای مشاهده نوسانات سریع و نرخ لحظه‌ای (نمودار زنده).
* 🔸 `increase`: برای محاسبه تعداد یا مقدار تجمعی رخدادها در یک بازه زمانی.



## Summary

یک `metric` از نوع `summary` معمولاً شامل `_sum` و `_count` و گاهی اوقات یک `time series` بدون پسوند با یک `label` به نام `quantile` است. `_sum` و `_count` هر دو `counter` هستند.
‏`Prometheus` شما یک `summary` به نام `http_response_size_bytes` برای مقدار داده‌ای که برخی از `HTTP APIs` آن بازمی‌گردانند، ارائه می‌دهد. `http_response_size_bytes_count` تعداد درخواست‌ها را ردیابی می‌کند و از آنجایی که یک `counter` است، باید قبل از تجمیع `label` با نام `handler` آن، از `rate` استفاده کنید:

`sum without(handler)(rate(http_response_size_bytes_count[5m]))`

این به شما نرخ کل درخواست HTTP در ثانیه را می‌دهد و از آنجایی که `Node Exporter` نیز این `metric` را برمی‌گرداند، هر دو `job` را در نتیجه خواهید دید:

```
{instance="localhost:9090",job="prometheus"} 0.26868836781609196
{instance="localhost:9100",job="node"} 0.1
```

به طور مشابه، `http_response_size_bytes_sum` یک `counter` با تعداد بایت‌هایی است که هر `handler` بازگردانده است، بنابراین همان الگو اعمال می‌شود:

`sum without(handler)(rate(http_response_size_bytes_sum[5m]))`

این نتایجی با همان `labels` کوئری قبلی برمی‌گرداند، اما مقادیر بزرگتر هستند زیرا پاسخ‌ها معمولاً بایت‌های زیادی را برمی‌گردانند:

```
{instance="localhost:9090",job="prometheus"} 796.0015958275862
{instance="localhost:9100",job="node"} 1581.6103448275862
```

قدرت یک `summary` در این است که به شما امکان می‌دهد اندازه متوسط یک رویداد را محاسبه کنید، در این مورد، میانگین مقدار بایت‌هایی که در هر پاسخ بازگردانده می‌شود. اگر سه پاسخ با اندازه‌های ۱، ۴ و ۷ داشتید، میانگین آن‌ها مجموع تقسیم بر تعدادشان خواهد بود، یعنی ۱۲ تقسیم بر ۳. همین امر در مورد `summary` نیز صدق می‌کند. شما `_sum` را بر `_count` (پس از اعمال `rate`) تقسیم می‌کنید تا میانگین را در یک دوره زمانی به دست آورید:

```
sum without(handler)(rate(http_response_size_bytes_sum[5m]))
/
sum without(handler)(rate(http_response_size_bytes_count[5m]))
```

‏`operator` تقسیم (`division operator`) `time series` با `labels` یکسان را `matches` (تطبیق می‌دهد) و تقسیم می‌کند، و همان دو `time series` را به شما می‌دهد اما با مقدار میانگین اندازه پاسخ در ۵ دقیقه گذشته:

```
{instance="localhost:9090",job="prometheus"} 2962.54580091246150133317
{instance="localhost:9100",job="node"} 15816.10344827586200000000
```

هنگام محاسبه میانگین، مهم است که ابتدا `sum` و `count` را تجمیع کنید و تنها به عنوان آخرین مرحله، تقسیم را انجام دهید. در غیر این صورت، ممکن است به میانگین‌گیری از میانگین‌ها ختم شوید که از نظر آماری معتبر نیست.

به عنوان مثال، اگر می‌خواستید میانگین اندازه پاسخ را در تمام `instances` یک `job` به دست آورید، می‌توانید از دستور زیر استفاده کنید:

```
sum without(instance)(
  sum without(handler)(rate(http_response_size_bytes_sum[5m]))
)
/
sum without(instance)(
  sum without(handler)(rate(http_response_size_bytes_count[5m]))
)
```

با این حال، انجام این کار نادرست خواهد بود:

```
avg without(instance)(
  sum without(handler)(rate(http_response_size_bytes_sum[5m]))
  /
  sum without(handler)(rate(http_response_size_bytes_count[5m]))
)
```

میانگین گرفتن از یک میانگین نادرست است و هم تقسیم و هم `avg` در حال محاسبه میانگین خواهند بود.

از دیدگاه آماری، امکان تجمیع `quantiles` یک `summary` (یعنی `time series` با `label` به نام `quantile`) وجود ندارد.

### Histogram

‏`Metrics` از نوع `Histogram` به شما امکان می‌دهند توزیع (`distribution`) اندازه رویدادها را ردیابی کنید و به شما امکان محاسبه `quantiles` از آن‌ها را می‌دهند. به عنوان مثال، می‌توانید از `histograms` برای محاسبه `quantile` ۰.۹ (که به عنوان `percentile` نودم نیز شناخته می‌شود) `latency` استفاده کنید. Prometheus 2.37.1 یک `metric` از نوع `histogram` به نام `prometheus_tsdb_compaction_duration_seconds` را ارائه می‌دهد که مدت زمان `compaction` (فشرده‌سازی) پایگاه داده `time series` را بر حسب ثانیه ردیابی می‌کند. این `metric` از نوع `histogram` دارای `time series` با پسوند `_bucket` به نام `prometheus_tsdb_compaction_duration_seconds_bucket` است. هر `bucket` دارای یک `label` به نام `le` است که یک `counter` از تعداد رویدادهایی است که اندازه‌ای کمتر یا مساوی مرز `bucket` (`bucket boundary`) دارند. این یک جزئیات پیاده‌سازی است که عمدتاً نیازی به نگرانی در مورد آن ندارید زیرا تابع `histogram_quantile` هنگام محاسبه `quantiles` این موضوع را مدیریت می‌کند. به عنوان مثال، `quantile` ۰.۹۰ به صورت زیر خواهد بود:

```
histogram_quantile(
  0.90,
  rate(prometheus_tsdb_compaction_duration_seconds_bucket[1d])
)
```

از آنجایی که `prometheus_tsdb_compaction_duration_seconds_bucket` یک `counter` است، ابتدا باید از `rate` استفاده کنید. `Compaction` معمولاً هر دو ساعت یکبار اتفاق می‌افتد، بنابراین از یک بازه زمانی یک روزه (`[1d]`) در اینجا استفاده می‌شود و نتیجه‌ای مانند این را در مرورگر `expression` (`expression browser`) خواهید دید:

`{instance="localhost:9090",job="prometheus"} 7.720000000000001`

این نشان می‌دهد که `percentile` نودم `latency` مربوط به `compaction` حدود ۷.۷۲ ثانیه است. از آنجایی که معمولاً فقط ۱۲ `compaction` در یک روز وجود خواهد داشت، `percentile` نودم می‌گوید که ۱۰٪ از `compaction`ها بیشتر از این طول می‌کشند، یعنی یک یا دو `compaction`. این نکته‌ای است که هنگام استفاده از `quantiles` باید به آن توجه داشت. به عنوان مثال، اگر می‌خواهید `quantile` ۰.۹۹۹ را محاسبه کنید، باید چندین هزار نقطه داده برای کار داشته باشید تا پاسخ نسبتاً دقیقی تولید شود. اگر کمتر از آن دارید، تک داده‌های پرت (`outliers`) می‌توانند به شدت بر نتیجه تأثیر بگذارند و باید به استفاده از `quantiles` پایین‌تر فکر کنید تا از بیان اظهاراتی در مورد سیستم خود که داده‌های کافی برای پشتیبانی از آن‌ها ندارید، اجتناب کنید.
معمولاً با `histograms` از `rate` ۵ یا ۱۰ دقیقه‌ای استفاده می‌کنید. تمام `time series` مربوط به `bucket`ها همراه با هر `label` و یک بازه طولانی در `rate`، می‌تواند منجر به تعداد زیادی `sample` شود که نیاز به پردازش دارند. مراقب عبارات (`expressions`) `PromQL` باشید که از بازه‌هایی به طول ساعت‌ها یا روزها استفاده می‌کنند، زیرا محاسبه آن‌ها می‌تواند نسبتاً پرهزینه باشد.

مشابه زمانی که میانگین می‌گیرید، استفاده از `histogram_quantile` باید آخرین مرحله در یک عبارت کوئری باشد. `Quantiles` از دیدگاه آماری قابل تجمیع نیستند یا نمی‌توان روی آن‌ها عملیات حسابی انجام داد. بر این اساس، هنگامی که می‌خواهید `histogram` یک تجمیع را بگیرید، ابتدا با `sum` تجمیع کنید و سپس از `histogram_quantile` استفاده کنید:

```
histogram_quantile(
  0.90,
  sum without(instance)(rate(prometheus_tsdb_compaction_duration_bucket[1d]))
)
```

این `quantile` ۰.۹ مدت زمان `compaction` را در تمام سرورهای `Prometheus` شما محاسبه می‌کند و نتیجه‌ای بدون `label` با نام `instance` تولید می‌کند:

`{job="prometheus"} 7.720000000000001`

‏`Metrics` از نوع `Histogram` همچنین شامل `metrics` با نام `_sum` و `_count` هستند که دقیقاً مانند `metric` از نوع `summary` کار می‌کنند. می‌توانید از این‌ها برای محاسبه میانگین اندازه رویدادها، مانند میانگین مدت زمان `compaction` استفاده کنید:

```
sum without(instance)(rate(prometheus_tsdb_compaction_duration_sum[1d]))
/
sum without(instance)(rate(prometheus_tsdb_compaction_duration_count[1d]))
```

این نتیجه‌ای مانند این تولید می‌کند:

`{job="prometheus"} 3.1766430400714287`

### Selectors

کار با تمام `time series`های مختلف با مقادیر `label` متفاوت برای یک `metric` می‌تواند کمی طاقت‌فرسا و بالقوه گیج‌کننده باشد اگر `metric` از انواع مختلف سرورها  می‌آید. معمولاً می‌خواهید `time series`هایی را که روی آن‌ها کار می‌کنید محدود کنید. تقریباً همیشه می‌خواهید بر اساس `label` با نام `job` محدود کنید و بسته به کاری که انجام می‌دهید، ممکن است بخواهید فقط به یک `instance` یا یک `handler` نگاه کنید، به عنوان مثال.

این محدود کردن بر اساس `labels` با استفاده از `selectors` (انتخابگرها) انجام می‌شود. شما در هر مثالی تاکنون `selectors` را دیده‌اید و اکنون می‌خواهیم آن‌ها را به تفصیل برای شما توضیح دهیم. به عنوان مثال:

`process_resident_memory_bytes{job="node"}`

یک `selector` است که تمام `time series` با نام `process_resident_memory_bytes` و `label` با نام `job` و مقدار `node` را برمی‌گرداند. این `selector` خاص به درستی `instant vector selector` نامیده می‌شود، زیرا مقادیر `time series` داده شده را در یک لحظه معین برمی‌گرداند. `Vector` در اینجا اساساً به معنای یک لیست یک بعدی است، زیرا یک `selector` می‌تواند صفر یا چند `time series` را برگرداند و هر `time series` یک `sample` خواهد داشت.

‏`job="node"` یک `matcher` (تطبیق‌دهنده) نامیده می‌شود و شما می‌توانید چندین `matcher` در یک `selector` داشته باشید که با هم `AND` می‌شوند.

## Matchers

چهار `matcher` وجود دارد (شما قبلاً `equality matcher` را دیده‌اید که رایج‌ترین نیز هست):

=
این `equality matcher` است؛ به عنوان مثال، `job="node"`. با این کار می‌توانید مشخص کنید که `time series` بازگشتی دارای نام `label` با مقدار دقیق `label` داده شده باشد. از آنجایی که یک مقدار `label` خالی همانند نداشتن آن `label` است، می‌توانید از `foo=""` برای مشخص کردن عدم وجود `label` با نام `foo` استفاده کنید.

!=
این `negative equality matcher` است؛ به عنوان مثال، `job!="node"`. با این کار می‌توانید مشخص کنید که `time series` بازگشتی نام `label` با مقدار دقیق `label` داده شده را نداشته باشند.

=~
این `regular expression matcher` است؛ به عنوان مثال، `job=~"n.*"`. با این کار مشخص می‌کنید که برای `time series` بازگشتی، مقدار `label` داده شده با `regular expression` (عبارت منظم) مطابقت داشته باشد. `regular expression` کاملاً `anchored` (مهار شده) است، به این معنی که `regular expression` `a` فقط با رشته `a` مطابقت دارد و نه `xa` یا `ax`. اگر این رفتار را نمی‌خواهید، می‌توانید پیشوند یا پسوند `.*` را به `regular expression` خود اضافه کنید.⁶ همانند `relabeling`، از موتور `regular expression` `RE2` استفاده می‌شود، همانطور که در "Regular Expressions" قبلا پوشش داده شده است.

!~
این `negative regular expression matcher` است. `RE2` از عبارات `negative lookahead` پشتیبانی نمی‌کند، بنابراین این روش جایگزینی برای حذف مقادیر `label` بر اساس یک `regular expression` به شما ارائه می‌دهد. می‌توانید چندین `matcher` با نام `label` یکسان در یک `selector` داشته باشید که می‌تواند جایگزینی برای عبارات `negative lookahead` باشد. به عنوان مثال، برای یافتن اندازه تمام فایل‌سیستم‌های `mount` شده زیر `/run` اما نه `/run/user`، می‌توانید از دستور زیر استفاده کنید:
`node_filesystem_size_bytes{job="node",mountpoint=~"/run/.*",mountpoint!~"/run/user/.*"}`

در داخل، نام `metric` در یک `label` به نام `__name__` ذخیره می‌شود (همانطور که در "Reserved Labels and __name__" در گذشته بحث شد)، بنابراین `process_resident_memory_bytes{job="node"}` یک `syntactic sugar` برای `{name="process_resident_memory_bytes",job="node"}` است. حتی می‌توانید روی نام `metric` از `regular expressions` استفاده کنید، اما این کار خارج از زمانی که در حال `debugging` (اشکال‌زدایی) `performance` (عملکرد) سرور `Prometheus` هستید، هوشمندانه نیست.

نیاز به استفاده از `regular expression matchers` کمی نشانه خوبی نیست. اگر متوجه شدید که زیاد از آن‌ها روی یک `label` خاص استفاده می‌کنید، در نظر بگیرید که آیا باید مقادیر `label` تطبیق داده شده را در یک مقدار ترکیب کنید. به عنوان مثال، برای کدهای وضعیت `HTTP` به جای انجام `code~="4.."` برای گرفتن 401ها، 404ها، 405ها و غیره، ممکن است آن‌ها را در یک مقدار `label` `4xx` ترکیب کنید و از `equality matcher` `code="4xx"` استفاده کنید.

‏`selector` `{}` یک `error` برمی‌گرداند که یک اقدام ایمنی (`safety measure`) برای جلوگیری از بازگرداندن تصادفی تمام `time series`های داخل سرور `Prometheus` است زیرا این کار می‌تواند پرهزینه باشد. به طور دقیق‌تر، حداقل یکی از `matchers` در یک `selector` نباید با رشته خالی مطابقت داشته باشد. بنابراین `{foo=""}` و `{foo=~".*"}` یک `error` برمی‌گردانند، در حالی که `{foo="",bar="x"}`, `{foo!=""}` یا `{foo=~".+"}` مجاز هستند.



## Durations
‏`Durations` (مدت زمان‌ها) در `Prometheus` همانطور که در `PromQL` و فایل پیکربندی استفاده می‌شوند، از چندین واحد پشتیبانی می‌کنند. شما قبلاً `m` برای دقیقه را دیده‌اید.

![[pic14.png]]

### Subqueries

در حالی که `range vectors` روی `time series` عمل می‌کنند، نمی‌توان آن‌ها را در ترکیب با توابع استفاده کرد. اگر می‌خواهید `max_over_time` را با `rate` ترکیب کنید، می‌توانید یا از `recording rules` استفاده کنید که نتیجه تابع `rate` را ضبط کرده و آن را به تابع `vector` منتقل می‌کند، یا می‌توانید از `subquery` استفاده کنید. `subquery` بخشی از یک کوئری است که به شما امکان می‌دهد یک `range query` را در داخل یک کوئری انجام دهید. سینتکس (`syntax`) برای `subquery` از براکت‌های مربع، مانند `range selectors` استفاده می‌کند. اما دو `duration` مختلف می‌گیرد: `range` و `resolution`. `range` بازه‌ای است که توسط `subquery` بازگردانده می‌شود و `resolution` به عنوان `step` (گام) عمل می‌کند:

`max_over_time( rate(http_requests_total[5m])[30m:1m] )`

کوئری فوق `rate(http_requests_total[5m])` را هر دقیقه (`1m`) برای ۳۰ دقیقه گذشته (`30m`) اجرا می‌کند، سپس نتیجه را به تابع `max_over_time()` می‌دهد. `resolution` می‌تواند حذف شود، مانند `[30m:]`. در این حالت، فاصله ارزیابی سراسری (`global evaluation interval`) به عنوان `resolution` استفاده می‌شود.

### Offset

یک اصلاح‌کننده (`modifier`) وجود دارد که می‌توانید با هر دو نوع `vector selector` به نام `offset` استفاده کنید. `offset` به شما امکان می‌دهد زمان ارزیابی یک کوئری را بگیرید و بر اساس هر `selector`، آن را در زمان به عقب ببرید. به عنوان مثال:

`process_resident_memory_bytes{job="node"} offset 1h`
استفاده از حافظه را یک ساعت قبل از زمان ارزیابی کوئری دریافت می‌کند. `offset` در کوئری‌های ساده مانند این زیاد استفاده نمی‌شود، زیرا تغییر زمان ارزیابی برای کل کوئری آسان‌تر خواهد بود. جایی که این می‌تواند مفید باشد زمانی است که فقط می‌خواهید یک `selector` را در یک عبارت کوئری تنظیم کنید. به عنوان مثال:

`process_resident_memory_bytes{job="node"}`
`-`
`process_resident_memory_bytes{job="node"} offset 1h`
تغییر در استفاده از حافظه در `Node Exporter` طی ساعت گذشته¹³ را نشان می‌دهد. همین رویکرد با `range vectors` نیز کار می‌کند:
`rate(process_cpu_seconds_total{job="node"}[5m])`
`-`
`rate(process_cpu_seconds_total{job="node"}[5m] offset 1h)`
‏`offset` به شما امکان می‌دهد به گذشته دورتر نگاه کنید، اما همچنین به آینده، با استفاده از `offset` منفی. این می‌تواند هنگام انجام `prediction` یا زمانی که `sample` `metrics` با واقعیت هم‌تراز نیست، استفاده شود:
`rate(process_cpu_seconds_total{job="node"}[5m]) offset -1h`
`-`
`rate(process_cpu_seconds_total{job="node"}[5m])`
توجه داشته باشید که این کوئری احتمالاً برای ساعت گذشته چیزی برنمی‌گرداند. `Grafana` ویژگی‌ای دارد که یک `panel` را در زمان به بازه زمانی متفاوتی نسبت به بقیه `dashboard` که بخشی از آن است، منتقل کند. در Grafana 5.0.0 می‌توانید این را در تب `Time range` ویرایشگر `panel` (`panel editor`) پیدا کنید.

### @ Modifier

مشابه `offset modifier`، `PromQL` از یک `@ modifier` پشتیبانی می‌کند که به شما امکان می‌دهد ارزیابی `vector selectors`، `range selectors` و `subqueries` را به یک زمان ارزیابی مجدد ثابت تغییر دهید. `@ modifier` می‌تواند با یک `Unix timestamp` استفاده شود. کوئری `http_requests_total @ 1667491200` مقدار `http_requests_total` را در 2022-11-03T16:00:00+00:00 برمی‌گرداند. کوئری `rate(http_requests_total[5m] @1667491200)` نرخ ۵ دقیقه‌ای `http_requests_total` را در همان زمان برمی‌گرداند. علاوه بر این، `start()` و `end()` می‌توانند به عنوان مقادیر برای `@ modifier` استفاده شوند. برای یک `range query`، آنها به ترتیب با شروع و پایان `range query` تطبیق می‌یابند. برای یک کوئری لحظه‌ای (`instant query`)، هر دو به زمان ارزیابی تطبیق می‌یابند. در عمل، می‌توان از `@ modifier` برای رسم نمودار تکامل `http_request_total` که دارای `rate` بالایی در انتهای بازه ارزیابی است، استفاده کرد:

```
rate(http_requests_total[1m])
and
topk(5, rate(http_requests_total[1h] @ end()))
```

ء`topk(5, rate(http_requests_total[1h] @ end()))` به عنوان یک تابع رتبه‌بندی (`ranking function`) عمل می‌کند و فقط مقادیر بالاتر را در انتهای بازه ارزیابی فیلتر می‌کند.

### HTTP API

‏`Prometheus` تعدادی `HTTP API` ارائه می‌دهد. آنهایی که بیشتر با آنها تعامل خواهید داشت `query` و `query_range` هستند که به شما امکان دسترسی به `PromQL` را می‌دهند و می‌توانند توسط ابزارهای `dashboarding` یا اسکریپت‌های گزارش‌گیری سفارشی استفاده شوند. تمام `endpoints` مورد علاقه زیر `/api/v1/` قرار دارند و فراتر از اجرای `PromQL`، می‌توانید `metadata` (فراداده) `time series` را جستجو کرده و اقدامات مدیریتی (`administrative actions`) مانند گرفتن `snapshots` و حذف (`deleting`) `time series` را انجام دهید. این `API`های دیگر عمدتاً برای ابزارهای `dashboarding` مانند `Grafana` که می‌توانند از `metadata` برای بهبود `UI` خود استفاده کنند و برای کسانی که `Prometheus` را مدیریت می‌کنند، مورد علاقه هستند، اما به اجرای `PromQL` مربوط نمی‌شوند.

### query

‏`endpoint` با نام `query`، یا به طور رسمی `/api/v1/query`، یک عبارت `PromQL` را در یک زمان معین اجرا می‌کند و نتیجه را برمی‌گرداند. به عنوان مثال، `http://localhost:9090/api/v1/query?query=process_resident_memory_bytes` نتایجی مانند این را برمی‌گرداند:

```json
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "__name__": "process_resident_memory_bytes",
          "instance": "localhost:9090",
          "job": "prometheus"
        },
        "value": [1517929228.782, "91656192"]
      },
      {
        "metric": {
          "__name__": "process_resident_memory_bytes",
          "instance": "localhost:9100",
          "job": "node"
        },
        "value": [1517929228.782, "15507456"]
      }
    ]
  }
}
```

‏`status` برابر `success` است، به این معنی که کوئری کار کرده است. اگر ناموفق بود، `status` برابر `error` بود و فیلد `error` جزئیات بیشتری را ارائه می‌داد. این نتیجه خاص یک `instant vector` است که می‌توانید از `"resultType": "vector"` متوجه شوید. برای هر یک از `samples` در نتیجه، `labels` در `map`ی به نام `metric` قرار دارند و مقدار `sample` در `list`ی به نام `value` است. عدد اول در `list` `value`، `timestamp` مربوط به `sample` بر حسب ثانیه است و دومی مقدار واقعی `sample` است. مقدار داخل یک `string` است، زیرا `JSON` نمی‌تواند مقادیر غیرواقعی مانند `NaN` و `+Inf` را نشان دهد. زمان تمام `samples` زمان ارزیابی کوئری خواهد بود، حتی اگر عبارت فقط از یک `instant vector selector` تشکیل شده باشد. در اینجا زمان ارزیابی کوئری به طور پیش‌فرض زمان فعلی است، اما می‌توانید با پارامتر `URL` به نام `time`، زمانی را مشخص کنید که می‌تواند `Unix time`، بر حسب ثانیه، یا زمان `RFC 3339` باشد. به عنوان مثال، `http://localhost:9090/api/v1/query?query=process_resident_memory_bytes&time=1514764800` کوئری را در نیمه شب اول ژانویه ۲۰۱۸ ارزیابی می‌کند.
شما همچنین می‌توانید از `range vectors` با `endpoint` `query` استفاده کنید. به عنوان مثال، 
`http://localhost:9090/api/v1/query=prometheus_tsdb_head_samples_appended_total[1m]` نتایجی مانند این را برمی‌گرداند:

```json
{
  "status": "success",
  "data": {
    "resultType": "matrix",
    "result": [
      {
        "metric": {
          "__name__": "process_resident_memory_bytes",
          "instance": "localhost:9090",
          "job": "prometheus"
        },
        "values": [
          [1518008453.662, "87318528"],
          [1518008463.662, "87318528"],
          [1518008473.662, "87318528"]
        ]
      },
      {
        "metric": {
          "__name__": "process_resident_memory_bytes",
          "instance": "localhost:9100",
          "job": "node"
        },
        "values": [
          [1518008444.819, "17043456"],
          [1518008454.819, "17043456"],
          [1518008464.819, "17043456"]
        ]
      }
    ]
  }
}
```

این با نتیجه `instant vector` قبلی متفاوت است، زیرا `resultType` اکنون `matrix` است و هر `time series` چندین مقدار دارد. هنگامی که با `range vector` استفاده می‌شود، `endpoint` `query` `samples` خام را برمی‌گرداند، اما مراقب باشید که داده‌های زیادی را به یکباره درخواست نکنید زیرا ممکن است حافظه یک طرف یا طرف دیگر تمام شود.
یک نوع نتیجه دیگر به نام `scalar` وجود دارد. `Scalars` `labels` ندارند، آنها فقط اعداد هستند.¹⁷ `http://localhost:9090/api/v1/query?query=42` نتیجه زیر را تولید می‌کند:

```json
{
  "status": "success",
  "data": {
    "resultType": "scalar",
    "result": [1518008879.023, "42"]
  }
}
```

### query_range

‏`endpoint` با نام `query range` در `/api/v1/query_range` `endpoint` اصلی `HTTP` `Prometheus` است که از آن استفاده خواهید کرد، زیرا این `endpoint` برای `graphing` (رسم نمودار) است. در پشت صحنه، `query_range` یک `syntactic sugar` (به علاوه برخی بهینه‌سازی‌های عملکردی) برای چندین فراخوانی به `endpoint` `query` است. علاوه بر پارامتر `URL` `query`، شما به `query_range` یک `start time` (زمان شروع)، یک `end time` (زمان پایان) و یک `step` (گام) ارائه می‌دهید. کوئری ابتدا در `start time` اجرا می‌شود. سپس `step` ثانیه پس از `start time` اجرا می‌شود. سپس دو برابر `step` ثانیه پس از `start time` اجرا می‌شود و به همین ترتیب، تا زمانی که زمان ارزیابی کوئری از `end time` تجاوز کند، متوقف می‌شود. تمام نتایج `instant vector`   از اجراهای مختلف در یک `range vector` ترکیب شده و بازگردانده می‌شوند.
به عنوان مثال، اگر می‌خواستید تعداد `samples`ی که `Prometheus` در ۱۵ دقیقه اول سال ۲۰۱۸ دریافت کرده است را کوئری کنید، می‌توانید دستور زیر را اجرا کنید: `http://localhost:9090/api/v1/query_range?query=rate(prometheus_tsdb_head_samples_appended_total[5m])&start=1514764800&end=1514765700&step=60`، که نتیجه‌ای مانند این تولید می‌کند:

```json
{
  "status": "success",
  "data": {
    "resultType": "matrix",
    "result": [
      {
        "metric": {
          "instance": "localhost:9090",
          "job": "prometheus"
        },
        "values": [
          [1514764800, "85.07241379310345"],
          [1514764860, "102.6793103448276"],
          [1514764920, "120.30344827586208"],
          [1514764980, "137.93103448275863"],
          [1514765040, "146.7586206896552"],
          [1514765100, "146.7793103448276"],
          [1514765160, "146.8"],
          [1514765220, "146.8"],
          [1514765280, "146.8"],
          [1514765340, "146.8"],
          [1514765400, "146.8"],
          [1514765460, "146.8"],
          [1514765520, "146.8"],
          [1514765580, "146.8"],
          [1514765640, "146.8"],
          [1514765700, "146.8"]
        ]
      }
    ]
  }
}
```

چند جنبه از این وجود دارد که باید به آنها توجه کنید. اول اینکه `timestamps` `sample`ها با `start time` و `step` هم‌تراز هستند، زیرا هر نتیجه از یک ارزیابی کوئری لحظه‌ای متفاوت می‌آید و نتایج کوئری لحظه‌ای همیشه از زمان ارزیابی خود به عنوان `timestamp` نتایج استفاده می‌کنند. دوم اینکه آخرین `sample` در اینجا در `end time` است، به این معنی که بازه `inclusive` (شامل) است و آخرین نقطه `end time` خواهد بود اگر اتفاقاً با `step` هم‌تراز شود. سوم اینکه ما یک بازه ۵ دقیقه‌ای برای تابع `rate` انتخاب کردیم که بزرگتر از `step` است. از آنجایی که `query_range` ارزیابی‌های کوئری لحظه‌ای تکراری را انجام می‌دهد، هیچ وضعیتی بین ارزیابی‌ها منتقل نمی‌شود. اگر بازه کوچکتر از `step` بود، ما از روی داده‌ها می‌پریدیم. به عنوان مثال، یک بازه ۱ دقیقه‌ای با یک`step`  از نوع ۵ دقیقه‌ای، ۸۰٪ از `samples` را نادیده می‌گرفت. برای جلوگیری از این، باید از بازه‌هایی استفاده کنید که حداقل یک یا دو فاصله `scrape` (`scrape intervals`) بزرگتر از `step`ی باشند که استفاده می‌کنید.

هنگام استفاده از `range vectors` با `query_range`، معمولاً باید از بازه‌ای استفاده کنید که طولانی‌تر از `step` شما باشد تا داده‌ای را از دست ندهید.

چهارم اینکه برخی از `samples` به طور خاص گرد نیستند، و اینکه هر عددی اصلاً گرد باشد به این دلیل است که این یک تنظیم ساده از مقادیر `sample` است. هنگام کار با `metrics`، داده‌های شما به ندرت کاملاً تمیز هستند؛ `targets` مختلف در زمان‌های مختلف `scrape` می‌شوند و `scrapes` می‌توانند به تأخیر بیفتند. هنگام انجام کوئری‌هایی که کاملاً با داده‌های زیربنایی هم‌تراز نیستند یا هنگام تجمیع در چندین میزبان، به ندرت نتایج گردی دریافت خواهید کرد. علاوه بر این، ماهیت محاسبات `floating-point` می‌تواند منجر به اعدادی شود که تقریباً گرد هستند.

در اینجا، برای هر `step` یک `sample` وجود دارد. اگر اتفاق می‌افتاد که برای یک `time series` معین برای یک `step` نتیجه‌ای وجود نداشت، آن `sample` به سادگی در نتیجه نهایی وجود نداشت. اگر بیش از ۱۱۰۰۰ `step` برای یک `query_range` وجود داشته باشد، `Prometheus` کوئری را با یک `error` رد می‌کند. این برای جلوگیری از ارسال تصادفی کوئری‌های بسیار بزرگ به `Prometheus` است، مانند `step` ۱ ثانیه‌ای برای یک هفته. از آنجایی که مانیتورهایی با وضوح افقی بیش از ۱۱۰۰۰ پیکسل (`pixels`) نادر هستند، بعید است هنگام `graphing` با این مشکل مواجه شوید. اگر در حال نوشتن اسکریپت‌های گزارش‌گیری هستید، می‌توانید درخواست‌های `query_range` را که به این حد می‌رسند، تقسیم کنید. این حد امکان وضوح یک دقیقه‌ای برای یک هفته، یا وضوح یک ساعته برای یک سال را فراهم می‌کند، بنابراین در بیشتر مواقع نباید اعمال شود.
**داده‌های تراز شده (Aligned data)**
هنگام استفاده از ابزارهایی مانند `Grafana`، معمول است که `alignment` (ترازبندی) `query_range` بر اساس زمان فعلی باشد، و بنابراین نتایج شما کاملاً با دقیقه‌ها، ساعت‌ها یا روزها هم‌تراز نخواهند بود. در حالی که این هنگام نگاه کردن به `dashboards` خوب است، به ندرت چیزی است که با اسکریپت‌های گزارش‌گیری می‌خواهید. `query_range` گزینه‌ای برای مشخص کردن `alignment` ندارد، در عوض این به شما بستگی دارد که پارامتر `start` را با `alignment` صحیح مشخص کنید. به عنوان مثال، اگر می‌خواستید در `Python` `samples` را هر ساعت سر ساعت داشته باشید، عبارت `(time.time() // 3600) * 3600` شروع ساعت فعلی  را برمی‌گرداند، که می‌توانید آن را در گام‌های ۳۶۰۰ تنظیم کرده و به عنوان پارامترهای `URL` `start` و `end` استفاده کنید و سپس از پارامتر `step` با مقدار 3600 استفاده کنید.
اکنون که اصول اولیه نحوه استفاده از `PromQL` و اجرای کوئری‌ها از طریق `HTTP APIs` را می‌دانید، با جزئیات بیشتری به `aggregation` خواهیم پرداخت.