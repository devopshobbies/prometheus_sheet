
اکسپورتِر نود¹ (Node Exporter) احتمالاً یکی از اولین اکسپورترهایی است که استفاده خواهید کرد، همانطور که قبلاً در فصل ۲ مشاهده شد. این اکسپورتر متریک‌های سطح ماشین را، عمدتاً از هسته سیستم عامل شما، مانند CPU، حافظه، فضای دیسک، ورودی/خروجی دیسک، پهنای باند شبکه و دمای مادربرد، در معرض نمایش قرار می‌دهد. اکسپورتر نود با سیستم‌های یونیکس استفاده می‌شود؛ کاربران ویندوز باید به جای آن از اکسپورتر ویندوز² (Windows Exporter) استفاده کنند. اکسپورتر نود فقط برای نظارت بر خود ماشین در نظر گرفته شده است، نه فرآیندها یا سرویس‌های منفرد روی آن. سایر سیستم‌های مانیتورینگ اغلب چیزی دارند که ما دوست داریم آن را «عامل جامع» (uberagent) بنامیم؛ یعنی یک فرآیند واحد که همه چیز را روی ماشین نظارت می‌کند. در معماری پرومتئوس (Prometheus)، هر یک از سرویس‌های شما متریک‌های خود را، در صورت نیاز با استفاده از یک اکسپورتر، در معرض نمایش قرار می‌دهند که سپس مستقیماً توسط پرومتئوس اسکرِیپ (scrape) می‌شود. این کار از تبدیل شدن عامل جامع به یک گلوگاه عملیاتی یا عملکردی جلوگیری می‌کند و به شما امکان می‌دهد به جای ماشین‌ها، بیشتر در قالب سرویس‌های پویا فکر کنید.

دستورالعمل‌هایی که هنگام ایجاد متریک‌ها با ابزار دقیق مستقیم (direct instrumentation)، مانند موارد مورد بحث در «متریک‌هایم را چه نام‌گذاری کنم؟» در صفحه ۶۰، باید استفاده کنید، نسبتاً واضح هستند. این مورد در مورد اکسپورترها صدق نمی‌کند، جایی که بنا به تعریف، داده‌ها از منبعی می‌آیند که با در نظر گرفتن دستورالعمل‌های پرومتئوس طراحی نشده است. بسته به حجم و کیفیت متریک‌ها، توسعه‌دهندگان اکسپورتر باید بین تلاش مهندسی و به دست آوردن متریک‌های عالی، مصالحه‌هایی انجام دهند. در مورد لینوکس، هزاران متریک در دسترس است. برخی به خوبی مستند و درک شده‌اند، مانند استفاده از CPU؛ برخی دیگر، مانند استفاده از حافظه، با تغییر پیاده‌سازی از نسخه‌ای به نسخه دیگر هسته متفاوت بوده‌اند. حتی متریک‌هایی را خواهید یافت که کاملاً بدون مستندات هستند، جایی که برای فهمیدن کارکرد آن‌ها باید کد منبع هسته را بخوانید. اکسپورتر نود طوری طراحی شده است که به عنوان کاربر غیر روت (non-root) اجرا شود و باید مستقیماً روی ماشین به همان روشی که یک دیمون (daemon) سیستمی مانند sshd یا cron را اجرا می‌کنید، اجرا شود.

در حالی که اجرای اکسپورتر نود در داخل داکر³ (Docker) امکان‌پذیر است، برای مانت (mount) کردن فایل سیستم میزبان در داخل کانتینر، باید از برخی والیوم‌ها (volumes) و پارامترهای خط فرمان (`--path.procfs`, `--path.rootfs`, `--path.sysfs`) استفاده کنید. در صورت امکان، اکسپورتر نود را به عنوان یک سرویس روی نود، بدون داکر اجرا کنید. داکر تلاش می‌کند تا یک کانتینر را از عملکرد داخلی ماشین جدا کند، که این با تلاش اکسپورتر نود برای دسترسی به آن عملکردهای داخلی به خوبی کار نمی‌کند.

