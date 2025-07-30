#!/usr/bin/env python3

import requests
import time
import threading
import argparse
from collections import Counter
import sys

# Global tracking variables
response_codes = Counter()
lock = threading.Lock()
completed_requests = 0
total_requests = 0
start_time = None

def send_request():
    """Send a single request and record the response code"""
    global completed_requests

    url = 'http://payment-options-service.i.s-go-gp-eks-01-a6ff5941.gopay.sh/payment-options-service/payment-options'
    params = {'intent': 'GOPAY_WALLET_INTENT'}
    headers = {
        'Content-Type': 'application/json',
        'User-ID': '22102721471112340001250790',
        'X-User-Locale': 'en_ID'
    }

    try:
        response = requests.get(url, params=params, headers=headers, timeout=10)
        status_code = response.status_code
    except requests.exceptions.RequestException as e:
        status_code = 'Error'
        print(f"\nRequest error: {e}")

    with lock:
        response_codes[status_code] += 1
        completed_requests += 1

        # Show progress
        if completed_requests % 10 == 0 or completed_requests == total_requests:
            elapsed = time.time() - start_time
            rps = completed_requests / elapsed if elapsed > 0 else 0
            print(f"\rCompleted: {completed_requests}/{total_requests} ({completed_requests/total_requests*100:.1f}%) - {rps:.1f} req/sec", end="")
            sys.stdout.flush()

def run_load_test(num_requests, requests_per_second):
    """Execute the load test with rate limiting"""
    global total_requests, start_time
    total_requests = num_requests
    start_time = time.time()

    # Calculate time between requests to maintain rate limit
    delay = 1.0 / requests_per_second if requests_per_second > 0 else 0
    threads = []

    print(f"Starting load test with {num_requests} requests at {requests_per_second} req/sec")

    for _ in range(num_requests):
        thread = threading.Thread(target=send_request)
        thread.daemon = True
        thread.start()
        threads.append(thread)

        # Sleep to maintain the rate limit
        if delay > 0:
            time.sleep(delay)

    # Wait for all requests to complete
    for thread in threads:
        thread.join()

    duration = time.time() - start_time
    return duration

def print_report(duration):
    """Print test results and statistics"""
    print("\n\n===== Load Test Results =====")
    print(f"Total requests: {total_requests}")
    print(f"Test duration: {duration:.2f} seconds")
    print(f"Average throughput: {total_requests/duration:.2f} requests/second")

    # Calculate success rate
    success_count = response_codes.get(200, 0)
    success_rate = (success_count / total_requests * 100) if total_requests > 0 else 0

    print("\nResponse Code Summary:")
    # Convert keys to strings for sorting to handle both int and str types
    for code, count in sorted(response_codes.items(), key=lambda x: str(x[0])):
        percentage = (count / total_requests * 100)
        print(f"  HTTP {code}: {count} ({percentage:.2f}%)")

    print(f"\nSuccess rate (HTTP 200): {success_rate:.2f}%")
    print(f"Other responses: {100 - success_rate:.2f}%")

def main():
    parser = argparse.ArgumentParser(description="API Load Testing Tool")
    parser.add_argument("--requests", type=int, default=100,
                        help="Total number of requests to send")
    parser.add_argument("--rate", type=float, default=10,
                        help="Request rate (requests per second)")
    args = parser.parse_args()

    try:
        duration = run_load_test(args.requests, args.rate)
        print_report(duration)
    except KeyboardInterrupt:
        print("\nTest interrupted by user")
        if start_time:
            print_report(time.time() - start_time)
        sys.exit(1)

if __name__ == "__main__":
    main()