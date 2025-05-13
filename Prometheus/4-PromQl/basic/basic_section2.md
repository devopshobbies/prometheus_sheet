

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

این قابل پیش‌بینی بودن در مورد اینکه کدام `labels` بازگردانده می‌شوند، برای تطبیق `vector` (`vector matching`) با `operators` مهم است، همانطور که دردر گذشته بحث خواهد شد.

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

## Summary

یک `metric` از نوع `summary` معمولاً شامل `_sum` و `_count` و گاهی اوقات یک `time series` بدون پسوند با یک `label` به نام `quantile` است. `_sum` و `_count` هر دو `counter` هستند.
‏`Prometheus` شما یک `summary` به نام `http_response_size_bytes` برای مقدار داده‌ای که برخی از `HTTP APIs` آن بازمی‌گردانند، ارائه می‌دهد.² `http_response_size_bytes_count` تعداد درخواست‌ها را ردیابی می‌کند و از آنجایی که یک `counter` است، باید قبل از تجمیع `label` با نام `handler` آن، از `rate` استفاده کنید:

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

کار با تمام `time series`های مختلف با مقادیر `label` متفاوت برای یک `metric` می‌تواند کمی طاقت‌فرسا و بالقوه گیج‌کننده باشد اگر `metric` از انواع مختلف سرورها⁵ می‌آید. معمولاً می‌خواهید `time series`هایی را که روی آن‌ها کار می‌کنید محدود کنید. تقریباً همیشه می‌خواهید بر اساس `label` با نام `job` محدود کنید و بسته به کاری که انجام می‌دهید، ممکن است بخواهید فقط به یک `instance` یا یک `handler` نگاه کنید، به عنوان مثال.

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

`topk(5, rate(http_requests_total[1h] @ end()))` به عنوان یک تابع رتبه‌بندی (`ranking function`) عمل می‌کند و فقط مقادیر بالاتر را در انتهای بازه ارزیابی فیلتر می‌کند.

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
شما همچنین می‌توانید از `range vectors` با `endpoint` `query` استفاده کنید. به عنوان مثال، `http://localhost:9090/api/v1/query?query=prometheus_tsdb_head_samples_appended_total[1m]` نتایجی مانند این را برمی‌گرداند:

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

این با نتیجه `instant vector` قبلی متفاوت است، زیرا `resultType` اکنون `matrix` است و هر `time series` چندین مقدار دارد. هنگامی که با `range vector` استفاده می‌شود، `endpoint` `query` `samples` خام را برمی‌گرداند،¹⁶ اما مراقب باشید که داده‌های زیادی را به یکباره درخواست نکنید زیرا ممکن است حافظه یک طرف یا طرف دیگر تمام شود.
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

‏`endpoint` با نام `query range` در `/api/v1/query_range` `endpoint` اصلی `HTTP` `Prometheus` است که از آن استفاده خواهید کرد، زیرا این `endpoint` برای `graphing` (رسم نمودار) است. در پشت صحنه، `query_range` یک `syntactic sugar` (به علاوه برخی بهینه‌سازی‌های عملکردی) برای چندین فراخوانی به `endpoint` `query` است. علاوه بر پارامتر `URL` `query`، شما به `query_range` یک `start time` (زمان شروع)، یک `end time` (زمان پایان) و یک `step` (گام) ارائه می‌دهید. کوئری ابتدا در `start time` اجرا می‌شود. سپس `step` ثانیه پس از `start time` اجرا می‌شود. سپس دو برابر `step` ثانیه پس از `start time` اجرا می‌شود و به همین ترتیب، تا زمانی که زمان ارزیابی کوئری از `end time` تجاوز کند، متوقف می‌شود. تمام نتایج `instant vector`¹⁸ از اجراهای مختلف در یک `range vector` ترکیب شده و بازگردانده می‌شوند.
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

