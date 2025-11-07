
در این‌جا ترجمهٔ صفحهٔ «روندها/Recording Rule» در مستندات Prometheus را برایتان آورده‌ام:

---

### Recording Rule

یک اسکِم نام‌گذاری ثابت برای Recording Rule باعث می‌شود که معنیِ یک قاعده به‌راحتی در نگاه اول قابل فهم باشد. همچنین از بروز اشتباهات، مثلاً محاسبات نادرست یا بی‌معنی جلوگیری می‌کند.
در این صفحه، نام‌گذاری صحیح و گردآوری (Aggregation) برای Recording Rule توضیح داده شده است.


---

#### نام‌گذاری

* Recording Rule باید به شکل کلی زیر باشند:
  `level:metric:operations` 
* ‏`level` نشان‌دهندهٔ سطحِ گردآوری و برچسب‌های خروجیِ قاعده است. 
* ‏`metric` اسم متریک است و باید بدون تغییر باقی بماند، به‌جز این‌که اگر یک کانتر (`_total`) را با `rate()` یا `irate()` استفاده می‌کنید، آن قسمت `_total` را حذف می‌کنید. 
* ‏`operations` فهرستی از عملیاتی است که بر روی متریک انجام شده‌اند، که جدیدترین عملیات اول نوشته می‌شوند.
نگه‌داشتن اسم متریک بدون تغییر باعث می‌شود که تشخیص آن آسان‌تر شود و یافتنش در کد یا محاسبات راحت‌تر باشد. 
برای نگه‌داشتن عملیات واضح، اگر عملیاتی `sum()` است و متریک خودش `_sum` ندارد، آن را حذف می‌کنیم. مثلاً اگر `min_min` باشد، همان `min` حساب می‌شود. 
اگر عملیاتی واضح نیست، از `sum` استفاده کنید. زمانی که بخواهید نسبت بگیرید (division) بین متریک‌ها، از `_per_` برای جداسازی اسم متریک‌ها استفاده کنید و عملیات را `ratio` نام‌گذاری کنید. 


---

#### گردآوری (Aggregation)

* وقتی نسبت‌ها را گردآوری می‌کنید، ابتدا صورت و مخرج را جداگانه جمع بزنید، سپس تقسیم کنید.
* از گرفتن میانگین نسبت یا میانگین میانگین خودداری کنید، چون از لحاظ آماری معتبر نیست. 
* در زمانی که با Summary سروکار دارید (مثلاً `_count` و `_sum`) و می‌خواهید میانگین مشاهدات را حساب کنید، اسم متریک را بدون پسوند `_count` یا `_sum` نگه‌دارید و `rate` را با `mean` جایگزین کنید؛ چون این کار نمایانگر میانگین اندازهٔ مشاهده‌ها در دورهٔ زمانی است. 
* همیشه از عبارت `without` با برچسب‌هایی که در حال گردآوری‌شدن هستند استفاده کنید. این کار باعث می‌شود که سایر برچسب‌ها مانند `job` حفظ شوند، و از بروز تداخل‌ها جلوگیری شود و متریک‌ها و هشدارها مفیدتر شوند. 

---

#### مثال‌ها

در ادامه چند مثال آمده است: 

* گردآوری نرخ درخواست‌ها در ثانیه با برچسب `path`:

  ```yaml
  - record: instance_path:requests:rate5m
    expr: rate(requests_total{job="myjob"}[5m])

  - record: path:requests:rate5m
    expr: sum without (instance)(instance_path:requests:rate5m{job="myjob"})
  ```
