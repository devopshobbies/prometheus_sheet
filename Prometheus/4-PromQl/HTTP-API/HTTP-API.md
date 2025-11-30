

---

### HTTP API

‏API پایدار فعلی HTTP از طریق `/api/v1` در سرور پرومتئوس قابل دسترسی است. هرگونه افزودنی که ساختار موجود را برهم نزند، تحت همین نقطه پایانی اضافه خواهد شد.

نمای کلی قالب (Format overview)

قالب پاسخ API به صورت JSON است. هر درخواست موفق API یک کد وضعیت 2xx برمی‌گرداند.

درخواست‌های نامعتبری که به کنترل‌کننده‌های API می‌رسند، یک شیء خطای JSON و یکی از کدهای وضعیت HTTP زیر را برمی‌گردانند:

*   ‏400 Bad Request: زمانی که پارامترها ناقص یا نادرست باشند.
*   ‏422 Unprocessable Entity: زمانی که یک عبارت قابل اجرا نباشد (RFC4918).
*   ‏503 Service Unavailable: زمانی که کوئری‌ها زمانشان تمام شود یا لغو شوند.

سایر کدهای غیر 2xx ممکن است برای خطاهایی که قبل از رسیدن به نقطه پایانی API رخ می‌دهند، بازگردانده شوند.

در صورت وجود خطاهایی که مانع اجرای درخواست نمی‌شوند، ممکن است آرایه‌ای از هشدارها (warnings) بازگردانده شود. همچنین، آرایه دیگری از یادداشت‌های سطح اطلاعاتی (info-level annotations) ممکن است برای مشکلات احتمالی کوئری‌ که ممکن است مثبت کاذب باشند یا نباشند، بازگردانده شود. تمام داده‌هایی که با موفقیت جمع‌آوری شده‌اند، در فیلد `data` بازگردانده خواهند شد.

قالب کلی پاسخ JSON به شرح زیر است:

```json
{
  "status": "success" | "error",
  "data": <data>,

  // فقط در صورتی تنظیم می‌شود که status "error" باشد.
  // فیلد data ممکن است همچنان حاوی داده‌های اضافی باشد.
  "errorType": "<string>",
  "error": "<string>",

  // فقط در صورتی تنظیم می‌شود که هنگام اجرای درخواست هشدارهایی وجود داشته باشد.
  // همچنان داده‌ها در فیلد data وجود خواهند داشت.
  "warnings": ["<string>"],
  // فقط در صورتی تنظیم می‌شود که هنگام اجرای درخواست یادداشت‌های سطح اطلاعاتی وجود داشته باشد.
  "infos": ["<string>"]
}
```

مقادیر نگه‌دارنده عمومی به شرح زیر تعریف می‌شوند:

 * ‏   <rfc3339 | unix_timestamp>: این timestamp ورودی می‌توانند به فرمت RFC3339 یا به عنوان timestamp یونیکس بر حسب ثانیه ارائه شوند، با امکان دقت زیرثانیه‌ای (با اعشار). timestamp خروجی همیشه به عنوان timestamp یونیکس بر حسب ثانیه نمایش داده می‌شوند.
* ‏ series_selector: انتخاب‌گرهای سری زمانی پرومتئوس مانند `http_requests_total` یا `http_requests_total{method=~"(GET|POST)"}` که باید کدگذاری URL شوند.
* ‏ duration: زیرمجموعه‌ای از مقادیر اعشاری پرومتئوس که از واحدهای زمانی استفاده می‌کنند. به عنوان مثال، `5m` به مدت زمان 5 دقیقه اشاره دارد.
* ‏bool: مقادیر بولی (رشته‌های "true" و "false").

توجه: نام پارامترهای کوئری‌ که ممکن است تکرار شوند، با `[]` به پایان می‌رسند.

#### کوئری‌های عبارتی (Expression queries)

عبارات زبان کوئری‌ را می‌توان در یک لحظه مشخص یا در یک بازه زمانی ارزیابی کرد. بخش‌های زیر، نقاط پایانی API را برای هر نوع کوئری عبارتی توضیح می‌دهند.

کوئری‌های لحظه‌ای (Instant queries)

نقطه پایانی زیر، یک کوئری لحظه‌ای را در یک نقطه زمانی خاص ارزیابی می‌کند:

`GET /api/v1/query`
`POST /api/v1/query`

پارامترهای کوئری URL:

‏ *  <query=<string: رشته کوئری عبارت پرومتئوس.
*‏   time=rfc3339 | unix_timestamp: مهر زمانی ارزیابی. اختیاری.
*‏   timeout=duration: مهلت ارزیابی. اختیاری. به طور پیش‌فرض، مقدار پرچم `-query.timeout` است و توسط آن محدود می‌شود.
* ‏  limit=number: حداکثر تعداد سری‌های برگشتی. روی مقادیر اسکالر یا رشته‌ها تأثیری ندارد اما تعداد سری‌ها را برای ماتریس‌ها و بردارها کوتاه می‌کند. اختیاری. `0` به معنای غیرفعال بودن است.
*  ‏ lookback_delta=number: دوره بازنگری (lookback period) را فقط برای این کوئری‌ بازنویسی می‌کند. اختیاری.

اگر پارامتر `time` حذف شود، زمان فعلی سرور استفاده می‌شود.



می‌توانید این پارامترها را مستقیماً در بدنه درخواست با استفاده از متد POST و هدر `Content-Type: application/x-www-form-urlencoded` کدگذاری URL کنید. این کار زمانی مفید است که یک کوئری بزرگ را مشخص می‌کنید که ممکن است از محدودیت کاراکتر URL سمت سرور فراتر رود.

بخش `data` نتیجه کوئری‌ دارای قالب زیر است:

```json
{
  "resultType": "matrix" | "vector" | "scalar" | "string",
  "result": <value>
}
```

‏`<value>` به داده‌های نتیجه کوئری‌ اشاره دارد که بسته به `resultType` قالب‌های متفاوتی دارد. به فرمت‌های نتیجه کوئری عبارت مراجعه کنید.

مثال زیر عبارت `up` را در زمان `2015-07-01T20:10:51.781Z` ارزیابی می‌کند:

`curl 'http://localhost:9090/api/v1/query?query=up&time=2015-07-01T20:10:51.781Z'`
```json
{
   "status" : "success",
   "data" : {
      "resultType" : "vector",
      "result" : [
         {
            "metric" : {
               "__name__" : "up",
               "job" : "prometheus",
               "instance" : "localhost:9090"
            },
            "value": [ 1435781451.781, "1" ]
         },
         {
            "metric" : {
               "__name__" : "up",
               "job" : "node",
               "instance" : "localhost:9100"
            },
            "value" : [ 1435781451.781, "0" ]
         }
      ]
   }
}
```

کوئری‌های بازه‌ای (Range queries)

نقطه پایانی زیر یک کوئری عبارتی را در یک بازه زمانی ارزیابی می‌کند:

`GET /api/v1/query_range`
`POST /api/v1/query_range`

پارامترهای کوئری URL:

*‏   query=string : رشته کوئری عبارت پرومتئوس.
*‏   start=rfc3339 | unix_timestamp: مهر زمانی شروع، شامل.
*‏   end=rfc3339 | unix_timestamp: مهر زمانی پایان، شامل.
* ‏  step=duration | float: عرض گام (step width) وضوح کوئری‌ در قالب مدت زمان یا عدد اعشاری بر حسب ثانیه.
*‏   timeout=duration: مهلت ارزیابی. اختیاری. به طور پیش‌فرض، مقدار پرچم `-query.timeout` است و توسط آن محدود می‌شود.
* ‏  limit=number: حداکثر تعداد سری‌های برگشتی. اختیاری. `0` به معنای غیرفعال بودن است.
* ‏  lookback_delta=number: دوره بازنگری (lookback period) را فقط برای این کوئری‌ بازنویسی می‌کند. اختیاری.

می‌توانید این پارامترها را مستقیماً در بدنه درخواست با استفاده از متد POST و هدر `Content-Type: application/x-www-form-urlencoded` کدگذاری URL کنید. این کار زمانی مفید است که یک کوئری بزرگ را مشخص می‌کنید که ممکن است از محدودیت کاراکتر URL سمت سرور فراتر رود.

