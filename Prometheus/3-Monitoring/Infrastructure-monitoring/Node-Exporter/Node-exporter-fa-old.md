

ءNode Exporter1 احتمالاً یکی از اولین صادرکننده‌هایی (exporters) است که از آن استفاده خواهید کرد، همانطور که قبلاً در فصل 2 دیدیم. این صادرکننده متریک‌های سطح ماشین، عمدتاً از هسته سیستم عامل شما، مانند CPU، حافظه، فضای دیسک، ورودی/خروجی دیسک، پهنای باند شبکه و دمای مادربرد را ارائه می‌دهد. Node Exporter با سیستم‌های یونیکس استفاده می‌شود؛ کاربران ویندوز باید به جای آن از Windows Exporter2 استفاده کنند. Node Exporter فقط برای نظارت بر خود ماشین در نظر گرفته شده است، نه فرآیندها یا خدمات فردی روی آن. سیستم‌های نظارتی دیگر اغلب چیزی دارند که ما دوست داریم آن را uberagent بنامیم؛ یعنی یک فرآیند واحد که همه چیز را روی ماشین نظارت می‌کند. در معماری پرومتئوس، هر یک از خدمات شما متریک‌های خود را، در صورت نیاز با استفاده از یک صادرکننده، ارائه می‌دهند که سپس مستقیماً توسط پرومتئوس جمع‌آوری می‌شود. این کار از تبدیل شدن uberagent به یک گلوگاه عملیاتی یا عملکردی جلوگیری می‌کند و شما را قادر می‌سازد بیشتر به خدمات پویا فکر کنید تا ماشین‌ها. دستورالعمل‌هایی که هنگام ایجاد متریک‌ها با ابزار دقیق مستقیم (direct instrumentation) استفاده می‌کنید، مانند مواردی که در "نام متریک‌هایم را چه بگذارم؟" در صفحه 60 بحث شد، نسبتاً واضح هستند. این در مورد صادرکننده‌ها صدق نمی‌کند، جایی که طبق تعریف، داده‌ها از منبعی می‌آیند که با دستورالعمل‌های پرومتئوس طراحی نشده است. بسته به حجم و کیفیت متریک‌ها، توسعه‌دهندگان صادرکننده باید بین تلاش مهندسی و دریافت متریک‌های عالی، مصالحه‌ای انجام دهند. در مورد لینوکس، هزاران متریک ارائه می‌شود. برخی از آنها به خوبی مستند و درک شده‌اند، مانند استفاده از CPU؛ برخی دیگر، مانند استفاده از حافظه، با تغییر پیاده‌سازی، از نسخه هسته به نسخه هسته دیگر متفاوت بوده‌اند. حتی متریک‌هایی را پیدا خواهید کرد که کاملاً مستند نشده‌اند، جایی که باید کد منبع هسته را بخوانید تا سعی کنید بفهمید چه کاری انجام می‌دهند. Node Exporter به گونه‌ای طراحی شده است که به عنوان یک کاربر غیر ریشه (nonroot) اجرا شود و باید مستقیماً روی ماشین به همان روشی که یک دیمون سیستم مانند sshd یا cron را اجرا می‌کنید، اجرا شود.

در حالی که اجرای Node Exporter در داخل Docker3 امکان‌پذیر است، باید از برخی حجم‌ها و پارامترهای خط فرمان (--path.procfs، --path.rootfs، --path.sysfs) برای نصب سیستم فایل میزبان در داخل کانتینر استفاده کنید. در صورت امکان، Node Exporter را به عنوان یک سرویس روی گره، بدون Docker اجرا کنید. Docker سعی می‌کند یک کانتینر را از عملکرد داخلی ماشین جدا کند، که با Node Exporter که سعی در دسترسی به آن عملکردهای داخلی دارد، خوب کار نمی‌کند.