* محاسبهٔ نسبت خطاهای درخواست:

  ```yaml
  - record: instance_path:request_failures:rate5m
    expr: rate(request_failures_total{job="myjob"}[5m])

  - record: instance_path:request_failures_per_requests:ratio_rate5m
    expr: |
      instance_path:request_failures:rate5m{job="myjob"}
    /
      instance_path:requests:rate5m{job="myjob"}

  - record: path:request_failures_per_requests:ratio_rate5m
    expr: |
      sum without (instance)(instance_path:request_failures:rate5m{job="myjob"})
    /
      sum without (instance)(instance_path:requests:rate5m{job="myjob"})
  ```

  و سپس گردآوری تا سطح job:

  ```yaml
  - record: job:request_failures_per_requests:ratio_rate5m
    expr: |
      sum without (instance, path)(instance_path:request_failures:rate5m{job="myjob"})
    /
      sum without (instance, path)(instance_path:requests:rate5m{job="myjob"})
  ```
* محاسبهٔ میانگین تاخیر از یک Summary:

  ```yaml
  - record: instance_path:request_latency_seconds_count:rate5m
    expr: rate(request_latency_seconds_count{job="myjob"}[5m])

  - record: instance_path:request_latency_seconds_sum:rate5m
    expr: rate(request_latency_seconds_sum{job="myjob"}[5m])

  - record: instance_path:request_latency_seconds:mean5m
    expr: |
      instance_path:request_latency_seconds_sum:rate5m{job="myjob"}
    /
      instance_path:request_latency_seconds_count:rate5m{job="myjob"}

  - record: path:request_latency_seconds:mean5m
    expr: |
      sum without (instance)(instance_path:request_latency_seconds_sum:rate5m{job="myjob"})
    /
      sum without (instance)(instance_path:request_latency_seconds_count:rate5m{job="myjob"})
  ```
* محاسبهٔ نرخ میانگین پرس‌وجوها در طول نمونه‌ها و مسیرها با استفاده از `avg()`:

  ```yaml
  - record: job:request_latency_seconds_count:avg_rate5m
    expr: avg without (instance, path)(instance_path:request_latency_seconds_count:rate5m{job="myjob"})
  ```

نکته: وقتی گردآوری می‌کنید، برچسب‌هایی که در `without` فهرست شده‌اند، از سطح اسم متریک خروجی حذف می‌شوند؛ اگر درجۀ خروجی با درجۀ ورودی مطابقت ندارد، احتمالاً در قاعده اشتباهی رخ داده است.


---



---

### تعریف Recording Rule

#### پیکربندی قواعد

Prometheus از دو نوع قاعده پشتیبانی می‌کند که می‌توان آن‌ها را پیکربندی کرده و در فواصل منظم ارزیابی نمود: **Recording Rule** (recording rules) و **قواعد هشداردهی** (alerting rules). 
برای اینکه قواعد را در Prometheus وارد کنید، یک فایل حاوی دستورهای مربوطه بسازید و آن را از طریق فیلد `rule_files` در پیکربندی Prometheus بارگذاری کنید. فایل‌های قاعده باید به فرمت YAML باشند. 
فایل‌های قاعده در زمان اجرا نیز قابل بازخوانی هستند؛ کافی است به فرایند Prometheus سیگنال `SIGHUP` ارسال شود. تغییرات فقط زمانی اعمال می‌شوند که همه فایل‌های قاعده، درست فرمت شده باشند. 

#### بررسی سینتکس فایل قواعد

برای اینکه سریع بررسی کنید آیا فایل قواعد از نظر نحوی درست است بدون اینکه سرور Prometheus را راه‌اندازی کنید، می‌توانید از ابزار خط فرمان `promtool` استفاده نمایید: 

```
promtool check rules /path/to/example.rules.yml
```

اگر فایل از نظر سینتکس معتبر باشد، ابزار بازنمایی متنی قواعد را چاپ می‌کند و با کد خروجی 0 خاتمه می‌یابد. اگر خطا وجود داشته باشد، پیام خطا به stderr چاپ می‌شود و با کد خروجی 1 پایان می‌یابد. 

#### Recording Rule

