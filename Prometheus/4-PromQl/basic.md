

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


![[ReangeVec.png]]


![[RangeVector2.png]]

بسته به مورد استفاده (مثلاً هنگام ترسیم `graph` در مقابل نمایش خروجی یک `expression`)، فقط برخی از این انواع به عنوان نتیجه یک `expression` مشخص شده توسط کاربر مجاز هستند. برای `instant queries`، هر یک از انواع داده بالا به عنوان ریشه `expression` مجاز است. `Range queries` فقط از `expression`های نوع `scalar` و `instant-vector` پشتیبانی می‌کنند.

**یادداشت‌هایی درباره `native histograms` آزمایشی:**

*   فعال کردن `Ingesting` `native histograms` باید از طریق یک `feature flag` انجام شود.
*   هنگامی که `native histograms` در `TSDB` `ingest` شدند (و حتی پس از غیرفعال کردن دوباره `feature flag`)، هم `instant vectors` و هم `range vectors` ممکن است اکنون حاوی `sample`هایی باشند که اعداد ساده `floating point` (`float samples`) نیستند، بلکه `histograms` کامل (`histogram samples`) هستند. یک `vector` ممکن است ترکیبی از `float samples` و `histogram samples` را شامل شود. توجه داشته باشید که اصطلاح «`histogram sample`» در مستندات `PromQL` همیشه به یک `native histogram` اشاره دارد. `Classic histograms` به تعدادی `series` از `float samples` تقسیم می‌شوند. از دیدگاه `PromQL`، چیزی به نام «`classic histogram samples`» وجود ندارد.
*   مانند `float samples`، `histogram samples` می‌توانند `counters` یا `gauges` باشند که به ترتیب `counter histograms` یا `gauge histograms` نیز نامیده می‌شوند.
*  ‏ `Native histograms` می‌توانند `bucket layout`های مختلفی داشته باشند، اما به طور کلی برای اعمال عملیات `binary` و `aggregation` روی آنها، قابل تبدیل به نسخه‌های سازگار هستند. این برای همه `bucketing schemas` صادق نیست. اگر `histograms` ناسازگار در یک عملیات مواجه شوند، عنصر `vector` خروجی مربوطه از نتیجه حذف شده و با یک `annotation` سطح هشدار (warn-level) علامت‌گذاری می‌شود. جزئیات بیشتر را می‌توان در `native histogram specification` یافت.

**`Literals`**

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

`Float literals` همچنین برای مشخص کردن `durations` بر حسب ثانیه استفاده می‌شوند. برای راحتی، اعداد `integer` `decimal` ممکن است با `time units` زیر ترکیب شوند:

*   `ms` – milliseconds
*   `s` – seconds – 1s برابر است با 1000ms
*   `m` – minutes – 1m برابر است با 60s (نادیده گرفتن leap seconds)
*   `h` – hours – 1h برابر است با 60m
*   `d` – days – 1d برابر است با 24h (نادیده گرفتن daylight saving time)
*   `w` – weeks – 1w برابر است با 7d
*   `y` – years – 1y برابر است با 365d (نادیده گرفتن leap days)

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

**`Instant vector selectors`**
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

تطابق‌های `Regex` کاملاً `anchored` هستند. تطابق `env=~"foo"` به عنوان `env=~"^foo$"` در نظر گرفته می‌شود.

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

`query` `http_requests_total{environment=""}` مطابقت داده و بازمی‌گرداند:

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

`Vector selectors` باید یا یک نام مشخص کنند یا حداقل یک `label matcher` که با `string` خالی مطابقت نداشته باشد. `expression` زیر غیرقانونی است:

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