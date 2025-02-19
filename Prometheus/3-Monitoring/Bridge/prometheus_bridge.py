from flask import Flask
import requests

app = Flask(__name__)

# URL of the app's JSON metrics endpoint
TARGET_URL = "http://localhost:5000/stats"

@app.route('/metrics')
def metrics():
    try:
        # Fetch JSON metrics from the app
        response = requests.get(TARGET_URL)
        data = response.json()

        # Convert JSON to Prometheus format
        prom_metrics = []
        prom_metrics.append("# HELP app_requests_total Total HTTP requests.")
        prom_metrics.append("# TYPE app_requests_total counter")
        prom_metrics.append(f"app_requests_total {data['requests_total']}")

        prom_metrics.append("# HELP app_errors_total Total HTTP errors.")
        prom_metrics.append("# TYPE app_errors_total counter")
        prom_metrics.append(f"app_errors_total {data['errors_total']}")

        prom_metrics.append("# HELP app_latency_seconds HTTP request latency in seconds.")
        prom_metrics.append("# TYPE app_latency_seconds gauge")
        prom_metrics.append(f"app_latency_seconds {data['latency_seconds']}")

        return "\n".join(prom_metrics), 200, {'Content-Type': 'text/plain'}

    except Exception as e:
        return f"Error: {str(e)}", 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