چند جنبه از این وجود دارد که باید به آنها توجه کنید. اول اینکه `timestamps` `sample`ها با `start time` و `step` هم‌تراز هستند، زیرا هر نتیجه از یک ارزیابی کوئری لحظه‌ای متفاوت می‌آید و نتایج کوئری لحظه‌ای همیشه از زمان ارزیابی خود به عنوان `timestamp` نتایج استفاده می‌کنند. دوم اینکه آخرین `sample` در اینجا در `end time` است، به این معنی که بازه `inclusive` (شامل) است و آخرین نقطه `end time` خواهد بود اگر اتفاقاً با `step` هم‌تراز شود. سوم اینکه ما یک بازه ۵ دقیقه‌ای برای تابع `rate` انتخاب کردیم که بزرگتر از `step` است. از آنجایی که `query_range` ارزیابی‌های کوئری لحظه‌ای تکراری را انجام می‌دهد، هیچ وضعیتی بین ارزیابی‌ها منتقل نمی‌شود. اگر بازه کوچکتر از `step` بود، ما از روی داده‌ها می‌پریدیم. به عنوان مثال، یک بازه ۱ دقیقه‌ای با `step` ۵ دقیقه‌ای، ۸۰٪ از `samples` را نادیده می‌گرفت. برای جلوگیری از این، باید از بازه‌هایی استفاده کنید که حداقل یک یا دو فاصله `scrape` (`scrape intervals`) بزرگتر از `step`ی باشند که استفاده می‌کنید.

هنگام استفاده از `range vectors` با `query_range`، معمولاً باید از بازه‌ای استفاده کنید که طولانی‌تر از `step` شما باشد تا داده‌ای را از دست ندهید.

چهارم اینکه برخی از `samples` به طور خاص گرد نیستند، و اینکه هر عددی اصلاً گرد باشد به این دلیل است که این یک تنظیم ساده از مقادیر `sample` است. هنگام کار با `metrics`، داده‌های شما به ندرت کاملاً تمیز هستند؛ `targets` مختلف در زمان‌های مختلف `scrape` می‌شوند و `scrapes` می‌توانند به تأخیر بیفتند. هنگام انجام کوئری‌هایی که کاملاً با داده‌های زیربنایی هم‌تراز نیستند یا هنگام تجمیع در چندین میزبان، به ندرت نتایج گردی دریافت خواهید کرد. علاوه بر این، ماهیت محاسبات `floating-point` می‌تواند منجر به اعدادی شود که تقریباً گرد هستند.

در اینجا، برای هر `step` یک `sample` وجود دارد. اگر اتفاق می‌افتاد که برای یک `time series` معین برای یک `step` نتیجه‌ای وجود نداشت، آن `sample` به سادگی در نتیجه نهایی وجود نداشت. اگر بیش از ۱۱۰۰۰ `step` برای یک `query_range` وجود داشته باشد، `Prometheus` کوئری را با یک `error` رد می‌کند. این برای جلوگیری از ارسال تصادفی کوئری‌های بسیار بزرگ به `Prometheus` است، مانند `step` ۱ ثانیه‌ای برای یک هفته. از آنجایی که مانیتورهایی با وضوح افقی بیش از ۱۱۰۰۰ پیکسل (`pixels`) نادر هستند، بعید است هنگام `graphing` با این مشکل مواجه شوید. اگر در حال نوشتن اسکریپت‌های گزارش‌گیری هستید، می‌توانید درخواست‌های `query_range` را که به این حد می‌رسند، تقسیم کنید. این حد امکان وضوح یک دقیقه‌ای برای یک هفته، یا وضوح یک ساعته برای یک سال را فراهم می‌کند، بنابراین در بیشتر مواقع نباید اعمال شود.
**داده‌های تراز شده (Aligned data)**
هنگام استفاده از ابزارهایی مانند `Grafana`، معمول است که `alignment` (ترازبندی) `query_range` بر اساس زمان فعلی باشد، و بنابراین نتایج شما کاملاً با دقیقه‌ها، ساعت‌ها یا روزها هم‌تراز نخواهند بود. در حالی که این هنگام نگاه کردن به `dashboards` خوب است، به ندرت چیزی است که با اسکریپت‌های گزارش‌گیری می‌خواهید. `query_range` گزینه‌ای برای مشخص کردن `alignment` ندارد، در عوض این به شما بستگی دارد که پارامتر `start` را با `alignment` صحیح مشخص کنید. به عنوان مثال، اگر می‌خواستید در `Python` `samples` را هر ساعت سر ساعت داشته باشید، عبارت `(time.time() // 3600) * 3600` شروع ساعت فعلی¹⁹ را برمی‌گرداند، که می‌توانید آن را در گام‌های ۳۶۰۰ تنظیم کرده و به عنوان پارامترهای `URL` `start` و `end` استفاده کنید و سپس از پارامتر `step` با مقدار 3600 استفاده کنید.
اکنون که اصول اولیه نحوه استفاده از `PromQL` و اجرای کوئری‌ها از طریق `HTTP APIs` را می‌دانید، با جزئیات بیشتری به `aggregation` خواهیم پرداخت.