Recording Rule به شما اجازه می‌دهند که عبارات پرس‌و‌جو (PromQL) که به‌طور مکرر مورد نیازند یا محاسباتشان سنگین است را از پیش محاسبه کنید و نتیجه را به‌صورت یک مجموعهٔ جدید از سری‌های زمانی ذخیره نمایید. سپس پرس‌وجو بر روی نتیجهٔ از پیش محاسبه‌شده معمولاً بسیار سریع‌تر خواهد بود نسبت به اجرای عبارت اصلی هر بار. این موضوع به‌ویژه برای داشبوردهایی مفید است که با هر تازه‌سازی همان عبارت را مجدداً اجرا می‌کنند. 
Recording Rule و قواعد هشداردهی هر دو درون «گروه قاعده» (rule group) قرار می‌گیرند. قواعد درون یک گروه به‌صورت متوالی اجرا می‌شوند و زمان ارزیابی یکسانی دارند. 
نام Recording Rule باید **یک نام معتبر برای متریک‌ها** باشد. 

#### ساختار فایل قاعده

ساختار کلی فایل YAML قاعده به این شکل است: 

```yaml
groups:
  - name: <rule_group>
    rules:
      - record: <string>
        expr: <string>
        labels:
          <labelname>: <labelvalue>
```

و همچنین پارامترهای اختیاری برای گروه‌ها: 

* `interval: <duration>` | پیش‌فرض = `global.evaluation_interval`
* `limit: <int>` | پیش‌فرض = 0 — محدودیت تعداد سر‌ی‌های تولیدشده توسط Recording Rule یا هشداردهی
* `query_offset: <duration>` | پیش‌فرض = `global.rule_query_offset` — افست زمانی برای پرس‌وجو
* `labels:` — برچسب‌هایی که قبل از ذخیرهٔ نتیجه اضافه یا بازنویسی می‌شوند.

#### مثال

یک مثال ساده از فایل قواعد: 

```yaml
groups:
  - name: example
    rules:
      - record: code:prometheus_http_requests_total:sum
        expr: sum by (code) (prometheus_http_requests_total)
```

#### محدود کردن هشدارها و سری‌ها

محدودیتی برای تعداد هشدارهایی که قواعد هشداردهی تولید می‌کنند و سری‌هایی که Recording Rule تولید می‌کنند قابل تنظیم است. زمانی که محدودیت تجاوز شود، تمام سری تولیدشده توسط آن قاعده حذف خواهند شد، و در مورد قواعد هشداردهی، تمام هشدارها (فعال، در انتظار، یا غیرفعال) نیز پاک می‌شوند. این رویداد به‌عنوان یک خطا در ارزیابی ثبت می‌شود.

#### افست پرس‌وجوی قاعده (Rule Query Offset)

این گزینه زمانی مفید است که بخواهید مطمئن شوید متریک‌های زمینه‌ای قبل از اجرای قاعده دریافت و ذخیره شده‌اند. تأخیر در در دسترس شدن متریک‌ها بیشتر زمانی رخ می‌دهد که Prometheus به‌عنوان یک هدف “remote write” اجرا شود یا زمانی که اسکرِیپ یا ارزیابی به‌صورت غیرمعمولی دیر انجام شود.

#### ارزیابی ناموفق قواعد به‌خاطر ارزیابی کند

اگر گروه قاعده، ارزیابی خود را قبل از زمان بعدی ارزیابی (بر اساس `evaluation_interval`) به پایان نرساند، آن ارزیابی بعدی رد می‌شود. ارزیابی‌های بعدی گروه قاعده نیز تا زمانی که ارزیابی فعلی یا به پایان برسد یا تایم‌اوت شود، رد خواهند شد. در این حالت شکاف (gap) در متریک تولیدشده توسط قاعده ضبط دیده خواهد شد. متریک `rule_group_iterations_missed_total` برای هر اجرای از دست‌رفته افزایش می‌یابد. 

[1]: https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/ "Defining recording rules | Prometheus"
