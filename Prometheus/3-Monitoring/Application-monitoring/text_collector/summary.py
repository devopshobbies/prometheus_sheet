from prometheus_client import Summary
from prometheus_client import CollectorRegistry, Gauge, write_to_textfile

s = Summary('request_latency_seconds', 'Description of summary')

s.observe(4.7)    # Observe 4.7 (seconds in this case)

write_to_textfile('./TextCollect/summary.prom', s)