بخش `data` نتیجه کوئری‌ دارای قالب زیر است:

```json
{
  "resultType": "matrix",
  "result": <value>
}
```

برای قالب نگه‌دارنده `<value>`، به فرمت نتیجه بردار بازه‌ای مراجعه کنید.

مثال زیر عبارت `up` را در یک بازه 30 ثانیه‌ای با وضوح کوئری 15 ثانیه ارزیابی می‌کند.

`curl 'http://localhost:9090/api/v1/query_range?query=up&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s'`
```json
{
   "status" : "success",
   "data" : {
      "resultType" : "matrix",
      "result" : [
         {
            "metric" : {
               "__name__" : "up",
               "job" : "prometheus",
               "instance" : "localhost:9090"
            },
            "values" : [
               [ 1435781430.781, "1" ],
               [ 1435781445.781, "1" ],
               [ 1435781460.781, "1" ]
            ]
         },
         {
            "metric" : {
               "__name__" : "up",
               "job" : "node",
               "instance" : "localhost:9091"
            },
            "values" : [
               [ 1435781430.781, "0" ],
               [ 1435781445.781, "0" ],
               [ 1435781460.781, "1" ]
            ]
         }
      ]
   }
}
```

قالب‌بندی عبارات کوئری‌ (Formatting query expressions)

نقطه پایانی زیر یک عبارت PromQL را به شیوه‌ای زیبا قالب‌بندی می‌کند:

`GET /api/v1/format_query`
`POST /api/v1/format_query`

پارامترهای کوئری URL:

*   ‏query=string: رشته کوئری عبارت پرومتئوس.

می‌توانید این پارامترها را مستقیماً در بدنه درخواست با استفاده از متد POST و هدر `Content-Type: application/x-www-form-urlencoded` کدگذاری URL کنید. این کار زمانی مفید است که یک کوئری بزرگ را مشخص می‌کنید که ممکن است از محدودیت کاراکتر URL سمت سرور فراتر رود.

بخش `data` نتیجه کوئری‌ یک رشته حاوی عبارت کوئری قالب‌بندی شده است. توجه داشته باشید که هرگونه نظر در رشته قالب‌بندی شده حذف می‌شود.

مثال زیر عبارت `foo/bar` را قالب‌بندی می‌کند:

`curl 'http://localhost:9090/api/v1/format_query?query=foo/bar'`


```json
{
   "status" : "success",
   "data" : "foo / bar"
}
```

تجزیه عبارات PromQL به درخت نحو انتزاعی (AST) (Parsing a PromQL expressions into a abstract syntax tree (AST))

این نقطه پایانی آزمایشی است و ممکن است در آینده تغییر کند. در حال حاضر، فقط برای استفاده توسط رابط کاربری وب خود پرومتئوس در نظر گرفته شده است و نام نقطه پایانی و قالب دقیق بازگردانده شده ممکن است از یک نسخه پرومتئوس به نسخه دیگر تغییر کند. همچنین ممکن است در صورت عدم نیاز توسط رابط کاربری، دوباره حذف شود.

نقطه پایانی زیر یک عبارت PromQL را تجزیه کرده و آن را به عنوان یک نمایش AST (درخت نحو انتزاعی) با فرمت JSON بازمی‌گرداند:

`GET /api/v1/parse_query`
`POST /api/v1/parse_query`

پارامترهای کوئری URL:

*   ‏query=string: رشته کوئری عبارت پرومتئوس.

می‌توانید این پارامترها را مستقیماً در بدنه درخواست با استفاده از متد POST و هدر `Content-Type: application/x-www-form-urlencoded` کدگذاری URL کنید. این کار زمانی مفید است که یک کوئری بزرگ را مشخص می‌کنید که ممکن است از محدودیت کاراکتر URL سمت سرور فراتر رود.

بخش `data` نتیجه کوئری‌ یک رشته حاوی AST عبارت کوئری تجزیه شده است.

مثال زیر عبارت `foo/bar` را تجزیه می‌کند:

`curl 'http://localhost:9090/api/v1/parse_query?query=foo/bar'`
```json
{
   "data" : {
      "bool" : false,
      "lhs" : {
         "matchers" : [
            {
               "name" : "__name__",
               "type" : "=",
               "value" : "foo"
            }
         ],
         "name" : "foo",
         "offset" : 0,
         "startOrEnd" : null,
         "timestamp" : null,
         "type" : "vectorSelector"
      },
      "matching" : {
         "card" : "one-to-one",
         "include" : [],
         "labels" : [],
         "on" : false
      },
      "op" : "/",
      "rhs" : {
         "matchers" : [
            {
               "name" : "__name__",
               "type" : "=",
               "value" : "bar"
            }
         ],
         "name" : "bar",
         "offset" : 0,
         "startOrEnd" : null,
         "timestamp" : null,
         "type" : "vectorSelector"
      },
      "type" : "binaryExpr"
   },
   "status" : "success"
}
```

### کوئری فراداده (Querying metadata)

پرومتئوس مجموعه‌ای از نقاط پایانی API را برای کوئری فراداده در مورد سری‌ها و برچسب‌های آن‌ها ارائه می‌دهد.

نکته: این نقاط پایانی API ممکن است فراداده‌هایی را برای سری‌هایی بازگردانند که در محدوده زمانی انتخاب شده هیچ نمونه‌ای ندارند، و/یا برای سری‌هایی که نمونه‌های آن‌ها از طریق نقطه پایانی API حذف، علامت‌گذاری شده‌اند. میزان دقیق فراداده‌های اضافی بازگردانده شده، جزئیات پیاده‌سازی است که ممکن است در آینده تغییر کند.

### یافتن سری‌ها بر اساس تطابق‌دهنده‌های برچسب (Finding series by label matchers)

نقطه پایانی زیر فهرستی از سری‌های زمانی را که با یک مجموعه برچسب خاص مطابقت دارند، بازمی‌گرداند.

`GET /api/v1/series`
`POST /api/v1/series`

پارامترهای کوئری URL:

* ‏  match[]=<series_selector>: آرگومان انتخاب‌گر سری تکراری که سری‌های مورد نظر را برای بازگرداندن انتخاب می‌کند. حداقل یک آرگومان `match[]` باید ارائه شود.
* ‏  start=<rfc3339 | unix_timestamp>: مهر زمانی شروع.
*   ‏end=<rfc3339 | unix_timestamp>: مهر زمانی پایان.
* ‏  limit=number: حداکثر تعداد سری‌های برگشتی. اختیاری. `0` به معنای غیرفعال بودن است.

می‌توانید این پارامترها را مستقیماً در بدنه درخواست با استفاده از متد POST و هدر `Content-Type: application/x-www-form-urlencoded` کدگذاری URL کنید. این کار زمانی مفید است که تعداد زیادی یا پویایی از انتخاب‌گرهای سری را مشخص می‌کنید که ممکن است از محدودیت کاراکتر URL سمت سرور فراتر رود.

بخش `data` نتیجه کوئری‌ شامل لیستی از اشیا است که شامل جفت‌های نام/مقدار برچسب هستند که هر سری را شناسایی می‌کنند. توجه داشته باشید که زمان‌های شروع و پایان تقریبی هستند و نتیجه ممکن است شامل مقادیر برچسب برای سری‌هایی باشد که در بازه زمانی مشخص شده هیچ نمونه‌ای ندارند.

مثال زیر تمام سری‌هایی را که با هر یک از انتخاب‌گرهای `up` یا `process_start_time_seconds{job="prometheus"}` مطابقت دارند، بازمی‌گرداند:

`curl -g 'http://localhost:9090/api/v1/series?' --data-urlencode 'match[]=up' --data-urlencode 'match[]=process_start_time_seconds{job="prometheus"}'`

```json
{
   "status" : "success",
   "data" : [
      {
         "__name__" : "up",
         "job" : "prometheus",
         "instance" : "localhost:9090"
      },
      {
         "__name__" : "up",
         "job" : "node",
         "instance" : "localhost:9091"
      },
      {
         "__name__" : "process_start_time_seconds",
         "job" : "prometheus",
         "instance" : "localhost:9090"
      }
   ]
}
```

### دریافت نام برچسب‌ها (Getting label names)

نقطه پایانی زیر لیستی از نام برچسب‌ها را بازمی‌گرداند:

`GET /api/v1/labels`
`POST /api/v1/labels`

پارامترهای کوئری URL:

*   ‏start=<rfc3339 | unix_timestamp>: مهر زمانی شروع. اختیاری.
*   ‏end=<rfc3339 | unix_timestamp>: مهر زمانی پایان. اختیاری.
*  ‏ match[]=<series_selector>: آرگومان انتخاب‌گر سری تکراری که سری‌ها را برای خواندن نام برچسب‌ها انتخاب می‌کند. اختیاری.
* ‏  limit=number: حداکثر تعداد سری‌های برگشتی. اختیاری. `0` به معنای غیرفعال بودن است.

بخش `data` پاسخ JSON لیستی از نام‌های برچسب رشته‌ای است. توجه داشته باشید که زمان‌های شروع و پایان تقریبی هستند و نتیجه ممکن است شامل نام‌های برچسب برای سری‌هایی باشد که در بازه زمانی مشخص شده هیچ نمونه‌ای ندارند.

در اینجا یک مثال آورده شده است.

`curl 'localhost:9090/api/v1/labels'`
```json
{
    "status": "success",
    "data": [
        "__name__",
        "call",
        "code",
        "config",
        "dialer_name",
        "endpoint",
        "event",
        "goversion",
        "handler",
        "instance",
        "interval",
        "job",
        "le",
        "listener_name",
        "name",
        "quantile",
        "reason",
        "role",
        "scrape_job",
        "slice",
        "version"
    ]
}
```

### کوئری مقادیر برچسب (Querying label values)

نقطه پایانی زیر لیستی از مقادیر برچسب را برای یک نام برچسب ارائه شده بازمی‌گرداند:

`GET /api/v1/label/<label_name>/values`

پارامترهای کوئری URL:

*   ‏start=<rfc3339 | unix_timestamp>: مهر زمانی شروع. اختیاری.
*   ‏end=<rfc3339 | unix_timestamp>: مهر زمانی پایان. اختیاری.
*  ‏ match[]=<series_selector>: آرگومان انتخاب‌گر سری تکراری که سری‌ها را برای خواندن مقادیر برچسب انتخاب می‌کند. اختیاری.
*  ‏ limit=number: حداکثر تعداد سری‌های برگشتی. اختیاری. `0` به معنای غیرفعال بودن است.

بخش `data` پاسخ JSON لیستی از مقادیر برچسب رشته‌ای است. توجه داشته باشید که زمان‌های شروع و پایان تقریبی هستند و نتیجه ممکن است شامل مقادیر برچسب برای سری‌هایی باشد که در بازه زمانی مشخص شده هیچ نمونه‌ای ندارند.

این مثال تمام مقادیر برچسب را برای برچسب `http_status_code` کوئری‌ می‌کند:

`curl http://localhost:9090/api/v1/label/http_status_code/values`
```json
{
   "status" : "success",
   "data" : [
      "200",
      "504"
   ]
}
```

نام‌های برچسب را می‌توان به صورت اختیاری با استفاده از روش فرار از مقادیر (Values Escaping) کدگذاری کرد، و در صورتی که نام شامل کاراکتر `/` باشد، این کار ضروری است. برای کدگذاری یک نام به این روش:

*   `U__` را به ابتدای برچسب اضافه کنید.
*   حروف، اعداد و دونقطه (:) به همان صورت ظاهر می‌شوند.
*   تک زیرخط‌ها را به دو زیرخط تبدیل کنید.
*   برای تمام کاراکترهای دیگر، از کدپوینت UTF-8 به عنوان یک عدد صحیح هگزا دسیمال استفاده کنید که با زیرخط احاطه شده باشد. بنابراین ` ` به `_20_` و `.` به `_2e_` تبدیل می‌شود.

اطلاعات بیشتر در مورد فرار از متن را می‌توانید در سند پیشنهاد اصلی UTF-8 بیابید.

این مثال تمام مقادیر برچسب را برای برچسب `http.status_code` کوئری‌ می‌کند:

`curl http://localhost:9090/api/v1/label/U__http_2e_status_code/values`
```json
{
   "status" : "success",
   "data" : [
      "200",
      "404"
   ]
}
```

### کوئری نمونه‌ها (Querying exemplars)

این بخش آزمایشی است و ممکن است در آینده تغییر کند. نقطه پایانی زیر لیستی از نمونه‌ها را برای یک کوئری PromQL معتبر در یک بازه زمانی خاص بازمی‌گرداند:

`GET /api/v1/query_exemplars`
`POST /api/v1/query_exemplars`

پارامترهای کوئری URL:

*  ‏ query=string: رشته کوئری عبارت پرومتئوس.
* ‏  start=<rfc3339 | unix_timestamp>: مهر زمانی شروع.
*  ‏ end=<rfc3339 | unix_timestamp>: مهر زمانی پایان.

`curl -g 'http://localhost:9090/api/v1/query_exemplars?query=test_exemplar_metric_total&start=2020-09-14T15:22:25.479Z&end=2020-09-14T15:23:25.479Z'`
```json
{
    "status": "success",
    "data": [
        {
            "seriesLabels": {
                "__name__": "test_exemplar_metric_total",
                "instance": "localhost:8090",
                "job": "prometheus",
                "service": "bar"
            },
            "exemplars": [
                {
                    "labels": {
                        "trace_id": "EpTxMJ40fUus7aGY"
                    },
                    "value": "6",
                    "timestamp": 1600096945.479
                }
            ]
        },
        {
            "seriesLabels": {
                "__name__": "test_exemplar_metric_total",
                "instance": "localhost:8090",
                "job": "prometheus",
                "service": "foo"
            },
            "exemplars": [
                {
                    "labels": {
                        "trace_id": "Olp9XHlq763ccsfa"
                    },
                    "value": "19",
                    "timestamp": 1600096955.479
                },
                {
                    "labels": {
                        "trace_id": "hCtjygkIHwAN9vs4"
                    },
                    "value": "20",
                    "timestamp": 1600096965.489
                }
            ]
        }
    ]
}
```

### قالب‌های نتیجه کوئری عبارت (Expression query result formats)

کوئری‌های عبارتی ممکن است مقادیر پاسخ زیر را در ویژگی `result` بخش `data` بازگردانند. نگه‌دارنده‌های `<sample_value>` مقادیر نمونه عددی هستند. JSON از مقادیر اعشاری خاص مانند NaN، Inf و -Inf پشتیبانی نمی‌کند، بنابراین مقادیر نمونه به عنوان رشته‌های JSON نقل‌قول شده منتقل می‌شوند تا اعداد خام.

کلیدهای "histogram" و "histograms" فقط در صورتی ظاهر می‌شوند که هیستوگرام‌های بومی آزمایشی در پاسخ وجود داشته باشند. نگه‌دارنده `<histogram>` آن‌ها در بخش مربوط به خودش با جزئیات توضیح داده شده است.

### بردارهای بازه‌ای (Range vectors)

بردارهای بازه‌ای به عنوان نوع نتیجه `matrix` بازگردانده می‌شوند. ویژگی `result` مربوطه دارای قالب زیر است:

```json
[
  {
    "metric": { "<label_name>": "<label_value>", ... },
    "values": [ [ <unix_time>, "<sample_value>" ], ... ],
    "histograms": [ [ <unix_time>, <histogram> ], ... ]
  },
  ...
]
```

هر سری می‌تواند کلید "values"، یا کلید "histograms"، یا هر دو را داشته باشد. برای یک مهر زمانی مشخص، فقط یک نمونه از نوع اعشاری یا هیستوگرام وجود خواهد داشت.

سری‌ها بر اساس معیار مرتب شده بازگردانده می‌شوند. توابع مانند `sort` و `sort_by_label` برای بردارهای بازه‌ای تأثیری ندارند.

### بردارهای لحظه‌ای (Instant vectors)

بردارهای لحظه‌ای به عنوان نوع نتیجه `vector` بازگردانده می‌شوند. ویژگی `result` مربوطه دارای قالب زیر است:

```json
[
  {
    "metric": { "<label_name>": "<label_value>", ... },
    "value": [ <unix_time>, "<sample_value>" ],
    "histogram": [ <unix_time>, <histogram> ]
  },
  ...
]
```