برخلاف اکثر اکسپورترهای دیگر، به دلیل تنوع گسترده متریک‌های موجود از سیستم‌عامل‌ها، اکسپورتر نود به شما امکان می‌دهد پیکربندی کنید که کدام دسته‌بندی از متریک‌ها را جمع‌آوری کند. می‌توانید این کار را با فلگ‌های خط فرمان مانند `--collector.wifi` که کالکتور (collector) وای‌فای را فعال می‌کند، و `--no-collector.wifi` که آن را غیرفعال می‌کند، انجام دهید. `--collector.disable-defaults` همه کالکتورها را به جز آن‌هایی که به صراحت با فلگ‌های خط فرمان فعال شده‌اند، غیرفعال می‌کند. پیش‌فرض‌های معقولی تنظیم شده‌اند، بنابراین این چیزی نیست که در شروع کار نگران آن باشید. هسته‌های مختلف متریک‌های متفاوتی را در معرض نمایش قرار می‌دهند، زیرا، به عنوان مثال، لینوکس و FreeBSD کارها را به روش‌های متفاوتی انجام می‌دهند. ممکن است متریک‌ها با گذشت زمان و با بازسازی (refactor) اکسپورتر نود بین کالکتورها جابجا شوند. اگر از سیستم یونیکس دیگری استفاده می‌کنید، متوجه خواهید شد که متریک‌ها و کالکتورهای ارائه شده متفاوت هستند. در این فصل، ما برخی از متریک‌های کلیدی را که اکسپورتر نود نسخه 1.4.0 با هسته لینوکس 5.18.0 در معرض نمایش قرار می‌دهد، توضیح می‌دهیم. این قرار نیست لیست کاملی از متریک‌های موجود باشد. مانند اکثر اکسپورترها و برنامه‌ها، شما می‌خواهید مسیر `/metrics` را بررسی کنید تا ببینید چه چیزی در دسترس است. می‌توانید عبارات PromQL نمونه را با استفاده از تنظیمات خود از فصل ۲ امتحان کنید.

## کالکتور CPU

متریک اصلی از کالکتور CPU، `node_cpu_seconds_total` است، که یک شمارنده (counter) است که نشان می‌دهد هر CPU چقدر زمان در هر حالت (mode) صرف کرده است. لیبل‌ها (labels) `cpu` و `mode` هستند:

```
# HELP node_cpu_seconds_total Seconds the CPUs spent in each mode.
# TYPE node_cpu_seconds_total counter
node_cpu_seconds_total{cpu="0",mode="idle"} 13024.48
node_cpu_seconds_total{cpu="0",mode="iowait"} 9.53
node_cpu_seconds_total{cpu="0",mode="irq"} 0
node_cpu_seconds_total{cpu="0",mode="nice"} 0.11
node_cpu_seconds_total{cpu="0",mode="softirq"} 109.74
node_cpu_seconds_total{cpu="0",mode="steal"} 0
node_cpu_seconds_total{cpu="0",mode="system"} 566.67
node_cpu_seconds_total{cpu="0",mode="user"} 1220.36
node_cpu_seconds_total{cpu="1",mode="idle"} 13501.28
node_cpu_seconds_total{cpu="1",mode="iowait"} 5.96
node_cpu_seconds_total{cpu="1",mode="irq"} 0
node_cpu_seconds_total{cpu="1",mode="nice"} 0.09
node_cpu_seconds_total{cpu="1",mode="softirq"} 23.74
node_cpu_seconds_total{cpu="1",mode="steal"} 0
node_cpu_seconds_total{cpu="1",mode="system"} 423.84
node_cpu_seconds_total{cpu="1",mode="user"} 936.05
```

برای هر CPU، مجموع حالت‌ها در هر ثانیه یک ثانیه افزایش می‌یابد. این به شما امکان می‌دهد نسبت زمان بیکاری (idle) را در تمام CPUها با استفاده از عبارت PromQL محاسبه کنید:
`avg without(cpu, mode)(rate(node_cpu_seconds_total{mode="idle"}[1m]))`

این کار می‌کند زیرا زمان بیکاری در هر ثانیه برای هر CPU را محاسبه می‌کند و سپس آن را در تمام CPUهای ماشین میانگین می‌گیرد. می‌توانید این را برای محاسبه نسبت زمان صرف شده در هر حالت برای یک ماشین با استفاده از عبارت زیر تعمیم دهید:
`avg without(cpu)(rate(node_cpu_seconds_total[1m]))`

استفاده از CPU توسط مهمان‌ها (یعنی ماشین‌های مجازی که تحت هسته اجرا می‌شوند) قبلاً در حالت‌های user و nice گنجانده شده است. می‌توانید زمان مهمان را به طور جداگانه در متریک `node_cpu_guest_seconds_total` مشاهده کنید.

### کالکتور Filesystem (فایل سیستم)

