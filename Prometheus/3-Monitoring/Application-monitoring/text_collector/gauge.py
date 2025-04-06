from prometheus_client import Gauge
from prometheus_client import CollectorRegistry, Gauge, write_to_textfile


g = Gauge('my_inprogress_requests', 'Description of gauge')
g.inc()      # Increment by 1
g.dec(10)    # Decrement by given value
g.set(4.2)   # Set to a given value
write_to_textfile('./TextCollect/gauge.prom', g)

