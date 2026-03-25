#!/usr/bin/env python3
"""Celery worker entry point for Mini Bank API."""
import os
import sys

# Add project root to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from workers.celery_app import app
from workers.tasks import (
    process_topup,
    process_withdraw,
    process_transfer,
    send_notification,
)

if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(description="Mini Bank Celery Worker")
    parser.add_argument(
        "--loglevel",
        default="info",
        choices=["debug", "info", "warning", "error", "critical"],
        help="Log level",
    )
    parser.add_argument(
        "--concurrency",
        type=int,
        default=4,
        help="Number of concurrent worker processes",
    )
    parser.add_argument(
        "--queues",
        nargs="+",
        default=["topup", "withdraw", "transfer", "notification"],
        help="Queues to consume from",
    )
    args = parser.parse_args()

    app.worker_main(
        [
            "worker",
            f"--loglevel={args.loglevel}",
            f"--concurrency={args.concurrency}",
            "-Q",
            ",".join(args.queues),
            "--autoscale=10,2",
        ]
    )