کالکتور فایل سیستم، همانطور که انتظار می‌رود، متریک‌هایی درباره فایل سیستم‌های مانت شده شما جمع‌آوری می‌کند، درست همانطور که از دستور `df` به دست می‌آورید. فلگ‌های `--collector.filesystem.mount-points-exclude` و `--collector.filesystem.fs-types-exclude` اجازه محدود کردن فایل سیستم‌های گنجانده شده را می‌دهند (پیش‌فرض‌ها شامل شبه‌فایل‌سیستم‌های (pseudofilesystems) مختلف نمی‌شوند). از آنجایی که اکسپورتر نود را به عنوان روت اجرا نخواهید کرد، باید اطمینان حاصل کنید که مجوزهای فایل به آن اجازه می‌دهد از فراخوانی سیستمی `statfs` روی نقاط مانت (mountpoints) مورد نظر شما استفاده کند.

تمام متریک‌های این کالکتور با پیشوند `node_filesystem_` شروع می‌شوند و دارای لیبل‌های `device`، `fstype` و `mountpoint` هستند:

```sh
# HELP node_filesystem_size_bytes Filesystem size in bytes.
# TYPE node_filesystem_size_bytes gauge
node_filesystem_size_bytes{device="/dev/sda5",fstype="ext4",mountpoint="/"} 9e+10
```

متریک‌های فایل سیستم عمدتاً بدیهی هستند. نکته ظریفی که باید از آن آگاه باشید تفاوت بین `node_filesystem_avail_bytes` و `node_filesystem_free_bytes` است. در فایل سیستم‌های یونیکس مقداری فضا برای کاربر روت رزرو شده است، تا زمانی که کاربران تمام فضای موجود را پر می‌کنند، همچنان بتوانند کارهایی انجام دهند. `node_filesystem_avail_bytes` فضای در دسترس کاربران است و هنگام تلاش برای محاسبه فضای دیسک استفاده شده، باید بر این اساس از فرمول زیر استفاده کنید:

```sh
node_filesystem_avail_bytes
/
node_filesystem_size_bytes
```

`node_filesystem_files` و `node_filesystem_files_free` تعداد آی‌نودها (inodes) و تعداد آی‌نودهای آزاد را نشان می‌دهند، که تقریباً تعداد فایل‌هایی است که فایل سیستم شما دارد. می‌توانید این را با `df -i` نیز مشاهده کنید.

### کالکتور Diskstats (آمار دیسک)

کالکتور diskstats متریک‌های ورودی/خروجی دیسک را از `/proc/diskstats` در معرض نمایش قرار می‌دهد. به طور پیش‌فرض، فلگ `--collector.diskstats.device-exclude` سعی می‌کند چیزهایی را که دیسک واقعی نیستند، مانند پارتیشن‌ها و دستگاه‌های لوپ‌بک (loopback)، مستثنی کند:

```
# HELP node_disk_io_now The number of I/Os currently in progress.
# TYPE node_disk_io_now gauge
node_disk_io_now{device="sda"} 0
```

تمام متریک‌ها دارای لیبل `device` هستند و تقریباً همه شمارنده (counter) می‌باشند، به شرح زیر:

node_disk_io_now
تعداد ورودی/خروجی‌های در حال انجام

node_disk_io_time_seconds_total
زمانی که ورودی/خروجی در حال انجام است افزایش می‌یابد

node_disk_read_bytes_total
بایت‌های خوانده شده توسط ورودی/خروجی‌ها

node_disk_read_time_seconds_total
زمان صرف شده توسط ورودی/خروجی‌های خواندن

node_disk_reads_completed_total
تعداد ورودی/خروجی‌های خواندن کامل شده

node_disk_written_bytes_total
بایت‌های نوشته شده توسط ورودی/خروجی‌ها

node_disk_write_time_seconds_total
زمان صرف شده توسط ورودی/خروجی‌های نوشتن

node_disk_writes_completed_total
تعداد ورودی/خروجی‌های نوشتن کامل شده

### کالکتور Netdev (دستگاه شبکه)

این‌ها عمدتاً همان معنایی را دارند که فکر می‌کنید، اما برای جزئیات بیشتر به مستندات هسته⁴ نگاهی بیندازید.

می‌توانید از `node_disk_io_time_seconds_total` برای محاسبه بهره‌وری ورودی/خروجی دیسک استفاده کنید، همانطور که توسط `iostat -x` نشان داده می‌شود:

`rate(node_disk_io_time_seconds_total[1m])`

می‌توانید میانگین زمان برای یک ورودی/خروجی خواندن را با استفاده از فرمول زیر محاسبه کنید:

```sh
rate(node_disk_read_time_seconds_total[1m])
/
rate(node_disk_reads_completed_total[1m])
```

### کالکتور Netdev (دستگاه شبکه)

کالکتور netdev متریک‌هایی درباره دستگاه‌های شبکه شما با پیشوند `node_network_` و یک لیبل `device` در معرض نمایش قرار می‌دهد:

```sh
# HELP node_network_receive_bytes_total Network device statistic receive_bytes.
# TYPE node_network_receive_bytes_total counter
node_network_receive_bytes_total{device="lo"} 8.3213967e+07
node_network_receive_bytes_total{device="wlan0"} 7.0854462e+07
```

`node_network_receive_bytes_total` و `node_network_transmit_bytes_total` متریک‌های اصلی هستند که به آن‌ها اهمیت خواهید داد زیرا می‌توانید پهنای باند ورودی و خروجی شبکه را با آن‌ها محاسبه کنید:

`rate(node_network_receive_bytes_total[1m])`

ممکن است به `node_network_receive_packets_total` و `node_network_transmit_packets_total` نیز علاقه‌مند باشید که به ترتیب بسته‌های ورودی و خروجی را ردیابی می‌کنند.

## کالکتور Meminfo (اطلاعات حافظه)

کالکتور meminfo تمام متریک‌های استاندارد مربوط به حافظه شما را با پیشوند `node_memory_` دارد. همه این‌ها از `/proc/meminfo` شما می‌آیند، و این اولین کالکتوری است که معناشناسی (semantics) کمی مبهم می‌شود. کالکتور کیلوبایت‌ها را به بایت‌های ترجیحی تبدیل می‌کند، اما فراتر از آن، به شما بستگی دارد که به اندازه کافی از مستندات و تجربه با موارد داخلی لینوکس بدانید تا بفهمید این متریک‌ها چه معنایی دارند:

```sh
# HELP node_memory_MemTotal_bytes Memory information field MemTotal.
# TYPE node_memory_MemTotal_bytes gauge
node_memory_MemTotal_bytes 3.285016576e+10
```

به عنوان مثال، `node_memory_MemTotal_bytes` مقدار کل⁵ حافظه فیزیکی در ماشین است – خوب و واضح. اما توجه داشته باشید که هیچ متریک حافظه استفاده شده (used memory) وجود ندارد، بنابراین باید به نوعی آن را محاسبه کنید و در نتیجه بفهمید چقدر حافظه از متریک‌های دیگر استفاده نشده است.

‏`node_memory_MemFree_bytes` مقدار حافظه‌ای است که توسط هیچ چیزی استفاده نمی‌شود، اما این بدان معنا نیست که تمام حافظه‌ای است که در اختیار دارید. در تئوری، کش صفحه (page cache) شما (`node_memory_Cached_bytes`) و همچنین بافرهای نوشتن (write buffers) شما (`node_memory_Buffers_bytes`) قابل بازپس‌گیری هستند، اما این می‌تواند بر عملکرد برخی برنامه‌ها تأثیر منفی بگذارد⁶. علاوه بر این، ساختارهای مختلف هسته دیگری وجود دارند که از حافظه استفاده می‌کنند مانند slab و جداول صفحه (page tables).

‏`node_memory_MemAvailable` یک روش اکتشافی (heuristic) از هسته برای میزان حافظه واقعاً در دسترس است، اما تنها در نسخه 3.14 لینوکس اضافه شد. اگر از هسته به اندازه کافی جدید استفاده می‌کنید، این متریکی است که می‌توانید برای تشخیص اتمام حافظه از آن استفاده کنید.

## کالکتور Hwmon (نظارت سخت‌افزار)

هنگامی که روی سخت‌افزار فیزیکی (bare metal) هستید، کالکتور hwmon متریک‌هایی مانند دما و سرعت فن‌ها را با پیشوند `node_hwmon_` ارائه می‌دهد. این همان اطلاعاتی است که می‌توانید با دستور `sensors` به دست آورید:

```sh
# HELP node_hwmon_sensor_label Label for given chip and sensor
# TYPE node_hwmon_sensor_label gauge
node_hwmon_sensor_label{chip="platform_coretemp_0", label="core_0",sensor="temp2"} 1
node_hwmon_sensor_label{chip="platform_coretemp_0", label="core_1",sensor="temp3"} 1
# HELP node_hwmon_temp_celsius Hardware monitor for temperature (input)
# TYPE node_hwmon_temp_celsius gauge
node_hwmon_temp_celsius{chip="platform_coretemp_0",sensor="temp1"} 42
node_hwmon_temp_celsius{chip="platform_coretemp_0",sensor="temp2"} 42
node_hwmon_temp_celsius{chip="platform_coretemp_0",sensor="temp3"} 41
```

‏`node_hwmon_temp_celsius` دمای اجزای مختلف شما است، که ممکن است دارای لیبل‌های سنسور⁷ نیز باشند که در `node_hwmon_sensor_label` نمایش داده می‌شوند.

