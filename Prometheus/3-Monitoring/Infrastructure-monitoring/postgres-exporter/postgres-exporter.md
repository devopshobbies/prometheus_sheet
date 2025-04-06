
# Reference

concept:
https://grafana.com/oss/prometheus/exporters/postgres-exporter/
https://github.com/prometheus-community/postgres_exporter

dashboard:
https://grafana.com/grafana/dashboards/9628-postgresql-database/
https://grafana.com/grafana/dashboards/14114-postgres-overview/
https://grafana.com/grafana/dashboards/455-postgres-overview/
https://grafana.com/grafana/dashboards/6742-postgresql-statistics/

data source:
https://grafana.com/docs/grafana/latest/datasources/postgres/
https://grafana.com/docs/grafana/latest/datasources/postgres/configure/



# **PostgreSQL Exporter برای Prometheus**  

ءء **Postgres Exporter** یک ابزار متن‌باز است که داده‌های متریک مربوط به پایگاه داده **PostgreSQL** را جمع‌آوری کرده و در قالب قابل خواندن توسط **Prometheus** ارائه می‌دهد. این ابزار برای مانیتورینگ سلامت، کارایی و استفاده از منابع پایگاه داده PostgreSQL بسیار مفید است.  

---

## **1. کاربردهای Postgres Exporter**  
- **مانیتورینگ کوئری‌های کند** (pg_stat_statements)  
- **ردیابی اتصالات فعال و idle**  
- **بررسی لاگ‌های ریپلیکیشن** (در سیستم‌های Primary/Replica)  
- **تحلیل استفاده از فضای دیسک و I/O**  
- **مانیتورینگ locks و deadlocks**  
- **بررسی وضعیت autovacuum**  

---

## **2. نحوه نصب و راه‌اندازی**  

### **الف. نصب Postgres Exporter**  
#### **روش ۱: دانلود از GitHub (پیش‌ساخته)**  
```bash
wget https://github.com/prometheus-community/postgres_exporter/releases/download/v0.17.1/postgres_exporter-0.17.1.linux-amd64.tar.gz
tar xvf postgres_exporter-*.tar.gz
cd postgres_exporter-*/
./postgres_exporter
```

#### **روش ۲: نصب با Docker**  
```bash
# Start an example database
docker run --net=host -it --rm -e POSTGRES_PASSWORD=password postgres
# Connect to it
docker run \
  --net=host \
  -e DATA_SOURCE_URI="localhost:5432/postgres?sslmode=disable" \
  -e DATA_SOURCE_USER=postgres \
  -e DATA_SOURCE_PASS=password \
  quay.io/prometheuscommunity/postgres-exporter
```

---

### **ب. تنظیمات اتصال به PostgreSQL**  
برای اتصال **Exporter** به **PostgreSQL**، باید یک کاربر مخصوص مانیتورینگ ایجاد کنید:  

```sql
CREATE USER postgres_exporter WITH PASSWORD 'secure_password';
ALTER USER postgres_exporter SET SEARCH_PATH TO postgres_exporter,pg_catalog;

-- دسترسی‌های لازم برای خواندن متریک‌ها
GRANT pg_monitor TO postgres_exporter;
GRANT CONNECT ON DATABASE postgres TO postgres_exporter;
```

---

### **پ. تنظیمات Prometheus برای جمع‌آوری داده‌ها**  
در فایل `prometheus.yml`، یک **scrape_config** جدید اضافه کنید:  

```yaml
scrape_configs:
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres_exporter:9187']  # آدرس Postgres Exporter
    metrics_path: '/metrics'
```

---

## **3. مهم‌ترین متریک‌های قابل مانیتورینگ**  

| **متریک** | **توضیح** |
|-----------|-----------|
| `pg_up` | وضعیت اتصال به PostgreSQL (1 = فعال، 0 = قطع) |
| `pg_stat_activity_count` | تعداد اتصالات فعال |
| `pg_stat_database_xact_commit` | تعداد تراکنش‌های commit شده |
| `pg_stat_database_xact_rollback` | تعداد تراکنش‌های rollback شده |
| `pg_stat_user_tables_n_live_tup` | تعداد رکوردهای زنده در هر جدول |
| `pg_stat_user_indexes_idx_scan` | تعداد اسکن‌های ایندکس |
| `pg_replication_lag` | تاخیر ریپلیکیشن (بر حسب بایت) |
| `pg_stat_bgwriter_checkpoints_timed` | تعداد checkpointهای زمان‌بندی شده |