هر سری می‌تواند کلید "value" یا کلید "histogram" را داشته باشد، اما نه هر دو.

سری‌ها تضمینی برای بازگردانده شدن به ترتیب خاصی ندارند، مگر اینکه از تابعی مانند `sort` یا `sort_by_label` استفاده شود.

### مقادیر اسکالر (Scalars)

نتایج اسکالر به عنوان نوع نتیجه `scalar` بازگردانده می‌شوند. ویژگی `result` مربوطه دارای قالب زیر است:

```json
[ <unix_time>, "<scalar_value>" ]
```

رشته‌ها (Strings)

نتایج رشته‌ای به عنوان نوع نتیجه `string` بازگردانده می‌شوند. ویژگی `result` مربوطه دارای قالب زیر است:

```json
[ <unix_time>, "<string_value>" ]
```

### هیستوگرام‌های بومی (Native histograms)

نگه‌دارنده `<histogram>` که در بالا استفاده شد، به صورت زیر قالب‌بندی می‌شود.

توجه داشته باشید که هیستوگرام‌های بومی یک ویژگی آزمایشی هستند و قالب زیر ممکن است همچنان تغییر کند.

```json
{
  "count": "<count_of_observations>",
  "sum": "<sum_of_observations>",
  "buckets": [ [ <boundary_rule>, "<left_boundary>", "<right_boundary>", "<count_in_bucket>" ], ... ]
}
```

نگه‌دارنده `<boundary_rule>` یک عدد صحیح بین 0 تا 3 با معانی زیر است:

*   0: "باز از چپ" (حد چپ انحصاری، حد راست شامل)
*   1: "باز از راست" (حد چپ شامل، حد راست انحصاری)
*   2: "باز از هر دو" (هر دو حد انحصاری)
*   3: "بسته از هر دو" (هر دو حد شامل)

توجه داشته باشید که با طرحواره‌های سطل (bucket schemas) فعلی پیاده‌سازی شده، سطل‌های مثبت "باز از چپ"، سطل‌های منفی "باز از راست"، و سطل صفر (با حد چپ منفی و حد راست مثبت) "بسته از هر دو" هستند.

### هدف‌ها (Targets)

نقطه پایانی زیر نمای کلی از وضعیت فعلی کشف هدف‌های پرومتئوس را بازمی‌گرداند:

`GET /api/v1/targets`

هم هدف‌های فعال و هم هدف‌های حذف شده به طور پیش‌فرض بخشی از پاسخ هستند. هدف‌های حذف شده مشمول محدودیت `keep_dropped_targets`، در صورت تنظیم، هستند. `labels` مجموعه برچسب‌ها را پس از انجام عملیات بازبرچسب‌گذاری (relabeling) نشان می‌دهد. `discoveredLabels` برچسب‌های بدون تغییر را که در طول کشف سرویس قبل از بازبرچسب‌گذاری بازیابی شده‌اند، نشان می‌دهد.

`curl http://localhost:9090/api/v1/targets`
```json
{
  "status": "success",
  "data": {
    "activeTargets": [
      {
        "discoveredLabels": {
          "__address__": "127.0.0.1:9090",
          "__metrics_path__": "/metrics",
          "__scheme__": "http",
          "job": "prometheus"
        },
        "labels": {
          "instance": "127.0.0.1:9090",
          "job": "prometheus"
        },
        "scrapePool": "prometheus",
        "scrapeUrl": "http://127.0.0.1:9090/metrics",
        "globalUrl": "http://example-prometheus:9090/metrics",
        "lastError": "",
        "lastScrape": "2017-01-17T15:07:44.723715405+01:00",
        "lastScrapeDuration": 0.050688943,
        "health": "up",
        "scrapeInterval": "1m",
        "scrapeTimeout": "10s"
      }
    ],
    "droppedTargets": [
      {
        "discoveredLabels": {
          "__address__": "127.0.0.1:9100",
          "__metrics_path__": "/metrics",
          "__scheme__": "http",
          "__scrape_interval__": "1m",
          "__scrape_timeout__": "10s",
          "job": "node"
        },
        "scrapePool": "node"
      }
    ]
  }
}
```

پارامتر کوئری `state` به تماس‌گیرنده اجازه می‌دهد تا بر اساس هدف‌های فعال یا حذف شده فیلتر کند (به عنوان مثال، `state=active`، `state=dropped`، `state=any`). توجه داشته باشید که یک آرایه خالی همچنان برای هدف‌هایی که فیلتر شده‌اند، بازگردانده می‌شود. سایر مقادیر نادیده گرفته می‌شوند.

`curl 'http://localhost:9090/api/v1/targets?state=active'`
```json
  "status": "success",
  "data": {
    "activeTargets": [
      {
        "discoveredLabels": {
          "__address__": "127.0.0.1:9090",
          "__metrics_path__": "/metrics",
          "__scheme__": "http",
          "job": "prometheus"
        },
        "labels": {
          "instance": "127.0.0.1:9090",
          "job": "prometheus"
        },
        "scrapePool": "prometheus",
        "scrapeUrl": "http://127.0.0.1:9090/metrics",
        "globalUrl": "http://example-prometheus:9090/metrics",
        "lastError": "",
        "lastScrape": "2017-01-17T15:07:44.723715405+01:00",
        "lastScrapeDuration": 50688943,
        "health": "up"
      }
    ],
    "droppedTargets": []
  }
}
```

پارامتر کوئری `scrapePool` به تماس‌گیرنده اجازه می‌دهد تا بر اساس نام مجموعه جمع‌آوری (scrape pool) فیلتر کند.

`curl 'http://localhost:9090/api/v1/targets?scrapePool=node_exporter'`
```json
  "status": "success",
  "data": {
    "activeTargets": [
      {
        "discoveredLabels": {
          "__address__": "127.0.0.1:9091",
          "__metrics_path__": "/metrics",
          "__scheme__": "http",
          "job": "node_exporter"
        },
        "labels": {
          "instance": "127.0.0.1:9091",
          "job": "node_exporter"
        },
        "scrapePool": "node_exporter",
        "scrapeUrl": "http://127.0.0.1:9091/metrics",
        "globalUrl": "http://example-prometheus:9091/metrics",
        "lastError": "",
        "lastScrape": "2017-01-17T15:07:44.723715405+01:00",
        "lastScrapeDuration": 50688943,
        "health": "up"
      }
    ],
    "droppedTargets": []
  }
}
```

قوانین (Rules)

نقطه پایانی API `/rules` لیستی از قوانین هشداردهنده (alerting rules) و ثبت‌کننده (recording rules) را که در حال حاضر بارگذاری شده‌اند، بازمی‌گرداند. علاوه بر این، هشدارهای فعال فعلی را که توسط نمونه پرومتئوس هر قانون هشداردهنده فعال شده‌اند، برمی‌گرداند.

از آنجا که نقطه پایانی `/rules` نسبتاً جدید است، تضمین‌های پایداری مشابه API v1 اصلی را ندارد.

`GET /api/v1/rules`

پارامترهای کوئری URL:

*   ‏type=alert|record: فقط قوانین هشداردهنده (مانند `type=alert`) یا قوانین ثبت‌کننده (مانند `type=record`) را بازگردانید. وقتی پارامتر وجود ندارد یا خالی است، هیچ فیلتری انجام نمی‌شود.
*   ‏rule_name[]=string: فقط قوانینی را با نام قانون داده شده بازگردانید. اگر پارامتر تکرار شود، قوانینی با هر یک از نام‌های ارائه شده بازگردانده می‌شوند. اگر تمام قوانین یک گروه را فیلتر کرده باشیم، گروه بازگردانده نمی‌شود. وقتی پارامتر وجود ندارد یا خالی است، هیچ فیلتری انجام نمی‌شود.
*  ‏ rule_group[]=string: فقط قوانینی را با نام گروه قانون داده شده بازگردانید. اگر پارامتر تکرار شود، قوانینی با هر یک از نام‌های گروه قانون ارائه شده بازگردانده می‌شوند. وقتی پارامتر وجود ندارد یا خالی است، هیچ فیلتری انجام نمی‌شود.
*   ‏file[]=string: فقط قوانینی را با مسیر فایل داده شده بازگردانید. اگر پارامتر تکرار شود، قوانینی با هر یک از مسیرهای فایل ارائه شده بازگردانده می‌شوند. وقتی پارامتر وجود ندارد یا خالی است، هیچ فیلتری انجام نمی‌شود.
*  ‏ exclude_alerts=bool: فقط قوانین را بازگردانید، هشدارهای فعال را بازنگردانید.
*   ‏match[]=<label_selector>: فقط قوانینی را بازگردانید که دارای برچسب‌های پیکربندی شده‌ای هستند که انتخاب‌گرهای برچسب را ارضا می‌کنند. اگر پارامتر تکرار شود، قوانینی که با هر یک از مجموعه‌های انتخاب‌گر برچسب مطابقت دارند، بازگردانده می‌شوند. توجه داشته باشید که مطابقت بر روی برچسب‌های موجود در تعریف هر قانون است، نه بر روی مقادیر پس از توسعه قالب (برای قوانین هشداردهنده). اختیاری.
*  ‏ group_limit=number: پارامتر `group_limit` به شما اجازه می‌دهد تا محدودیتی برای تعداد گروه‌های قانون که در یک پاسخ واحد بازگردانده می‌شوند، تعیین کنید. اگر تعداد کل گروه‌های قانون از مقدار `group_limit` مشخص شده بیشتر شود، پاسخ شامل یک ویژگی `groupNextToken` خواهد بود. می‌توانید از مقدار این ویژگی `groupNextToken` در درخواست‌های بعدی در پارامتر `group_next_token` برای صفحه‌بندی بر روی گروه‌های قانون باقی‌مانده استفاده کنید. ویژگی `groupNextToken` در پاسخ نهایی وجود نخواهد داشت که نشان می‌دهد تمام گروه‌های قانون موجود را بازیابی کرده‌اید. لطفاً توجه داشته باشید که هیچ تضمینی در مورد سازگاری پاسخ وجود ندارد اگر گروه‌های قانون در طول فرآیند صفحه‌بندی در حال تغییر باشند.
*   ‏group_next_token: توکن صفحه‌بندی که در درخواست قبلی زمانی که ویژگی `group_limit` تنظیم شده بود، بازگردانده شد. توکن صفحه‌بندی برای صفحه‌بندی تکراری بر روی تعداد زیادی از گروه‌های قانون استفاده می‌شود. برای استفاده از پارامتر `group_next_token`، پارامتر `group_limit` نیز باید وجود داشته باشد. اگر یک گروه قانون که با توکن بعدی همزمان است، در حالی که شما در حال صفحه‌بندی بر روی گروه‌های قانون هستید، حذف شود، یک پاسخ با کد وضعیت 400 بازگردانده خواهد شد.

`curl http://localhost:9090/api/v1/rules`
```json
{
    "data": {
        "groups": [
            {
                "rules": [
                    {
                        "alerts": [
                            {
                                "activeAt": "2018-07-04T20:27:12.60602144+02:00",
                                "annotations": {
                                    "summary": "High request latency"
                                },
                                "labels": {
                                    "alertname": "HighRequestLatency",
                                    "severity": "page"
                                },
                                "state": "firing",
                                "value": "1e+00"
                            }
                        ],
                        "annotations": {
                            "summary": "High request latency"
                        },
                        "duration": 600,
                        "health": "ok",
                        "labels": {
                            "severity": "page"
                        },
                        "name": "HighRequestLatency",
                        "query": "job:request_latency_seconds:mean5m{job=\"myjob\"} > 0.5",
                        "type": "alerting"
                    },
                    {
                        "health": "ok",
                        "name": "job:http_inprogress_requests:sum",
                        "query": "sum by (job) (http_inprogress_requests)",
                        "type": "recording"
                    }
                ],
                "file": "/rules.yaml",
                "interval": 60,
                "limit": 0,
                "name": "example"
            }
        ]
    },
    "status": "success"
}
```

### هشدارها (Alerts)

نقطه پایانی `/alerts` لیستی از تمام هشدارهای فعال را بازمی‌گرداند.

از آنجا که نقطه پایانی `/alerts` نسبتاً جدید است، تضمین‌های پایداری مشابه API v1 اصلی را ندارد.

`GET /api/v1/alerts`

`curl http://localhost:9090/api/v1/alerts`
```json
{
    "data": {
        "alerts": [
            {
                "activeAt": "2018-07-04T20:27:12.60602144+02:00",
                "annotations": {},
                "labels": {
                    "alertname": "my-alert"
                },
                "state": "firing",
                "value": "1e+00"
            }
        ]
    },
    "status": "success"
}
```

### کوئری فراداده هدف (Querying target metadata)

نقطه پایانی زیر فراداده مربوط به معیارهایی را که در حال حاضر از هدف‌ها جمع‌آوری (scraped) می‌شوند، بازمی‌گرداند. این بخش آزمایشی است و ممکن است در آینده تغییر کند.

`GET /api/v1/targets/metadata`

پارامترهای کوئری URL:

*  ‏ match_target=<label_selectors>: انتخاب‌گرهای برچسب که هدف‌ها را بر اساس مجموعه‌های برچسب‌شان مطابقت می‌دهند. اگر خالی گذاشته شود، تمام هدف‌ها انتخاب می‌شوند.
*  ‏ metric=string: نام معیاری که فراداده آن بازیابی شود. اگر خالی گذاشته شود، تمام فراداده معیارها بازیابی می‌شوند.
*  ‏ limit=number: حداکثر تعداد هدف برای مطابقت.

بخش `data` نتیجه کوئری‌ شامل لیستی از اشیا است که حاوی فراداده معیار و مجموعه برچسب هدف هستند.

مثال زیر تمام ورودی‌های فراداده را برای معیار `go_goroutines` از دو هدف اول با برچسب `job="prometheus"` بازمی‌گرداند.

`curl -G http://localhost:9091/api/v1/targets/metadata \`
`    --data-urlencode 'metric=go_goroutines' \`
`    --data-urlencode 'match_target={job="prometheus"}' \`
`    --data-urlencode 'limit=2'`
```json
{
  "status": "success",
  "data": [
    {
      "target": {
        "instance": "127.0.0.1:9090",
        "job": "prometheus"
      },
      "type": "gauge",
      "help": "Number of goroutines that currently exist.",
      "unit": ""
    },
    {
      "target": {
        "instance": "127.0.0.1:9091",
        "job": "prometheus"
      },
      "type": "gauge",
      "help": "Number of goroutines that currently exist.",
      "unit": ""
    }
  ]
}
```

مثال زیر فراداده را برای تمام معیارها برای تمام هدف‌ها با برچسب `instance="127.0.0.1:9090"` بازمی‌گرداند.

`curl -G http://localhost:9091/api/v1/targets/metadata \`
`    --data-urlencode 'match_target={instance="127.0.0.1:9090"}'`
```json
{
  "status": "success",
  "data": [
    // ...
    {
      "target": {
        "instance": "127.0.0.1:9090",
        "job": "prometheus"
      },
      "metric": "prometheus_treecache_zookeeper_failures_total",
      "type": "counter",
      "help": "The total number of ZooKeeper failures.",
      "unit": ""
    },
    {
      "target": {
        "instance": "127.0.0.1:9090",
        "job": "prometheus"
      },
      "metric": "prometheus_tsdb_reloads_total",
      "type": "counter",
      "help": "Number of times the database reloaded block data from disk.",
      "unit": ""
    },
    // ...
  ]
}
```

کوئری فراداده معیار (Querying metric metadata)

این نقطه پایانی فراداده مربوط به معیارهایی را که در حال حاضر از هدف‌ها جمع‌آوری می‌شوند، بازمی‌گرداند. با این حال، هیچ اطلاعات هدفی را ارائه نمی‌دهد. این بخش آزمایشی در نظر گرفته می‌شود و ممکن است در آینده تغییر کند.

`GET /api/v1/metadata`

پارامترهای کوئری URL:

*   ‏limit=number: حداکثر تعداد معیار برای بازگرداندن.
*   ‏limit_per_metric=number: حداکثر تعداد فراداده برای بازگرداندن به ازای هر معیار.
*   ‏metric=string: نام معیاری برای فیلتر کردن فراداده. اگر خالی گذاشته شود، تمام فراداده معیارها بازیابی می‌شوند.