برخلاف اکثر صادرکننده‌های دیگر، به دلیل تنوع گسترده متریک‌های موجود از سیستم‌های عامل، Node Exporter به شما امکان می‌دهد دسته‌های متریک‌هایی را که واکشی می‌کند، پیکربندی کنید. می‌توانید این کار را با پرچم‌های خط فرمان مانند --collector.wifi انجام دهید، که جمع‌آورنده WiFi را فعال می‌کند، و --no-collector.wifi، که آن را غیرفعال می‌کند. --collector.disable-defaults تمام جمع‌آورنده‌ها را به جز مواردی که به صراحت به عنوان پرچم‌های خط فرمان فعال شده‌اند، غیرفعال می‌کند. پیش‌فرض‌های منطقی تنظیم شده‌اند، بنابراین این چیزی نیست که هنگام شروع کار نگران آن باشید. هسته‌های مختلف متریک‌های مختلفی را ارائه می‌دهند، زیرا، به عنوان مثال، لینوکس و FreeBSD کارها را به روش‌های مختلف انجام می‌دهند. ممکن است با بازسازی Node Exporter، متریک‌ها در طول زمان بین جمع‌آورنده‌ها جابجا شوند. اگر از یک سیستم یونیکس متفاوت استفاده می‌کنید، خواهید دید که متریک‌ها و جمع‌آورنده‌های ارائه شده متفاوت هستند. در این فصل، برخی از متریک‌های کلیدی ارائه شده توسط Node Exporter نسخه 1.4.0 با هسته لینوکس 5.18.0 را توضیح می‌دهیم. این فهرست جامعی از متریک‌های موجود نیست. مانند اکثر صادرکننده‌ها و برنامه‌ها، می‌خواهید مسیر /metrics را بررسی کنید تا ببینید چه چیزی در دسترس است. می‌توانید عبارات PromQL مثال را با استفاده از تنظیمات خود از فصل 2 امتحان کنید.

بله، در اینجا ترجمه فارسی متن شما آمده است:

## جمع‌آورنده CPU (CPU Collector)

متریک اصلی از جمع‌آورنده CPU، `node_cpu_seconds_total` است که یک شمارنده است که نشان می‌دهد هر CPU چه مدت زمانی را در هر حالت صرف کرده است. برچسب‌ها `cpu` و `mode` هستند:

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

برای هر CPU، حالت‌ها در مجموع در هر ثانیه یک ثانیه افزایش می‌یابند. این به شما امکان می‌دهد نسبت زمان بیکاری را در تمام CPUها با استفاده از عبارت PromQL زیر محاسبه کنید:

avg without(cpu, mode)(rate(node_cpu_seconds_total{mode="idle"}[1m]))

این کار به این صورت انجام می‌شود که زمان بیکاری را در هر ثانیه در هر CPU محاسبه می‌کند و سپس آن را در تمام CPUهای دستگاه میانگین می‌گیرد. می‌توانید این را تعمیم دهید تا نسبت زمان صرف شده در هر حالت را برای یک دستگاه با استفاده از این محاسبه کنید:

avg without(cpu)(rate(node_cpu_seconds_total[1m]))

استفاده از CPU توسط مهمان‌ها (یعنی ماشین‌های مجازی که تحت هسته اجرا می‌شوند) قبلاً در حالت‌های کاربر و nice گنجانده شده است. می‌توانید زمان مهمان را به طور جداگانه در متریک node_cpu_guest_seconds_total مشاهده کنید.

## جمع‌آورنده سیستم فایل (Filesystem Collector)

جمع‌آورنده سیستم فایل، همانطور که از نامش پیداست، متریک‌هایی را در مورد سیستم‌های فایل نصب شده شما جمع‌آوری می‌کند، همانطور که از دستور `df` دریافت می‌کنید. پرچم‌های `--collector.file system.mount-points-exclude` و `--collector.filesystem.fs-types-exclude` امکان محدود کردن سیستم‌های فایلی را که شامل می‌شوند، فراهم می‌کنند (پیش‌فرض‌ها سیستم‌های فایل شبه (pseudofilesystems) مختلف را حذف می‌کنند). از آنجا که Node Exporter را به عنوان ریشه اجرا نمی‌کنید، باید اطمینان حاصل کنید که مجوزهای فایل به آن اجازه می‌دهد تا از فراخوانی سیستم `statfs` در نقاط نصب مورد نظر شما استفاده کند.

تمام متریک‌های این جمع‌آورنده با `node_filesystem_` پیشوند می‌شوند و دارای برچسب‌های `device`، `fstype` و `mountpoint` هستند:

Bash

```
# HELP node_filesystem_size_bytes Filesystem size in bytes.
# TYPE node_filesystem_size_bytes gauge
node_filesystem_size_bytes{device="/dev/sda5",fstype="ext4",mountpoint="/"} 9e+10
```