در حالی که این مورد برای همه سخت‌افزارها صدق نمی‌کند، برای برخی⁸ برای درک اینکه سنسور چیست، به لیبل سنسور نیاز دارید. در متریک‌های قبلی، `temp3` هسته شماره ۱ CPU را نشان می‌دهد.

می‌توانید لیبل `label` از `node_hwmon_sensor_label` را با استفاده از `group_left` به `node_hwmon_temp_celsius` متصل کنید، که در «Many-to-One و group_left» در صفحه ۲۶۸ بیشتر مورد بحث قرار گرفته است:

```
node_hwmon_temp_celsius
* ignoring(label) group_left(label)
node_hwmon_sensor_label
```

### کالکتور Stat (آمار)

کالکتور stat کمی ترکیبی است، زیرا متریک‌هایی را از `/proc/stat`⁹ ارائه می‌دهد. `node_boot_time_seconds` زمانی است که هسته شروع به کار کرده است، که از آن می‌توانید محاسبه کنید چه مدت هسته فعال بوده است:

`time() - node_boot_time_seconds`
‏‏
‏`node_intr_total` تعداد وقفه‌های سخت‌افزاری (hardware interrupts) را که داشته‌اید نشان می‌دهد. به آن `node_interrupts_total` گفته نمی‌شود، زیرا این نام توسط کالکتور interrupts استفاده می‌شود که به دلیل کاردینالیتی بالا به طور پیش‌فرض غیرفعال است. متریک‌های دیگر مربوط به فرآیندها هستند. `node_forks_total` یک شمارنده برای تعداد فراخوانی‌های سیستمی fork است، `node_context_switches_total` تعداد تعویض‌های زمینه (context switches) است، در حالی که `node_procs_blocked` و `node_procs_running` تعداد فرآیندهایی را که مسدود (blocked) یا در حال اجرا (running) هستند، نشان می‌دهند.

### کالکتور Uname (اطلاعات Uname)

کالکتور uname یک متریک واحد به نام `node_uname_info` را در معرض نمایش قرار می‌دهد که قبلاً در «پنل Stat» در صفحه ۱۱۳ آن را دیده‌اید:

```sh
# HELP node_uname_info Labeled system information as provided by the uname system call.
# TYPE node_uname_info gauge
node_uname_info{domainname="(none)",machine="x86_64",nodename="kozo", release="4.4.0-101-generic",sysname="Linux", version="#124-Ubuntu SMP Fri Nov 10 18:29:59 UTC 2017"} 1
```

لیبل `nodename` نام میزبان (hostname) ماشین است، که ممکن است با لیبل هدف `instance` (به «لیبل‌های هدف» در صفحه ۱۵۳ مراجعه کنید) یا هر نام دیگری، مانند نام‌های موجود در DNS، که ممکن است برای آن داشته باشید، متفاوت باشد.

برای شمارش تعداد ماشین‌هایی که هر نسخه هسته را اجرا می‌کنند، می‌توانید از دستور زیر استفاده کنید:

`count by(release)(node_uname_info)`

## کالکتور OS (سیستم عامل)

کالکتور OS دو متریک `node_os_info` و `node_os_version` را در معرض نمایش قرار می‌دهد که اطلاعات سیستم عامل را به شما ارائه می‌دهند:

```sh
# HELP node_os_info A metric with a constant '1' value labeled by build_id, id, id_like, image_id, image_version, name, pretty_name, variant, variant_id, version, version_codename, version_id.
# TYPE node_os_info gauge
node_os_info{build_id="22.05.20220912.bf014ca",id="nixos", id_like="",image_id="",image_version="",name="NixOS", pretty_name="NixOS 22.05 (Quokka)",variant="", variant_id="",version="22.05 (Quokka)", version_codename="quokka",version_id="22.05"} 1
# HELP node_os_version Metric containing the major.minor part of the OS version.
# TYPE node_os_version gauge
node_os_version{id="nixos",id_like="",name="NixOS"} 22.05
```

برای شمارش تعداد ماشین‌هایی که هر نسخه توزیع (distro) را اجرا می‌کنند، می‌توانید از دستور زیر استفاده کنید:
`count by(name, version)(node_os_info)`

## کالکتور Loadavg (میانگین بار)

کالکتور loadavg میانگین بار ۱، ۵ و ۱۵ دقیقه‌ای را به ترتیب به صورت `node_load1`، `node_load5` و `node_load15` ارائه می‌دهد.

## کالکتور Loadavg (میانگین بار)

کالکتور loadavg میانگین بار ۱، ۵ و ۱۵ دقیقه‌ای را به ترتیب به صورت `node_load1`، `node_load5` و `node_load15` ارائه می‌دهد.