بخش `data` نتیجه کوئری‌ شامل یک شیء است که هر کلید آن یک نام معیار و هر مقدار آن لیستی از اشیاء فراداده منحصر به فرد است، همانطور که برای آن نام معیار در تمام هدف‌ها نمایش داده می‌شود.

مثال زیر دو معیار را بازمی‌گرداند. توجه داشته باشید که معیار `http_requests_total` بیش از یک شیء در لیست دارد. حداقل یک هدف دارای مقداری برای `HELP` است که با بقیه مطابقت ندارد.

`curl -G http://localhost:9090/api/v1/metadata?limit=2`
```json
{
  "status": "success",
  "data": {
    "cortex_ring_tokens": [
      {
        "type": "gauge",
        "help": "Number of tokens in the ring",
        "unit": ""
      }
    ],
    "http_requests_total": [
      {
        "type": "counter",
        "help": "Number of HTTP requests",
        "unit": ""
      },
      {
        "type": "counter",
        "help": "Amount of HTTP requests",
        "unit": ""
      }
    ]
  }
}
```

مثال زیر فقط یک ورودی فراداده را برای هر معیار بازمی‌گرداند.

`curl -G http://localhost:9090/api/v1/metadata?limit_per_metric=1`
```json
{
  "status": "success",
  "data": {
    "cortex_ring_tokens": [
      {
        "type": "gauge",
        "help": "Number of tokens in the ring",
        "unit": ""
      }
    ],
    "http_requests_total": [
      {
        "type": "counter",
        "help": "Number of HTTP requests",
        "unit": ""
      }
    ]
  }
}
```

مثال زیر فقط فراداده را برای معیار `http_requests_total` بازمی‌گرداند.

`curl -G http://localhost:9090/api/v1/metadata?metric=http_requests_total`
```json
{
  "status": "success",
  "data": {
    "http_requests_total": [
      {
        "type": "counter",
        "help": "Number of HTTP requests",
        "unit": ""
      },
      {
        "type": "counter",
        "help": "Amount of HTTP requests",
        "unit": ""
      }
    ]
  }
}
```

مدیران هشدار (Alertmanagers)

نقطه پایانی زیر نمای کلی از وضعیت فعلی کشف مدیران هشدار پرومتئوس را بازمی‌گرداند:

`GET /api/v1/alertmanagers`

هم مدیران هشدار فعال و هم مدیران هشدار حذف شده بخشی از پاسخ هستند.

`curl http://localhost:9090/api/v1/alertmanagers`
```json
{
  "status": "success",
  "data": {
    "activeAlertmanagers": [
      {
        "url": "http://127.0.0.1:9090/api/v1/alerts"
      }
    ],
    "droppedAlertmanagers": [
      {
        "url": "http://127.0.0.1:9093/api/v1/alerts"
      }
    ]
  }
}
```

وضعیت (Status)

نقاط پایانی وضعیت زیر، پیکربندی فعلی پرومتئوس را نمایش می‌دهند.

پیکربندی (Config)

نقطه پایانی زیر فایل پیکربندی بارگذاری شده فعلی را بازمی‌گرداند:

`GET /api/v1/status/config`

پیکربندی به عنوان فایل YAML خروجی شده بازگردانده می‌شود. به دلیل محدودیت کتابخانه YAML، نظرات YAML شامل نمی‌شوند.

`curl http://localhost:9090/api/v1/status/config`
```json
{
  "status": "success",
  "data": {
    "yaml": "<محتوای فایل پیکربندی بارگذاری شده در قالب YAML>",
  }
}
```

پرچم‌ها (Flags)

نقطه پایانی زیر مقادیر پرچم‌هایی را که پرومتئوس با آنها پیکربندی شده بود، بازمی‌گرداند:

`GET /api/v1/status/flags`

تمام مقادیر از نوع نتیجه رشته هستند.

`curl http://localhost:9090/api/v1/status/flags`
```json
{
  "status": "success",
  "data": {
    "alertmanager.notification-queue-capacity": "10000",
    "alertmanager.timeout": "10s",
    "log.level": "info",
    "query.lookback-delta": "5m",
    "query.max-concurrency": "20",
    ...
  }
}
```

جدید در نسخه 2.2

اطلاعات زمان اجرا (Runtime Information)

نقطه پایانی زیر ویژگی‌های مختلف اطلاعات زمان اجرا را در مورد سرور پرومتئوس بازمی‌گرداند:

`GET /api/v1/status/runtimeinfo`

مقادیر بازگردانده شده از انواع مختلفی هستند که به ماهیت ویژگی زمان اجرا بستگی دارد.

`curl http://localhost:9090/api/v1/status/runtimeinfo`
```json
{
  "status": "success",
  "data": {
    "startTime": "2019-11-02T17:23:59.301361365+01:00",
    "CWD": "/",
    "hostname" : "DESKTOP-717H17Q",
    "serverTime": "2025-01-05T18:27:33Z",
    "reloadConfigSuccess": true,
    "lastConfigTime": "2019-11-02T17:23:59+01:00",
    "timeSeriesCount": 873,
    "corruptionCount": 0,
    "goroutineCount": 48,
    "GOMAXPROCS": 4,
    "GOGC": "",
    "GODEBUG": "",
    "storageRetention": "15d"
  }
}
```

نکته: ویژگی‌های زمان اجرای دقیق بازگردانده شده ممکن است بدون اطلاع قبلی بین نسخه‌های پرومتئوس تغییر کند.
جدید در نسخه 2.14

اطلاعات ساخت (Build Information)

نقطه پایانی زیر ویژگی‌های مختلف اطلاعات ساخت را در مورد سرور پرومتئوس بازمی‌گرداند:

`GET /api/v1/status/buildinfo`

تمام مقادیر از نوع نتیجه رشته هستند.

`curl http://localhost:9090/api/v1/status/buildinfo`
```json
{
  "status": "success",
  "data": {
    "version": "2.13.1",
    "revision": "cb7cbad5f9a2823a622aaa668833ca04f50a0ea7",
    "branch": "master",
    "buildUser": "julius@desktop",
    "buildDate": "20191102-16:19:59",
    "goVersion": "go1.13.1"
  }
}
```

نکته: ویژگی‌های ساخت دقیق بازگردانده شده ممکن است بدون اطلاع قبلی بین نسخه‌های پرومتئوس تغییر کند.
جدید در نسخه 2.14

### آمار TSDB (TSDB Stats)

نقطه پایانی زیر آمارهای مختلف کاردینالیتی را در مورد TSDB پرومتئوس بازمی‌گرداند:

`GET /api/v1/status/tsdb`

پارامترهای کوئری URL:

*  ‏ limit=number: تعداد موارد بازگردانده شده را برای هر مجموعه آمار به یک عدد مشخص محدود می‌کند. به طور پیش‌فرض، 10 مورد بازگردانده می‌شود.

بخش `data` نتیجه کوئری‌ شامل موارد زیر است:

*   ‏headStats: این بخش داده‌های زیر را در مورد بلوک سر (head block) TSDB ارائه می‌دهد:
    *  ‏ numSeries: تعداد سری‌ها.
    *  ‏ chunkCount: تعداد قطعات (chunks).
    *  ‏ minTime: حداقل مهر زمانی فعلی بر حسب میلی‌ثانیه.
    *   ‏maxTime: حداکثر مهر زمانی فعلی بر حسب میلی‌ثانیه.
*   ‏seriesCountByMetricName: این بخش لیستی از نام‌های معیار و تعداد سری‌های آن‌ها را ارائه می‌دهد.
*  ‏ labelValueCountByLabelName: این بخش لیستی از نام‌های برچسب و تعداد مقادیر آن‌ها را ارائه می‌دهد.
*  ‏ memoryInBytesByLabelName: این بخش لیستی از نام‌های برچسب و حافظه استفاده شده بر حسب بایت را ارائه می‌دهد. مصرف حافظه با جمع کردن طول تمام مقادیر برای یک نام برچسب مشخص محاسبه می‌شود.
*  ‏ seriesCountByLabelPair: این بخش لیستی از جفت‌های مقدار برچسب و تعداد سری‌های آن‌ها را ارائه می‌دهد.