متریک‌های سیستم فایل تا حد زیادی خود توضیحی هستند. تنها نکته ظریفی که باید از آن آگاه باشید، تفاوت بین `node_filesystem_avail_bytes` و `node_filesystem_free_bytes` است. در سیستم‌های فایل یونیکس، مقداری فضا برای کاربر ریشه رزرو می‌شود، به طوری که آنها هنوز هم می‌توانند وقتی کاربران تمام فضای موجود را پر می‌کنند، کارها را انجام دهند. `node_filesystem_avail_bytes` فضایی است که برای کاربران در دسترس است و هنگام تلاش برای محاسبه فضای دیسک استفاده شده، باید بر این اساس از این استفاده کنید:

Bash

```
node_filesystem_avail_bytes
/
node_filesystem_size_bytes
```

`node_filesystem_files` و `node_filesystem_files_free` تعداد inodeها و تعداد آزاد آنها را نشان می‌دهند، که تقریباً تعداد فایل‌هایی است که سیستم فایل شما دارد. همچنین می‌توانید این را با `df -i` مشاهده کنید.

----

بله، در اینجا ترجمه فارسی متن شما آمده است:

### جمع‌آورنده Diskstats

جمع‌آورنده diskstats متریک‌های ورودی/خروجی دیسک را از /proc/diskstats ارائه می‌دهد. به طور پیش‌فرض، پرچم `--collector.diskstats.device-exclude` سعی می‌کند مواردی را که دیسک واقعی نیستند، مانند پارتیشن‌ها و دستگاه‌های loopback، حذف کند:

```
# HELP node_disk_io_now The number of I/Os currently in progress.
# TYPE node_disk_io_now gauge
node_disk_io_now{device="sda"} 0
```

همه متریک‌ها دارای برچسب دستگاه هستند و تقریباً همه شمارنده هستند، به شرح زیر:

- `node_disk_io_now`: تعداد ورودی/خروجی‌های در حال انجام
- `node_disk_io_time_seconds_total`: زمانی که ورودی/خروجی در حال انجام است افزایش می‌یابد
- `node_disk_read_bytes_total`: بایت‌های خوانده شده توسط ورودی/خروجی‌ها
- `node_disk_read_time_seconds_total`: زمان صرف شده توسط ورودی/خروجی‌های خواندن
- `node_disk_reads_completed_total`: تعداد ورودی/خروجی‌های کامل شده
- `node_disk_written_bytes_total`: بایت‌های نوشته شده توسط ورودی/خروجی‌ها
- `node_disk_write_time_seconds_total`: زمان صرف شده توسط ورودی/خروجی‌های نوشتن
- `node_disk_writes_completed_total`: تعداد ورودی/خروجی‌های نوشتن کامل شده

می‌توانید از `node_disk_io_time_seconds_total` برای محاسبه میزان استفاده از ورودی/خروجی دیسک، همانطور که توسط iostat -x نشان داده می‌شود، استفاده کنید:

`rate(node_disk_io_time_seconds_total[1m])`

می‌توانید میانگین زمان برای ورودی/خروجی خواندن را با این محاسبه کنید:

Bash

```
rate(node_disk_read_time_seconds_total[1m])
/
rate(node_disk_reads_completed_total[1m])
```

### جمع‌آورنده Netdev

جمع‌آورنده netdev متریک‌هایی را در مورد دستگاه‌های شبکه شما با پیشوند `node_network_` و برچسب دستگاه ارائه می‌دهد:

Bash

```
# HELP node_network_receive_bytes_total Network device statistic receive_bytes.
# TYPE node_network_receive_bytes_total counter
node_network_receive_bytes_total{device="lo"} 8.3213967e+07
node_network_receive_bytes_total{device="wlan0"} 7.0854462e+07
```

`node_network_receive_bytes_total` و `node_network_transmit_bytes_total` متریک‌های اصلی هستند که به آنها اهمیت می‌دهید، زیرا می‌توانید پهنای باند شبکه ورودی و خروجی را با آنها محاسبه کنید:

`rate(node_network_receive_bytes_total[1m])`

همچنین ممکن است به `node_network_receive_packets_total` و `node_network_transmit_packets_total` علاقه مند باشید، که به ترتیب بسته‌های ورودی و خروجی را ردیابی می‌کنند.