معنای این متریک در پلتفرم‌های مختلف متفاوت است و ممکن است آن چیزی نباشد که فکر می‌کنید. به عنوان مثال، در لینوکس، این فقط تعداد فرآیندهای منتظر در صف اجرا (run queue) نیست، بلکه شامل فرآیندهای غیرقابل وقفه (uninterruptible) مانند فرآیندهای منتظر ورودی/خروجی نیز می‌شود.

کالکتور loadavg میانگین بار ۱، ۵ و ۱۵ دقیقه‌ای را به ترتیب به صورت `node_load1`، `node_load5` و `node_load15` ارائه می‌دهد.

معنای این متریک در پلتفرم‌های مختلف متفاوت است و ممکن است آن چیزی نباشد که فکر می‌کنید. به عنوان مثال، در لینوکس، این فقط تعداد فرآیندهای منتظر در صف اجرا (run queue) نیست، بلکه شامل فرآیندهای غیرقابل وقفه (uninterruptible) مانند فرآیندهای منتظر ورودی/خروجی نیز می‌شود.

اگر هسته شما به اندازه کافی جدید است، توصیه می‌کنیم از کالکتور pressure استفاده کنید، همانطور که در «کالکتور Pressure» در صفحه ۱۳۳ توضیح داده شده است.

میانگین بار می‌تواند برای یک ایده سریع در مورد اینکه آیا یک ماشین اخیراً شلوغ‌تر شده است (برای برخی تعاریف از شلوغ‌تر) مفید باشد، اما انتخاب خوبی برای هشدار دادن نیستند. برای نگاهی دقیق‌تر، وبلاگ Brendan Gregg با عنوان «میانگین بار لینوکس: حل معما» را توصیه می‌کنیم.

> این یک عدد احمقانه است اما مردم فکر می‌کنند مهم است.
> — یک نظر در فایل loadavg.c لینوکس.

## کالکتور Pressure (فشار)

اطلاعات توقف فشار (Pressure Stall Information - PSI) در هسته لینوکس نسخه 4.20 معرفی شد. این متریک‌ها فشار منابع را برای سه منبع اندازه‌گیری می‌کنند: CPU، حافظه و ورودی/خروجی (I/O). این قابلیت باید در زمان کامپایل در هسته فعال شود.

ممکن است هسته شما با پشتیبانی PSI ساخته شده باشد اما به طور پیش‌فرض غیرفعال باشد، در این صورت می‌توانید در زمان بوت (boot) `psi=1` را در خط فرمان هسته (kernel command line) برای فعال کردن آن ارسال کنید. پنج متریک مختلف توسط کالکتور PSI در معرض نمایش قرار می‌گیرند:

```sh
# HELP node_pressure_cpu_waiting_seconds_total Total time in seconds that processes have waited for CPU time
# TYPE node_pressure_cpu_waiting_seconds_total counter
node_pressure_cpu_waiting_seconds_total 113.6605130
# HELP node_pressure_io_stalled_seconds_total Total time in seconds no process could make progress due to IO congestion
# TYPE node_pressure_io_stalled_seconds_total counter
node_pressure_io_stalled_seconds_total 8.630361
# HELP node_pressure_io_waiting_seconds_total Total time in seconds that processes have waited due to IO congestion
# TYPE node_pressure_io_waiting_seconds_total counter
node_pressure_io_waiting_seconds_total 9.609997
# HELP node_pressure_memory_stalled_seconds_total Total time in seconds no process could make progress
# TYPE node_pressure_memory_stalled_seconds_total counter
node_pressure_memory_stalled_seconds_total 0
# HELP node_pressure_memory_waiting_seconds_total Total time in seconds that processes have waited for memory
# TYPE node_pressure_memory_waiting_seconds_total counter
node_pressure_memory_waiting_seconds_total 0
```

متریک‌های `waiting` نشان‌دهنده کل زمان به ثانیه است که برخی وظایف منتظر بوده‌اند، و `stalled` به این معنی است که تمام وظایف به دلیل کمبود منابع به تأخیر افتاده‌اند. حافظه و I/O هر دو متریک `waiting` و `stalled` دارند، در حالی که CPU فقط `waiting` دارد. این به این دلیل است که CPU همیشه در حال اجرای یک فرآیند است.

از آنجایی که این‌ها شمارنده هستند، می‌توانید از آن‌ها با تابع `rate()` برای تعیین اینکه آیا برخی منابع بیش از حد بارگذاری شده‌اند استفاده کنید:

`rate(node_pressure_memory_waiting_seconds_total[1m])`

### کالکتور Textfile (فایل متنی)