`curl http://localhost:9090/api/v1/status/tsdb`
```json
{
  "status": "success",
  "data": {
    "headStats": {
      "numSeries": 508,
      "chunkCount": 937,
      "minTime": 1591516800000,
      "maxTime": 1598896800143,
    },
    "seriesCountByMetricName": [
      {
        "name": "net_conntrack_dialer_conn_failed_total",
        "value": 20
      },
      {
        "name": "prometheus_http_request_duration_seconds_bucket",
        "value": 20
      }
    ],
    "labelValueCountByLabelName": [
      {
        "name": "__name__",
        "value": 211
      },
      {
        "name": "event",
        "value": 3
      }
    ],
    "memoryInBytesByLabelName": [
      {
        "name": "__name__",
        "value": 8266
      },
      {
        "name": "instance",
        "value": 28
      }
    ],
    "seriesCountByLabelValuePair": [
      {
        "name": "job=prometheus",
        "value": 425
      },
      {
        "name": "instance=localhost:9090",
        "value": 425
      }
    ]
  }
}
```

جدید در نسخه 3.6.0

### بلوک‌های TSDB (TSDB Blocks)

نکته: این نقطه پایانی آزمایشی است و ممکن است در آینده تغییر کند. نام نقطه پایانی و قالب دقیق داده‌های بازگردانده شده ممکن است بین نسخه‌های پرومتئوس تغییر کند. فراداده دقیق بازگردانده شده توسط این نقطه پایانی یک جزئیات پیاده‌سازی است و ممکن است در نسخه‌های آینده پرومتئوس تغییر کند.

نقطه پایانی زیر لیستی از بلوک‌های TSDB بارگذاری شده فعلی و فراداده آن‌ها را بازمی‌گرداند.

`GET /api/v1/status/tsdb/blocks`

این نقطه پایانی اطلاعات زیر را برای هر بلوک بازمی‌گرداند:

*  ‏ ulid: شناسه منحصر به فرد بلوک.
*   ‏minTime: حداقل مهر زمانی (بر حسب میلی‌ثانیه) بلوک.
*  ‏ maxTime: حداکثر مهر زمانی (بر حسب میلی‌ثانیه) بلوک.
*  ‏ stats:
    ‏*   numSeries: تعداد سری‌ها در بلوک.
    ‏*   numSamples: تعداد نمونه‌ها در بلوک.
    ‏*   numChunks: تعداد قطعات (chunks) در بلوک.
*  ‏ compaction:
    ‏*   level: سطح فشرده‌سازی (compaction level) بلوک.
    ‏*   sources: لیستی از ULIDهای بلوک‌های منبع که برای فشرده‌سازی این بلوک استفاده شده‌اند.
* ‏  version: نسخه بلوک.

`curl http://localhost:9090/api/v1/status/tsdb/blocks`
```json
{
  "status": "success",
  "data": {
    "blocks": [
      {
        "ulid": "01JZ8JKZY6XSK3PTDP9ZKRWT60",
        "minTime": 1750860620060,
        "maxTime": 1750867200000,
        "stats": {
          "numSamples": 13701,
          "numSeries": 716,
          "numChunks": 716
        },
        "compaction": {
          "level": 1,
          "sources": [
            "01JZ8JKZY6XSK3PTDP9ZKRWT60"
          ]
        },
        "version": 1
      }
    ]
  }
}
```

جدید در نسخه 2.15

### آمار پخش مجدد WAL (WAL Replay Stats)

نقطه پایانی زیر اطلاعاتی در مورد پخش مجدد WAL را بازمی‌گرداند:

`GET /api/v1/status/walreplay`

*‏   read: تعداد سگمنت‌های پخش شده تا کنون.
* ‏  total: تعداد کل سگمنت‌هایی که نیاز به پخش مجدد دارند.
*  ‏ progress: پیشرفت پخش مجدد (0 - 100%).
*  ‏ state: وضعیت پخش مجدد. وضعیت‌های ممکن:
    ‏*   waiting: در انتظار شروع پخش مجدد.
    ‏*   in progress: پخش مجدد در حال انجام است.
    ‏*   done: پخش مجدد به پایان رسیده است.

`curl http://localhost:9090/api/v1/status/walreplay`
```json
{
  "status": "success",
  "data": {
    "min": 2,
    "max": 5,
    "current": 40,
    "state": "in progress"
  }
}
```

نکته: این نقطه پایانی قبل از اینکه سرور آماده علامت‌گذاری شود، در دسترس است و در زمان واقعی به روز می‌شود تا نظارت بر پیشرفت پخش مجدد WAL را تسهیل کند.
جدید در نسخه 2.28

‏### APIهای مدیریتی TSDB (TSDB Admin APIs)

اینها APIهایی هستند که قابلیت‌های پایگاه داده را برای کاربران پیشرفته نمایش می‌دهند. این APIها فعال نمی‌شوند مگر اینکه پرچم `--web.enable-admin-api` تنظیم شده باشد.

عکس فوری (Snapshot)

‏Snapshot یک Snapshot از تمام داده‌های فعلی را در `snapshots/<datetime>-<rand>` در دایرکتوری داده TSDB ایجاد می‌کند و آن دایرکتوری را به عنوان پاسخ بازمی‌گرداند. این عملیات به صورت اختیاری از ایجاد Snapshot از داده‌هایی که فقط در بلوک سر (head block) وجود دارند و هنوز به دیسک فشرده نشده‌اند، صرف نظر می‌کند.

`POST /api/v1/admin/tsdb/snapshot`
`PUT /api/v1/admin/tsdb/snapshot`

پارامترهای کوئری URL:

*   ‏skip_head=bool: از داده‌های موجود در بلوک سر صرف نظر کنید. اختیاری.

`curl -XPOST http://localhost:9090/api/v1/admin/tsdb/snapshot`
```json
{
  "status": "success",
  "data": {
    "name": "20171210T211224Z-2be650b6d019eb54"
  }
}
```

عکس فوری اکنون در `<data-dir>/snapshots/20171210T211224Z-2be650b6d019eb54` وجود دارد.

جدید در نسخه 2.1 و پشتیبانی از PUT از نسخه 2.9

حذف سری‌ها (Delete Series)

‏DeleteSeries داده‌های یک مجموعه انتخابی از سری‌ها را در یک بازه زمانی حذف می‌کند. داده‌های واقعی همچنان روی دیسک وجود دارند و در فشرده‌سازی‌های آینده پاک می‌شوند یا می‌توانند به صراحت با دسترسی به نقطه پایانی Clean Tombstones پاک شوند.

در صورت موفقیت، 204 بازگردانده می‌شود.

`POST /api/v1/admin/tsdb/delete_series`
`PUT /api/v1/admin/tsdb/delete_series`

پارامترهای کوئری URL:

*  ‏ match[]=<series_selector>: آرگومان تطابق‌دهنده برچسب تکراری که سری‌ها را برای حذف انتخاب می‌کند. حداقل یک آرگومان `match[]` باید ارائه شود.
*  ‏ start=<rfc3339 | unix_timestamp>: مهر زمانی شروع. اختیاری و به طور پیش‌فرض، حداقل زمان ممکن است.
* ‏  end=<rfc3339 | unix_timestamp>: مهر زمانی پایان. اختیاری و به طور پیش‌فرض، حداکثر زمان ممکن است.

ذکر نکردن هر دو زمان شروع و پایان، تمام داده‌های سری‌های مطابقت داده شده در پایگاه داده را پاک می‌کند.

مثال:

`curl -X POST \`
`  -g 'http://localhost:9090/api/v1/admin/tsdb/delete_series?match[]=up&match[]=process_start_time_seconds{job="prometheus"}'`

نکته: این نقطه پایانی نمونه‌ها را از سری‌ها به عنوان حذف شده علامت‌گذاری می‌کند، اما لزوماً از بازگرداندن فراداده‌های سری مرتبط در کوئری‌های فراداده برای محدوده زمانی مربوطه (حتی پس از پاک کردن سنگ‌قبرها) جلوگیری نمی‌کند. میزان دقیق حذف فراداده یک جزئیات پیاده‌سازی است که ممکن است در آینده تغییر کند.
جدید در نسخه 2.1 و پشتیبانی از PUT از نسخه 2.9

پاک کردن سنگ‌قبرها (Clean Tombstones)

‏CleanTombstones داده‌های حذف شده را از دیسک حذف کرده و سنگ‌قبرهای موجود را پاک می‌کند. این عملیات می‌تواند پس از حذف سری‌ها برای آزاد کردن فضا استفاده شود.

در صورت موفقیت، 204 بازگردانده می‌شود.

`POST /api/v1/admin/tsdb/clean_tombstones`
`PUT /api/v1/admin/tsdb/clean_tombstones`

این عملیات هیچ پارامتر یا بدنه ای نمی‌گیرد.

`curl -XPOST http://localhost:9090/api/v1/admin/tsdb/clean_tombstones`
جدید در نسخه 2.1 و پشتیبانی از PUT از نسخه 2.9

