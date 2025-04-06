from flask import Flask
import random
import time

app = Flask(__name__)

# Simulate metrics
request_count = 0
error_count = 0

@app.route('/')
def home():
    global request_count, error_count
    request_count += 1
    if random.random() < 0.2:  # 20% chance of error
        error_count += 1
    return "Hello, World!"

@app.route('/stats')
def stats():
    # Return metrics as JSON
    return {
        "requests_total": request_count,
        "errors_total": error_count,
        "latency_seconds": random.uniform(0.1, 1.0)
    }

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
