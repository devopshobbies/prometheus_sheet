می توان نمونه هایی را به متریک های شمارنده و هیستوگرام اضافه کرد. نمونه ها را می توان با ارسال دیکته ای از جفت های ارزش برچسب مشخص کرد تا به عنوان نمونه در معرض دید قرار گیرند. به عنوان مثال با یک شمارنده:


`./node_exporter --collector.textfile.directory="/var/lib/node_exporter/textfile_collector"`

-------------

این Enum ردیابی می کند که در حال حاضر چیزی در کدام یک از حالت ها قرار دارد.


------------

Exemplars

می توان Exemplars هایی را به متریک های counter و هیستوگرام اضافه کرد. Exemplars ها را می توان با ارسال dict ای از جفت های label value مشخص کرد تا به عنوان Exemplars در معرض دید قرار گیرند , exposed شوند. به عنوان مثال با یک counter:

from prometheus_client import Counter

c = Counter('my_requests_total', 'HTTP Failures', ['method', 'endpoint'])

c.labels('get', '/').inc(exemplar={'trace_id': 'abc123'})

c.labels('post', '/submit').inc(1.0, {'trace_id': 'def456'})

این Exemplars ها فقط در قالب نمایش OpenMetrics ارائه می شوند. در صورت استفاده از سرور یا برنامه های HTTP در این کتابخانه، می توان از مذاکره محتوا برای تعیین OpenMetrics استفاده کرد (که به طور پیش فرض در Prometheus انجام می شود). در غیر این صورت، برای مشاهده Exemplars ها، استفاده از gener_latest از prometheus_client.openmetrics.exposition ضروری است.

برای مشاهده نمونه ها در Prometheus، همچنین لازم است که پرچم ویژگی exemplar-storage را فعال کنید:

------

این Enum ردیابی می کند که در حال حاضر چیزی در کدام یک از حالت ها قرار دارد.

این Info اطلاعات کلید-مقدار را، معمولاً در مورد یک هدف کامل، ردیابی می کند.

-------------


Custom Collectors:

گاهی اوقات نمی توان به طور مستقیم  instrument code را تعیین کرد، زیرا در کنترل شما نیست. این به شما نیاز دارد که proxy metrics را از سایر سیستم ها انجام دهید. برای انجام این کار باید یک custom collector ایجاد کنید، به عنوان مثال:

SummaryMetricFamily, HistogramMetricFamily and InfoMetricFamily work similarly.

یک Collectors ممکن است یک describe method  را پیاده‌سازی کند که متریک‌ها را به همان format collect برمی‌گرداند (البته شما مجبور نیستید samples را اضافه کنید). این برای از پیش تعیین نام سری های زمانی استفاده می شود که CollectorRegistry در معرض نمایش قرار می دهد و بنابراین برای شناسایی collisions و ثبت های تکراری.


معمولاً custom collectors نیازی به پیاده سازی describe ندارند. اگر describe اجرا نشود و CollectorRegistry با auto_describe=True ایجاد شده باشد (که در مورد رجیستری پیش‌فرض صدق می‌کند)، collect در زمان registration به جای توصیف فراخوانی می‌شود. اگر این می تواند باعث ایجاد مشکل شود، یا یک توصیف مناسب را پیاده سازی کنید، یا اگر این عملی نیست، توصیف را یک لیست خالی برگردانید.

-------------------------------------
simple

```python
import http.server
from prometheus_client import start_http_server
class MyHandler(http.server.BaseHTTPRequestHandler):
	def do_GET(self):
		self.send_response(200)
		self.end_headers()
		self.wfile.write(b"Hello World")
		
if __name__ == "__main__":
	start_http_server(8000)
	server = http.server.HTTPServer(('localhost', 8001), MyHandler)
	server.serve_forever()
```


counter

```python
from prometheus_client import Counter
REQUESTS = Counter('hello_worlds_total','Hello Worlds requested.')
class MyHandler(http.server.BaseHTTPRequestHandler):
	def do_GET(self):
		REQUESTS.inc()
		self.send_response(200)
		self.end_headers()
		self.wfile.write(b"Hello World")
```


Counting Exceptions
Client libraries provide not just core functionality, but also utilities and methods for
common use cases. One of these in Python is the ability to count exceptions. You
don’t have to write your own instrumentation using a try…except; instead, you can
take advantage of the count_exceptions context manager and decorator, as shown in
Example 3-4.

```python
import random
from prometheus_client import Counter
REQUESTS = Counter('hello_worlds_total','Hello Worlds requested.')
EXCEPTIONS = Counter('hello_world_exceptions_total','Exceptions serving Hello World.')
class MyHandler(http.server.BaseHTTPRequestHandler):
	def do_GET(self):
	REQUESTS.inc()
		with EXCEPTIONS.count_exceptions():
		if random.random() < 0.2:
			raise Exception
		self.send_response(200)
		self.end_headers()
		self.wfile.write(b"Hello World")
```


-------
Buckets

The default buckets cover a range of latencies from 1 ms to 10 s. This is intended to
capture the typical range of latencies for a web application. But you can also override
them and provide your own buckets when defining metrics.

```python
LATENCY = Histogram('hello_world_latency_seconds','Time for a request Hello World.',buckets=[0.0001, 0.0002, 0.0005, 0.001, 0.01, 0.1])
```