### جمع‌آورنده Meminfo

جمع‌آورنده meminfo تمام متریک‌های استاندارد مربوط به حافظه شما را با پیشوند `node_memory_` دارد. همه اینها از /proc/meminfo شما می‌آیند، و این اولین جمع‌آورنده‌ای است که در آن معناشناسی کمی مبهم می‌شود. جمع‌آورنده کیلوبایت‌ها را به بایت‌های ترجیحی تبدیل می‌کند، اما فراتر از آن، این به شما بستگی دارد که از مستندات و تجربه با درونیات لینوکس به اندازه کافی بدانید تا بفهمید این متریک‌ها چه معنایی دارند:

Bash

```
# HELP node_memory_MemTotal_bytes Memory information field MemTotal.
# TYPE node_memory_MemTotal_bytes gauge
node_memory_MemTotal_bytes 3.285016576e+10
```

به عنوان مثال، `node_memory_MemTotal_bytes` کل مقدار حافظه فیزیکی در دستگاه است - خوب و واضح. اما توجه داشته باشید که هیچ متریک حافظه استفاده شده وجود ندارد، بنابراین باید به نوعی آن را محاسبه کنید و در نتیجه مقدار حافظه‌ای که از سایر متریک‌ها استفاده نمی‌شود.

‏`node_memory_MemFree_bytes` مقدار حافظه‌ای است که توسط هیچ چیزی استفاده نمی‌شود، اما این بدان معنا نیست که تمام حافظه‌ای است که می‌توانید صرفه‌جویی کنید. در تئوری، کش صفحه شما (`node_memory_Cached_bytes`) و بافرهای نوشتن شما (`node_memory_Buffers_bytes`) قابل بازیابی هستند، اما این می‌تواند بر عملکرد برخی برنامه‌ها تأثیر منفی بگذارد. علاوه بر این، ساختارهای مختلف هسته دیگری نیز وجود دارند که از حافظه استفاده می‌کنند، مانند slab و جداول صفحه.

`node_memory_MemAvailable` یک اکتشافی از هسته برای میزان حافظه واقعاً در دسترس است، اما فقط در نسخه 3.14 لینوکس اضافه شد. اگر یک هسته به اندازه کافی جدید را اجرا می‌کنید، این متریکی است که می‌توانید برای تشخیص اتمام حافظه از آن استفاده کنید.

### جمع‌آورنده Hwmon

هنگام استفاده از سخت‌افزار واقعی، جمع‌آورنده hwmon متریک‌هایی مانند دما و سرعت فن را با پیشوند `node_hwmon_` ارائه می‌دهد. این همان اطلاعاتی است که می‌توانید با دستور sensors بدست آورید:

Bash

```
# HELP node_hwmon_sensor_label Label for given chip and sensor
# TYPE node_hwmon_sensor_label gauge
node_hwmon_sensor_label{chip="platform_coretemp_0",
label="core_0",sensor="temp2"} 1
node_hwmon_sensor_label{chip="platform_coretemp_0",
label="core_1",sensor="temp3"} 1
# HELP node_hwmon_temp_celsius Hardware monitor for temperature (input)
# TYPE node_hwmon_temp_celsius gauge
node_hwmon_temp_celsius{chip="platform_coretemp_0",sensor="temp1"} 42
node_hwmon_temp_celsius{chip="platform_coretemp_0",sensor="temp2"} 42
node_hwmon_temp_celsius{chip="platform_coretemp_0",sensor="temp3"} 41
```

`node_hwmon_temp_celsius` دمای اجزای مختلف شما است، که ممکن است برچسب‌های حسگر نیز در `node_hwmon_sensor_label` ارائه شوند.

در حالی که این مورد برای همه سخت‌افزارها صدق نمی‌کند، برای برخی باید از برچسب حسگر برای درک اینکه حسگر چیست استفاده کنید. در متریک‌های قبلی، temp3 نشان دهنده هسته CPU شماره 1 است.

می‌توانید برچسب label را از `node_hwmon_sensor_label` به `node_hwmon_temp_celsius` با استفاده از group_left پیوند دهید، که در "چند به یک و group_left" در صفحه 268 بیشتر مورد بحث قرار گرفته است:

```
node_hwmon_temp_celsius
* ignoring(label) group_left(label)
node_hwmon_sensor_label
```

### جمع‌آورنده Stat

جمع‌آورنده stat کمی ترکیبی است، زیرا متریک‌هایی را از /proc/stat ارائه می‌دهد. `node_boot_time_seconds` زمانی است که هسته شروع به کار کرد، که از آن می‌توانید مدت زمان روشن بودن هسته را محاسبه کنید:

`time() - node_boot_time_seconds`

`node_intr_total` تعداد وقفه های سخت افزاری را که داشته اید نشان می دهد. `node_interrupts_total` نامیده نمی شود، زیرا توسط جمع آورنده interrupts استفاده می شود، که به دلیل کاردینالیتی بالا به طور پیش فرض غیرفعال است. سایر متریک ها مربوط به فرآیندها هستند. `node_forks_total` یک شمارنده برای تعداد فراخوانی های سیستمی fork است، `node_context_switches_total` تعداد سوئیچ های زمینه است، در حالی که `node_procs_blocked` و `node_procs_running` تعداد فرآیندهایی را که مسدود یا در حال اجرا هستند نشان می دهند.

### جمع‌آورنده Uname

جمع‌آورنده uname یک متریک واحد، `node_uname_info` را ارائه می‌دهد، که قبلاً در "پانل Stat" در صفحه 113 دیدید:

Bash

```
# HELP node_uname_info Labeled system information as provided by the uname
system call.
# TYPE node_uname_info gauge
node_uname_info{domainname="(none)",machine="x86_64",nodename="kozo",
release="4.4.0-101-generic",sysname="Linux",
version="#124-Ubuntu SMP Fri Nov 10 18:29:59 UTC 2017"} 1
```

برچسب nodename نام میزبان دستگاه است، که ممکن است با برچسب هدف نمونه (به "برچسب های هدف" در صفحه 153 مراجعه کنید) یا هر نام دیگری، مانند DNS، که ممکن است برای آن داشته باشید، متفاوت باشد.

برای شمارش تعداد ماشین هایی که نسخه هسته را اجرا می کنند، می توانید از:

`count by(release)(node_uname_info)`

### جمع آورنده OS

جمع آورنده OS دو متریک `node_os_info` و `node_os_version` را ارائه می دهد که اطلاعات سیستم عامل را در اختیار شما قرار می دهد:

Bash

```
# HELP node_os_info A metric with a constant '1' value labeled by
build_id, id, id_like, image_id, image_version, name,
pretty_name, variant, variant_id, version, version_codename,
version_id.
# TYPE node_os_info gauge
node_os_info{build_id="22.05.20220912.bf014ca",id="nixos",
id_like="",image_id="",image_version="",name="NixOS",
pretty_name="NixOS 22.05 (Quokka)",variant="",
variant_id="",version="22.05 (Quokka)",
version_codename="quokka",version_id="22.05"} 1
# HELP node_os_version Metric containing the major.minor
part of the OS version.
# TYPE node_os_version gauge
node_os_version{id="nixos",id_like="",name="NixOS"} 22.05
```

برای شمارش تعداد ماشین هایی که نسخه توزیع را اجرا می کنند، می توانید از:

`count by(name, version)(node_os_info)`

### جمع‌آورنده Loadavg

جمع‌آورنده loadavg میانگین بارهای 1، 5 و 15 دقیقه را به ترتیب به عنوان `node_load1`، `node_load5` و `node_load15` ارائه می‌کند.


----------



بله، در اینجا ترجمه فارسی متن شما آمده است:

## جمع‌آورنده Loadavg

جمع‌آورنده loadavg میانگین بارهای 1، 5 و 15 دقیقه را به ترتیب به عنوان `node_load1`، `node_load5` و `node_load15` ارائه می‌کند.

معنی این متریک در پلتفرم‌های مختلف متفاوت است و ممکن است آن چیزی نباشد که شما فکر می‌کنید. برای مثال، در لینوکس، این فقط تعداد فرآیندهایی نیست که در صف اجرا منتظر هستند، بلکه فرآیندهای غیرقابل قطع شدن مانند فرآیندهایی که منتظر ورودی/خروجی هستند نیز شامل می‌شود.

جمع‌آورنده loadavg میانگین بارهای 1، 5 و 15 دقیقه را به ترتیب به عنوان `node_load1`، `node_load5` و `node_load15` ارائه می‌کند.