### گیرنده ریموت رایت (Remote Write Receiver)

پرومتئوس را می‌توان به عنوان گیرنده‌ای برای پروتکل ریموت رایت پرومتئوس پیکربندی کرد. این روش برای ورود نمونه‌ها کارآمد در نظر گرفته نمی‌شود. با احتیاط برای موارد استفاده خاص با حجم کم از آن استفاده کنید. برای جایگزینی ورود از طریق جمع‌آوری (scraping) و تبدیل پرومتئوس به یک سیستم جمع‌آوری معیارهای مبتنی بر ارسال (push-based)، مناسب نیست.

گیرنده ریموت رایت را با تنظیم `--web.enable-remote-write-receiver` فعال کنید. هنگامی که فعال شود، نقطه پایانی گیرنده ریموت رایت `/api/v1/write` است. جزئیات بیشتر را اینجا بیابید.

جدید در نسخه 2.33

### گیرنده OTLP (OTLP Receiver)

پرومتئوس را می‌توان به عنوان گیرنده‌ای برای پروتکل معیارهای OTLP پیکربندی کرد. این روش برای ورود نمونه‌ها کارآمد در نظر گرفته نمی‌شود. با احتیاط برای موارد استفاده خاص با حجم کم از آن استفاده کنید. برای جایگزینی ورود از طریق جمع‌آوری (scraping)، مناسب نیست.

گیرنده OTLP را با تنظیم `--web.enable-otlp-receiver` فعال کنید. هنگامی که فعال شود، نقطه پایانی گیرنده OTLP `/api/v1/otlp/v1/metrics` است.

جدید در نسخه 2.47

‏### OTLP دلتا (OTLP Delta)

پرومتئوس می‌تواند معیارهای ورودی را از موقتیت دلتا به معادل تجمعی آن‌ها تبدیل کند. این کار با استفاده از `deltatocumulative` از OpenTelemetry Collector انجام می‌شود.

برای فعال کردن، `--enable-feature=otlp-deltatocumulative` را ارسال کنید.

جدید در نسخه 3.2

اعلانات (Notifications)

نقاط پایانی زیر اطلاعاتی در مورد اعلانات وضعیت فعال مربوط به خود سرور پرومتئوس ارائه می‌دهند. اعلانات در رابط کاربری وب استفاده می‌شوند.

این نقاط پایانی آزمایشی هستند. ممکن است در آینده تغییر کنند.

اعلانات فعال (Active Notifications)

نقطه پایانی `/api/v1/notifications` لیستی از تمام اعلانات فعال فعلی را بازمی‌گرداند.

`GET /api/v1/notifications`

مثال:

`curl http://localhost:9090/api/v1/notifications`
```json
{
  "status": "success",
  "data": [
    {
      "text": "Prometheus is shutting down and gracefully stopping all operations.",
      "date": "2024-10-07T12:33:08.551376578+02:00",
      "active": true
    }
  ]
}
```

جدید در نسخه 3.0

اعلانات زنده (Live Notifications)

نقطه پایانی `/api/v1/notifications/live` اعلانات زنده را به محض وقوع، با استفاده از رویدادهای ارسالی از سرور (Server-Sent Events) پخش می‌کند. اعلانات حذف شده با `active: false` ارسال می‌شوند. اعلانات فعال هنگام اتصال به نقطه پایانی ارسال خواهند شد.

`GET /api/v1/notifications/live`

مثال:

`curl http://localhost:9090/api/v1/notifications/live`
```json
data: {
  "status": "success",
  "data": [
    {
      "text": "Prometheus is shutting down and gracefully stopping all operations.",
      "date": "2024-10-07T12:33:08.551376578+02:00",
      "active": true
    }
  ]
}
```

نکته: نقطه پایانی `/notifications/live` در صورتی که حداکثر تعداد مشترکین (subscribers) به دست آمده باشد، یک پاسخ `204 No Content` را بازمی‌گرداند. می‌توانید حداکثر تعداد شنوندگان را با پرچم `--web.max-notifications-subscribers` تنظیم کنید که به طور پیش‌فرض 16 است.

`GET /api/v1/notifications/live`
`204 No Content`
جدید در نسخه 3.0


------------


این **Lookback Period** یا **Lookback Delta** در پرومتئوس یک بازه‌ٔ زمانی داخلی است که پرومتئوس برای تعیین *آخرین نمونه معتبر یک سری زمانی* از آن استفاده می‌کند.

به‌صورت خیلی ساده:

> **اگر یک سری زمانی در بازه Lookback Period داده‌ای نداشته باشد، پرومتئوس مقدار آن را `stale` در نظر می‌گیرد و دیگر برای کوئری‌ها استفاده نمی‌کند.**

---

## 🔎 مقدار پیش‌فرض Lookback Period

در پرومتئوس مقدار پیش‌فرض:

```
5m
```

یعنی:

* اگر یک متریک **بیش از ۵ دقیقه** نمونه جدید نداشته باشد، از نظر پرومتئوس **منقضی (stale)** است.
* کوئری‌هایی مثل `rate()`، `increase()` و حتی کوئری‌های ساده `[5m]` دیگر از آن استفاده نمی‌کنند.

---

## 📘 مثال ساده:

فرض کن یک متریک داریم:

```
my_counter_total
```

که هر 10 ثانیه یک بار نمونه جدید تولید می‌کند.

اگر به هر دلیلی تولید نمونه‌ها **۶ دقیقه قطع شود**، پرومتئوس:

* مقدار آخرین نمونه را stale تبدیل می‌کند.
* کوئری‌های بعدی این سری زمانی را صفر یا خالی نشان می‌دهند.

---

## 🔍 Lookback period در چه جاهایی استفاده می‌شود؟

### 1. **کوئری‌های بدون رنج (instant queries)**

پرومتئوس هنگام اجرای کوئری لحظه‌ای مثل:

```
my_metric
```

آخرین مثدار را تنها **در صورتی** برمی‌دارد که فاصله زمانی آن با now کمتر از **5m** باشد.

---

### 2. **کوئری‌های رنج مثل rate و increase**

```
rate(my_counter_total[1m])
```

پرومتئوس باید آخرین نمونه موجود در بازه را استفاده کند.
اگر آخرین نمونه *قدیمی‌تر از ۵ دقیقه* باشد → نرخ محاسبه نمی‌شود.

---

## ⚠️ اگر Lookback Period خیلی کم باشد؟

* سری‌هایی که دیر به‌دیر نمونه می‌دهند (مثلاً هر 10 دقیقه یک بار) مرتب stale می‌شوند.
* نمودارها "قطعی" و خالی دیده می‌شوند.

## ⚠️ اگر خیلی زیاد باشد؟

* سری‌هایی که واقعا مرده‌اند، هنوز زنده دیده می‌شوند.
* باعث گمراهی در alertها می‌شود.

---

## 🛠️ تغییر Lookback Period

در پرومتئوس با گزینه زیر در اجرا تنظیم می‌شود:

```
--query.lookback-delta=10m
```

---

## 🧠 خلاصه

| ویژگی                | توضیح                                                     |
| -------------------- | --------------------------------------------------------- |
| نام                  | Lookback Period / Lookback Delta                          |
| مقدار پیش‌فرض        | 5m                                                        |
| هدف                  | تعیین حداکثر فاصله زمانی معتبر بودن آخرین نمونه سری زمانی |
| اگر داده دیرتر بیاید | سری stale می‌شود و در کوئری‌ها نادیده گرفته می‌شود        |

---

اگر خواستی، می‌تونم با چند مثال تصویری یا نمودار مراحل دقیق تشخیص stale شدن در پرومتئوس را هم توضیح بدم.
