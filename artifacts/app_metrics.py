from flask import Flask, Response
from prometheus_client import (
    Counter,
    Histogram,
    Gauge,
    generate_latest,
    CONTENT_TYPE_LATEST
)

import time
import random

app = Flask(__name__)

http_requests_total = Counter(
    "http_requests_total",
    "Total HTTP requests",
    ["method", "endpoint", "status"]
)

http_request_duration_seconds = Histogram(
    "http_request_duration_seconds",
    "HTTP request duration",
    ["endpoint"]
)

active_requests = Gauge(
    "active_requests",
    "Current active requests"
)

users_requests_total = Counter(
    "users_requests_total",
    "Business metric: users endpoint calls"
)

@app.route("/metrics")
def metrics():
    return Response(
        generate_latest(),
        mimetype=CONTENT_TYPE_LATEST
    )

@app.route("/api/users")
def users():

    start = time.time()

    active_requests.inc()

    try:
        time.sleep(random.uniform(0.05, 0.15))

        users_requests_total.inc()

        http_requests_total.labels(
            "GET",
            "/api/users",
            "200"
        ).inc()

        return {"users": ["user1", "user2"]}

    finally:

        http_request_duration_seconds.labels(
            "/api/users"
        ).observe(time.time() - start)

        active_requests.dec()


@app.route("/api/fail")
def fail():

    http_requests_total.labels(
        "GET",
        "/api/fail",
        "500"
    ).inc()

    return {"error": "simulated"}, 500


if __name__ == "__main__":
    app.run(
        host="0.0.0.0",
        port=8080
    )