معنی این متریک در پلتفرم‌های مختلف متفاوت است و ممکن است آن چیزی نباشد که شما فکر می‌کنید. برای مثال، در لینوکس، این فقط تعداد فرآیندهایی نیست که در صف اجرا منتظر هستند، بلکه فرآیندهای غیرقابل قطع شدن مانند فرآیندهایی که منتظر ورودی/خروجی هستند نیز شامل می‌شود.

اگر هسته شما به اندازه کافی جدید است، توصیه می‌کنیم از جمع‌آورنده فشار (pressure collector) استفاده کنید، همانطور که در "جمع‌آورنده فشار" در صفحه 133 توضیح داده شده است.

میانگین بارها می‌توانند برای یک ایده سریع مفید باشند که آیا یک دستگاه اخیراً شلوغ‌تر شده است (برای برخی تعاریف شلوغ‌تر)، اما انتخاب خوبی برای هشدار دادن نیستند. برای نگاهی دقیق‌تر، وبلاگ برندان گرگ، "میانگین بارهای لینوکس: حل راز" را توصیه می‌کنیم.

این یک عدد احمقانه است اما مردم فکر می‌کنند مهم است.

—نظری در loadavg.c لینوکس.

## جمع‌آورنده فشار (Pressure Collector)

اطلاعات فشار استال (PSI) در هسته لینوکس 4.20 معرفی شد. این متریک‌ها فشار منابع را برای سه منبع اندازه‌گیری می‌کنند: CPU، حافظه و ورودی/خروجی. باید در زمان کامپایل در هسته فعال شود.

هسته شما ممکن است با پشتیبانی PSI ساخته شده باشد، اما ممکن است به طور پیش‌فرض غیرفعال باشد، در این صورت می‌توانید `psi=1` را در خط فرمان هسته در زمان بوت برای فعال کردن آن ارسال کنید. پنج متریک مختلف توسط جمع‌آورنده PSI ارائه می‌شود:

Bash

```
# HELP node_pressure_cpu_waiting_seconds_total
Total time in seconds that processes have waited for CPU time
# TYPE node_pressure_cpu_waiting_seconds_total counter
node_pressure_cpu_waiting_seconds_total 113.6605130
# HELP node_pressure_io_stalled_seconds_total
Total time in seconds no process could make progress due to IO congestion
# TYPE node_pressure_io_stalled_seconds_total counter
node_pressure_io_stalled_seconds_total 8.630361
# HELP node_pressure_io_waiting_seconds_total
Total time in seconds that processes have waited due to IO congestion
# TYPE node_pressure_io_waiting_seconds_total counter
node_pressure_io_waiting_seconds_total 9.609997
# HELP node_pressure_memory_stalled_seconds_total
Total time in seconds no process could make progress
# TYPE node_pressure_memory_stalled_seconds_total counter
node_pressure_memory_stalled_seconds_total 0
# HELP node_pressure_memory_waiting_seconds_total
Total time in seconds that processes have waited for memory
# TYPE node_pressure_memory_waiting_seconds_total counter
node_pressure_memory_waiting_seconds_total 0
```

متریک‌های waiting کل مقدار ثانیه‌هایی را نشان می‌دهند که برخی وظایف منتظر بوده‌اند، و stalled به این معنی است که همه وظایف به دلیل کمبود منابع به تأخیر افتاده‌اند. حافظه و ورودی/خروجی هر دو متریک waiting و stalled را دارند، در حالی که CPU فقط waiting را دارد. این به این دلیل است که CPU همیشه یک فرآیند را اجرا می‌کند.

از آنجا که اینها شمارنده‌ها هستند، می‌توانید از آنها با تابع `rate()` استفاده کنید تا تعیین کنید که آیا برخی منابع بیش از حد بارگذاری شده‌اند یا خیر:

`rate(node_pressure_memory_waiting_seconds_total[1m])`

### جمع‌آورنده Textfile