کالکتور textfile کمی با کالکتورهایی که قبلاً به شما نشان دادیم متفاوت است. این کالکتور متریک‌ها را از هسته به دست نمی‌آورد، بلکه از فایل‌هایی که شما تولید می‌کنید، به دست می‌آورد. اکسپورتر نود قرار نیست به عنوان روت اجرا شود، بنابراین متریک‌هایی مانند متریک‌های SMART¹⁰ برای اجرای دستور `smartctl` به امتیازات روت (root privileges) نیاز دارند. علاوه بر متریک‌هایی که به روت نیاز دارند، برخی اطلاعات را فقط با اجرای دستوری مانند `iptables` می‌توانید به دست آورید. برای اطمینان‌پذیری، اکسپورتر نود فرآیندها را شروع نمی‌کند.

برای استفاده از کالکتور textfile، شما یک کران‌جاب (cronjob) ایجاد می‌کنید که به طور منظم دستوراتی مانند `smartctl` یا `iptables` را اجرا می‌کند، خروجی آن را به فرمت نمایش متنی پرومتئوس (Prometheus text exposition format) تبدیل می‌کند و آن را به صورت اتمی (atomically) در فایلی در یک دایرکتوری مشخص می‌نویسد. در هر اسکرِیپ (scrape)، اکسپورتر نود فایل‌های موجود در آن دایرکتوری را می‌خواند و متریک‌های آن‌ها را در خروجی خود قرار می‌دهد.

سرور پرومتئوس نمی‌تواند فایل‌های متنی را مستقیماً بخواند، بنابراین شما به نرم‌افزاری نیاز دارید تا فایل را به صورت HTTP در معرض نمایش قرار دهد. در حالی که می‌توانید از هر سرور HTTP استفاده کنید، اکسپورتر نود همچنین صحت متریک‌ها را بررسی می‌کند و قادر است متریک‌های حاصل از چندین فایل را در معرض نمایش قرار دهد.

می‌توانید از این کالکتور برای افزودن متریک‌های خود از طریق کران‌جاب‌ها استفاده کنید، یا می‌توانید اطلاعات ایستا (static) بیشتری داشته باشید که از فایل‌هایی می‌آیند که توسط سیستم مدیریت پیکربندی ماشین (machine configuration management system) شما نوشته شده‌اند تا برخی متریک‌های اطلاعاتی (info metrics) (مورد بحث در «Info» در صفحه ۹۶) را ارائه دهند، مانند اینکه چه نقش‌های Chef¹¹ را دارد، درباره ماشین.

مانند اکسپورتر نود به طور کلی، کالکتور textfile برای متریک‌های مربوط به یک ماشین در نظر گرفته شده است. به عنوان مثال، ممکن است متریک هسته‌ای وجود داشته باشد که اکسپورتر نود هنوز آن را در معرض نمایش قرار نمی‌دهد، یا برای دسترسی به روت نیاز دارد. ممکن است بخواهید متریک‌های سطح سیستم عامل بیشتری را ردیابی کنید، مانند اینکه آیا به‌روزرسانی‌های بسته در انتظار هستند یا نیاز به راه‌اندازی مجدد (reboot) وجود دارد. در حالی که از نظر فنی یک سرویس است تا یک متریک سیستم عامل، ثبت زمان آخرین تکمیل کارهای دسته‌ای (batch jobs) مانند پشتیبان‌گیری (backups) برای نود Cassandra¹² که روی ماشین اجرا می‌شود نیز استفاده خوبی از کالکتور textfile خواهد بود، زیرا علاقه شما به اینکه آیا پشتیبان‌گیری‌ها روی آن ماشین کار کرده‌اند، با از بین رفتن ماشین از بین می‌رود. یعنی نود Cassandra چرخه عمر (lifecycle) مشابهی با ماشین دارد.

کالکتور textfile نباید برای تبدیل پرومتئوس به حالت پوش (push) استفاده شود. همچنین نباید از کالکتور textfile به عنوان راهی برای گرفتن متریک‌ها از سایر اکسپورترها و برنامه‌های در حال اجرا روی ماشین و نمایش همه آنها در `/metrics` اکسپورتر نود استفاده کنید، بلکه باید پرومتئوس هر اکسپورتر و برنامه را به صورت جداگانه اسکرِیپ کند.

### استفاده از کالکتور Textfile

کالکتور textfile به طور پیش‌فرض فعال است، اما برای کار کردن آن باید فلگ خط فرمان `--collector.textfile.directory` را به اکسپورتر نود ارائه دهید. این باید به دایرکتوری اشاره کند که صرفاً برای این منظور استفاده می‌کنید تا از اشتباهات جلوگیری شود.