---

## **4. تنظیم Dashboards در Grafana**  
برای نمایش بهتر داده‌ها، می‌توانید از **داشبوردهای آماده Grafana** استفاده کنید:  
- **PostgreSQL Overview** (ID: **9628**)  
- **PostgreSQL Database** (ID: **6742**)  

**نحوه اضافه کردن:**  
1. در Grafana به **Dashboards > Import** بروید.  
2. شناسه (ID) داشبورد را وارد کنید.  
3. منبع داده (Prometheus) را انتخاب کنید.  

---

## **5. خطاهای رایج و عیب‌یابی**  
- **`ERROR: permission denied for relation pg_stat_activity`**  
  → کاربر Exporter دسترسی کافی ندارد (`GRANT pg_monitor TO postgres_exporter`).  
- **`connection refused`**  
  ءء→ PostgreSQL به درستی listen نمی‌کند یا فایروال مسدود کرده است.  
- **`no metrics are being collected`**  
  → بررسی کنید که `DATA_SOURCE_NAME` به درستی تنظیم شده باشد.  

---

### **جمع‌بندی**  
✅ با **Postgres Exporter** ابزاری ضروری برای مانیتورینگ PostgreSQL در Prometheus است.  
✅ با **Grafana** می‌توانید داشبوردهای حرفه‌ای برای تحلیل عملکرد بسازید.  
✅ تنظیم صحیح دسترسی کاربر Exporter از خطاهای رایج جلوگیری می‌کند.  

اگر نیاز به تنظیمات پیشرفته‌تر دارید، می‌توانید **Custom Queries** در فایل `postgres_exporter.yml` تعریف کنید. 

### SQL for create table

```sql
-- Create a sample table with various data types

CREATE TABLE sample_data (

id SERIAL PRIMARY KEY,

customer_name VARCHAR(100) NOT NULL,

email VARCHAR(100) UNIQUE,

age INTEGER CHECK (age BETWEEN 18 AND 120),

registration_date DATE NOT NULL DEFAULT CURRENT_DATE,

membership_level VARCHAR(20) CHECK (membership_level IN ('Basic', 'Premium', 'VIP')),

account_balance DECIMAL(10, 2),

is_active BOOLEAN DEFAULT TRUE,

last_login TIMESTAMP,

random_code CHAR(5)

);

  

-- Create an index for better performance on common queries

-- CREATE INDEX idx_customer_name ON sample_data(customer_name);

-- CREATE INDEX idx_registration_date ON sample_data(registration_date);
```

## SQL for fill Table

```sql
-- Insert 1000 random records into the table

INSERT INTO sample_data (

customer_name,

email,

age,

registration_date,

membership_level,

account_balance,

is_active,

last_login,

random_code

)

SELECT

'Customer_' || i,

'user_' || i || '@' ||

CASE (i % 4)

WHEN 0 THEN 'gmail.com'

WHEN 1 THEN 'yahoo.com'

WHEN 2 THEN 'outlook.com'

ELSE 'example.com'

END,

(random() * 60 + 18)::integer, -- Age between 18 and 78

CURRENT_DATE - (random() * 365 * 3)::integer, -- Registration in last 3 years

CASE (i % 3)

WHEN 0 THEN 'Basic'

WHEN 1 THEN 'Premium'

ELSE 'VIP'

END,

(random() * 10000)::numeric(10,2), -- Balance between 0 and 10,000

random() > 0.2, -- 80% chance of being active

CASE

WHEN random() > 0.3 THEN CURRENT_TIMESTAMP - (random() * INTERVAL '30 days')

ELSE NULL -- 30% chance of never logging in

END,

substr(md5(random()::text), 0, 5) -- Random 5-character code

FROM generate_series(1, 1000) AS i;
```

## clear table

```sql
TRUNCATE TABLE sample_data;
```