جمع‌آورنده textfile کمی با جمع‌آورنده‌هایی که قبلاً به شما نشان دادیم متفاوت است. این متریک‌ها را از هسته دریافت نمی‌کند، بلکه از فایل‌هایی که شما تولید می‌کنید دریافت می‌کند. Node Exporter برای اجرا به عنوان ریشه در نظر گرفته نشده است، بنابراین متریک‌هایی مانند موارد حاصل از SMART10 برای اجرای دستور smartctl به امتیازات ریشه نیاز دارند. علاوه بر متریک‌هایی که به ریشه نیاز دارند، فقط می‌توانید برخی اطلاعات را با اجرای دستوری مانند iptables بدست آورید. برای قابلیت اطمینان، Node Exporter فرآیندها را شروع نمی‌کند.

برای استفاده از جمع‌آورنده textfile، یک cronjob ایجاد می‌کنید که به طور منظم دستوراتی مانند smartctl یا iptables را اجرا می‌کند، خروجی آن را به فرمت ارائه متن پرومتئوس تبدیل می‌کند و به طور اتمی آن را در فایلی در یک دایرکتوری خاص می‌نویسد. در هر scrape، Node Exporter فایل‌های آن دایرکتوری را می‌خواند و متریک‌های آنها را در خروجی خود قرار می‌دهد.

سرور پرومتئوس نمی‌تواند فایل‌های متنی را مستقیماً بخواند، بنابراین به نرم‌افزاری نیاز دارید که فایل را به عنوان HTTP ارائه دهد. در حالی که می‌توانید از هر سرور HTTP استفاده کنید، Node Exporter همچنین بررسی می‌کند که متریک‌ها صحیح هستند و می‌تواند متریک‌های حاصل از چندین فایل را ارائه دهد.

می‌توانید از این جمع‌آورنده برای افزودن متریک‌های خود از طریق cronjobها استفاده کنید، یا می‌توانید اطلاعات ثابت‌تری داشته باشید که از فایل‌هایی که توسط سیستم مدیریت پیکربندی دستگاه شما نوشته شده‌اند، برای ارائه برخی متریک‌های اطلاعاتی (که در "اطلاعات" در صفحه 96 بحث شده است)، مانند نقش‌های Chef11 که دارد، در مورد دستگاه، استفاده کنید.

مانند Node Exporter به طور کلی، جمع‌آورنده textfile برای متریک‌های مربوط به یک دستگاه در نظر گرفته شده است. به عنوان مثال، ممکن است برخی متریک‌های هسته وجود داشته باشد که Node Exporter هنوز ارائه نمی‌دهد، یا نیاز به دسترسی ریشه دارد. ممکن است بخواهید متریک‌های سطح سیستم عامل بیشتری را ردیابی کنید، مانند اینکه آیا به‌روزرسانی‌های بسته‌های معلق یا راه‌اندازی مجدد وجود دارد یا خیر. در حالی که از نظر فنی یک سرویس است نه یک متریک سیستم عامل، ثبت زمانی که کارهای دسته‌ای مانند پشتیبان‌گیری برای گره Cassandra12 در حال اجرا روی دستگاه آخرین بار کامل شده‌اند نیز استفاده خوبی از جمع‌آورنده textfile خواهد بود، زیرا علاقه شما به اینکه آیا پشتیبان‌گیری‌ها روی آن دستگاه کار کرده‌اند، با از بین رفتن دستگاه از بین می‌رود. یعنی گره Cassandra همان چرخه عمر دستگاه را دارد.

جمع‌آورنده textfile نباید برای تلاش برای تبدیل پرومتئوس به push استفاده شود. همچنین نباید از جمع‌آورنده textfile به عنوان راهی برای گرفتن متریک‌ها از سایر صادرکننده‌ها و برنامه‌های در حال اجرا روی دستگاه و ارائه همه آنها در /metrics Node Exporter استفاده کنید، بلکه پرومتئوس باید هر صادرکننده و برنامه را به صورت جداگانه scrape کند.

------


بله، در اینجا ترجمه فارسی متن شما آمده است:

### استفاده از جمع‌آورنده Textfile

جمع‌آورنده textfile به طور پیش‌فرض فعال است، اما برای اینکه کار کند، باید پرچم خط فرمان `--collector.textfile.directory` را به Node Exporter ارائه دهید. این باید به دایرکتوری‌ای اشاره کند که فقط برای این منظور استفاده می‌کنید تا از اشتباهات جلوگیری شود.