برای امتحان کردن این، باید یک دایرکتوری ایجاد کنید، یک فایل ساده در فرمت نمایش (exposition format) بنویسید (همانطور که در «فرمت نمایش متنی» در صفحه ۸۰ بحث شد)، و اکسپورتر نود را با پیکربندی برای استفاده از این دایرکتوری، همانطور که در مثال ۷-۱ نشان داده شده است، شروع کنید. کالکتور textfile فقط به فایل‌هایی با پسوند `.prom` نگاه می‌کند.

مثال ۷-۱. استفاده از کالکتور textfile با یک مثال ساده

```
hostname $ mkdir textfile
hostname $ echo example_metric 1 > textfile/example.prom
hostname $ ./node_exporter --collector.textfile.directory=$PWD/textfile
```

مثال ۷-۲ محتوای فایل ایجاد شده توسط مثال ۷-۱ را نشان می‌دهد.
مثال ۷-۲. محتوای textfile/example.prom

```
example_metric 1
```

اگر به `/metrics` اکسپورتر نود نگاه کنید، اکنون متریک خود را خواهید دید:

```
# HELP example_metric Metric read from /some/path/textfile/example.prom
# TYPE example_metric untyped
example_metric 1
```

اگر هیچ HELP ارائه نشود، کالکتور textfile یکی برای شما ارائه می‌دهد. اگر همان متریک را در چندین فایل (البته با لیبل‌های مختلف) قرار می‌دهید، باید HELP یکسانی برای هر کدام ارائه دهید، زیرا در غیر این صورت HELP نامطابق باعث خطا می‌شود.

معمولاً شما فایل‌های `.prom` را با یک کران‌جاب ایجاد و به‌روز می‌کنید. از آنجایی که یک اسکرِیپ می‌تواند در هر زمانی اتفاق بیفتد، مهم است که اکسپورتر نود فایل‌هایی را که تا حدی نوشته شده‌اند نبیند. برای این منظور باید ابتدا در یک فایل موقت (temporary file) در همان دایرکتوری بنویسید و سپس فایل کامل را به نام فایل نهایی منتقل (move) کنید.

مثال ۷-۳ یک کران‌جاب را نشان می‌دهد که به کالکتور textfile خروجی می‌دهد. این کران‌جاب متریک‌ها را در یک فایل موقت¹⁵ ایجاد می‌کند و نام آنها را به نام فایل نهایی تغییر می‌دهد. این یک مثال پیش پا افتاده است که از دستورات کوتاه استفاده می‌کند، اما در بیشتر موارد استفاده واقعی، می‌خواهید یک اسکریپت برای خوانا نگه داشتن موارد ایجاد کنید.

مثال ۷-۳. فایل `/etc/crontab` که تعداد خطوط در `/etc/shadow` را به عنوان متریک `shadow_entries` با استفاده از کالکتور textfile در معرض نمایش قرار می‌دهد.

`TEXTFILE=/path/to/textfile/directory`

```sh
# این باید همه در یک خط باشد
*/5 * * * * root (echo -n 'shadow_entries '; grep -c . /etc/shadow) > $TEXTFILE/shadow.prom.$$ && mv $TEXTFILE/shadow.prom.$$ $TEXTFILE/shadow.prom
```

تعدادی اسکریپت نمونه برای استفاده با کالکتور textfile در مخزن GitHub اسکریپت‌های نمونه کالکتور textfile موجود است.

### مُهرهای زمانی (Timestamps)

در حالی که فرمت نمایش از مُهرهای زمانی پشتیبانی می‌کند، نمی‌توانید از آنها با کالکتور textfile استفاده کنید. این به این دلیل است که از نظر معنایی منطقی نیست، زیرا متریک‌های شما با همان مُهر زمانی سایر متریک‌های اسکرِیپ ظاهر نمی‌شوند. در عوض، mtime¹⁶ (زمان آخرین تغییر) فایل در متریک `node_textfile_mtime_seconds` در دسترس شماست. می‌توانید از این برای هشدار دادن در مورد کار نکردن کران‌جاب‌های خود استفاده کنید، زیرا اگر این مقدار مربوط به زمان خیلی دوری باشد، می‌تواند نشان‌دهنده یک مشکل باشد:

```sh
# HELP node_textfile_mtime_seconds Unixtime mtime of textfiles successfully read.
# TYPE node_textfile_mtime_seconds gauge
node_textfile_mtime_seconds{file="example.prom"} 1.516205651e+09
```

اکنون که اکسپورتر نود را در حال اجرا دارید، بیایید ببینیم چگونه می‌توانید به پرومتئوس در مورد تمام ماشین‌هایی که آن را روی آن‌ها اجرا می‌کنید، اطلاع دهید.
