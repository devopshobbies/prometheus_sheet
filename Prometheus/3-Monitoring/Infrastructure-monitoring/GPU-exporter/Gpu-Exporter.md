
# **GPU Exporter برای Prometheus**  

‏**GPU Exporter** یک ابزار مانیتورینگ است که اطلاعات مربوط به کارت‌های گرافیک (**GPU**) را جمع‌آوری کرده و در قالب متریک‌های Prometheus ارائه می‌دهد. این ابزار برای سیستم‌هایی که از GPU برای پردازش‌های سنگین مانند **یادگیری ماشین (ML)، رندرینگ، ماینینگ یا محاسبات HPC** استفاده می‌کنند، بسیار مفید است.  

---

## **۱. کاربردهای GPU Exporter**  
✅ **مانیتورینگ استفاده از حافظه GPU (VRAM)**  
✅ **ردیابی میزان بار پردازشی (GPU Utilization)**  
✅ **بررسی دمای GPU و سرعت فن‌ها**  
✅ **تحلیل مصرف انرژی (Power Draw)**  
✅ **شناسایی bottleneckها در پردازش‌های موازی**  

---

## **۲. نصب و راه‌اندازی GPU Exporter**  

### **الف) روش‌های نصب**  

#### **۱. استفاده از NVIDIA DCGM Exporter (رسمی توسط انویدیا)**  
```bash
# دانلود و اجرای DCGM Exporter با Docker
docker run -d \
  --gpus all \
  --name nvidia-dcgm-exporter \
  -p 9400:9400 \
  nvcr.io/nvidia/k8s/dcgm-exporter:3.3.4-3.1.5-ubuntu22.04
```

#### **۲. استفاده از `nvidia_gpu_exporter` (راه‌حل ساده‌تر)**  
```bash
# دانلود از GitHub
wget https://github.com/utkuozdemir/nvidia_gpu_exporter/releases/download/v1.1.0/nvidia_gpu_exporter_1.1.0_linux_x86_64.tar.gz
tar -xvf nvidia_gpu_exporter_*.tar.gz
./nvidia_gpu_exporter
```

#### **۳. روش دستی با `nvidia-smi` و `node_exporter`**  
اگر Exporter اختصاصی ندارید، می‌توانید از `nvidia-smi` و یک اسکریپت `textfile_collector` در `node_exporter` استفاده کنید.  

---

### **ب) تنظیم Prometheus برای جمع‌آوری داده‌ها**  
در فایل `prometheus.yml`، یک **scrape_config** جدید اضافه کنید:  

```yaml
scrape_configs:
  - job_name: 'nvidia_gpu'
    static_configs:
      - targets: ['localhost:9400']  # پورت پیش‌فرض DCGM Exporter
    metrics_path: '/metrics'
```

---

## **۳. مهم‌ترین متریک‌های GPU**  

| **متریک** | **توضیح** |  
|-----------|-----------|  
| `dcgm_gpu_utilization` | درصد استفاده از GPU (0-100) |  
| `dcgm_mem_utilization` | درصد استفاده از VRAM |  
| `dcgm_gpu_temp` | دمای GPU (بر حسب °C) |  
| `dcgm_power_usage` | مصرف برق (وات) |  
| `dcgm_fb_used` | حافظه GPU مصرف شده (مگابایت) |  
| `dcgm_fb_free` | حافظه GPU آزاد (مگابایت) |  

---

## **۴. تنظیم Dashboards در Grafana**  
برای نمایش داده‌های GPU می‌توانید از **داشبوردهای آماده Grafana** استفاده کنید:  
‏- **NVIDIA DCGM Exporter Dashboard** (ID: **12239**)  
‏- **GPU Monitoring via node-exporter** (ID: **14574**)  

**نحوه اضافه کردن:**  
۱. در Grafana به **Dashboards > Import** بروید.  
۲. شناسه (ID) داشبورد را وارد کنید.  
۳. منبع داده (Prometheus) را انتخاب کنید.  

---

## **۵. عیب‌یابی و خطاهای رایج**  

### ** مشکل: `No GPUs detected`**  
**راه‌حل:**  
- مطمئن شوید **درایورهای NVIDIA** نصب هستند (`nvidia-smi` باید کار کند).  
- اگر از **Docker** استفاده می‌کنید، `--gpus all` را اضافه کنید.  

### ** مشکل: `Permission denied` برای `/dev/nvidia*`**  
**راه‌حل:**  
```bash
sudo chmod a+rw /dev/nvidia*
```
یا کاربر exporter را به گروه `video` اضافه کنید:  
```bash
sudo usermod -aG video prometheus
```

###  مشکل: `DCGM doesn't support this GPU architecture`**  
**راه‌حل:**  
نسخه جدیدتر **NVIDIA DCGM** را نصب کنید یا از `nvidia_gpu_exporter` جایگزین استفاده کنید.  

---

## **۶. جمع‌بندی**  
‏🔹 **GPU Exporter** برای مانیتورینگ کارت‌های گرافیک در **Prometheus** ضروری است.  
‏🔹 **DCGM Exporter** (راه‌حل انویدیا) بهترین گزینه برای سیستم‌های حرفه‌ای است.  
‏🔹 **Grafana** برای نمایش داده‌ها به صورت گرافیکی مناسب است.  
🔹 اگر مشکل دسترسی داشتید، بررسی کنید که **درایورهای NVIDIA** و **Docker (در صورت نیاز)** به درستی تنظیم شده باشند.  

اگر نیاز به تنظیمات پیشرفته‌تر دارید، می‌توانید **Custom Queries** در Prometheus تعریف کنید! 
## plan 1
https://github.com/NVIDIA/dcgm-exporter
https://stackoverflow.com/a/77269071/6274589

```sh
docker run -d --gpus all --cap-add SYS_ADMIN --rm -p 9400:9400 nvcr.io/nvidia/k8s/dcgm-exporter:4.1.1-4.0.4-ubuntu22.04
curl localhost:9400/metrics

```


### plan 2
https://github.com/utkuozdemir/nvidia_gpu_exporter?tab=readme-ov-file
https://github.com/utkuozdemir/nvidia_gpu_exporter/blob/master/INSTALL.md

```sh
$ docker run -d \
--name nvidia_smi_exporter \
--restart unless-stopped \
--device /dev/nvidiactl:/dev/nvidiactl \
--device /dev/nvidia0:/dev/nvidia0 \
-v /usr/lib/x86_64-linux-gnu/libnvidia-ml.so:/usr/lib/x86_64-linux-gnu/libnvidia-ml.so \
-v /usr/lib/x86_64-linux-gnu/libnvidia-ml.so.1:/usr/lib/x86_64-linux-gnu/libnvidia-ml.so.1 \
-v /usr/bin/nvidia-smi:/usr/bin/nvidia-smi \
-p 9835:9835 \
utkuozdemir/nvidia_gpu_exporter:1.1.0
```


## load test
https://www.geeks3d.com/20140304/gputest-0-7-0-opengl-benchmark-win-linux-osx-new-fp64-opengl-4-test-and-online-gpu-database/