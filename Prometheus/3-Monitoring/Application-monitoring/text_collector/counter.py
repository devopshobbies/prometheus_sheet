from prometheus_client import Counter
from prometheus_client import CollectorRegistry, Gauge, write_to_textfile

c = Counter('my_failures', 'Description of counter')

c.inc()     # Increment by 1

c.inc(1.6)  # Increment by given value
write_to_textfile('./TextCollect/counter.prom', c)