برای امتحان این موضوع، باید یک دایرکتوری ایجاد کنید، یک فایل ساده در قالب ارائه (همانطور که در "قالب ارائه متن" در صفحه 80 بحث شد) بنویسید و Node Exporter را با پیکربندی برای استفاده از این دایرکتوری شروع کنید، همانطور که در مثال 7-1 نشان داده شده است. جمع‌آورنده textfile فقط فایل‌هایی با پسوند .prom را بررسی می‌کند.

مثال 7-1. استفاده از جمع‌آورنده textfile با یک مثال ساده

```
hostname $ mkdir textfile
hostname $ echo example_metric 1 > textfile/example.prom
hostname $ ./node_exporter --collector.textfile.directory=$PWD/textfile
```

مثال 7-2 محتوای فایل ایجاد شده توسط مثال 7-1 را نشان می‌دهد.

مثال 7-2. محتوای textfile/example.prom

```
example_metric 1
```

اگر به /metrics Node Exporter نگاه کنید، اکنون متریک خود را خواهید دید:

```
# HELP example_metric Metric read from /some/path/textfile/example.prom
# TYPE example_metric untyped
example_metric 1
```

اگر هیچ HELP ارائه نشود، جمع‌آورنده textfile یکی را برای شما ارائه می‌کند. اگر همان متریک را در چندین فایل قرار می‌دهید (البته با برچسب‌های مختلف)، باید همان HELP را برای هر یک ارائه دهید، در غیر این صورت HELP ناسازگار باعث ایجاد خطا می‌شود.

معمولاً فایل‌های .prom را با یک cronjob ایجاد و به‌روزرسانی می‌کنید. از آنجا که scrape می‌تواند در هر زمان اتفاق بیفتد، مهم است که Node Exporter فایل‌های نیمه‌نوشته را نبیند. برای این منظور، ابتدا باید در یک فایل موقت در همان دایرکتوری بنویسید و سپس فایل کامل را به نام فایل نهایی منتقل کنید.

مثال 7-3 یک cronjob را نشان می‌دهد که به جمع‌آورنده textfile خروجی می‌دهد. این متریک‌ها را در یک فایل موقت ایجاد می‌کند و آنها را به نام فایل نهایی تغییر نام می‌دهد. این یک مثال ساده است که از دستورات کوتاه استفاده می‌کند، اما در بیشتر موارد استفاده واقعی، می‌خواهید یک اسکریپت ایجاد کنید تا موارد خوانا بمانند.

مثال 7-3. /etc/crontab که تعداد خطوط در `/etc/shadow` را به عنوان متریک shadow_entries با استفاده از جمع‌آورنده textfile ارائه می‌دهد

```
TEXTFILE=/path/to/textfile/directory
```

Bash

```
# This must all be on one line
*/5 * * * * root (echo -n 'shadow_entries '; grep -c . /etc/shadow) > $TEXTFILE/shadow.prom.$$ && mv $TEXTFILE/shadow.prom.$$ $TEXTFILE/shadow.prom
```

تعدادی اسکریپت نمونه برای استفاده با جمع‌آورنده textfile در مخزن GitHub اسکریپت‌های نمونه جمع‌آورنده textfile موجود است.

### مهر زمانی (Timestamps)

در حالی که قالب ارائه از مهر زمانی پشتیبانی می‌کند، نمی‌توانید از آنها با جمع‌آورنده textfile استفاده کنید. این به این دلیل است که از نظر معنایی منطقی نیست، زیرا متریک‌های شما با همان مهر زمانی سایر متریک‌های scrape ظاهر نمی‌شوند. در عوض، mtime16 فایل در متریک `node_textfile_mtime_seconds` برای شما در دسترس است. می‌توانید از این برای هشدار دادن در مورد کار نکردن cronjobهای خود استفاده کنید، زیرا اگر این مقدار از مدت زمان طولانی قبل باشد، می‌تواند نشان‌دهنده یک مشکل باشد:



```Bash
# HELP node_textfile_mtime_seconds Unixtime mtime of textfiles successfully read.
# TYPE node_textfile_mtime_seconds gauge
node_textfile_mtime_seconds{file="example.prom"} 1.516205651e+09
```

اکنون که Node Exporter را در حال اجرا دارید، بیایید نگاهی بیندازیم که چگونه می‌توانید به پرومتئوس در مورد تمام ماشین‌هایی که آن را روی آنها اجرا می‌کنید بگویید.