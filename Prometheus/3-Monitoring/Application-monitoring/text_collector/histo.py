from prometheus_client import Histogram
from prometheus_client import CollectorRegistry, Gauge, write_to_textfile

h = Histogram('request_latency_seconds', 'Description of histogram')
h.observe(4.7)    # Observe 4.7 (seconds in this case)

write_to_textfile('./TextCollect/histo.prom', h)
