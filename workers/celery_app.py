"""
Celery Application Configuration
"""
from celery import Celery
from celery.schedules import crontab
import os

# Redis URL from environment or default
REDIS_URL = os.environ.get('REDIS_URL', 'redis://localhost:6379/0')
RESULT_BACKEND = os.environ.get('RESULT_BACKEND', 'redis://localhost:6379/1')

# Initialize Celery app
app = Celery(
    'golib_api_tasks',
    broker=REDIS_URL,
    backend=RESULT_BACKEND,
    include=['workers.tasks']
)

# Celery configuration
app.conf.update(
    task_serializer='json',
    accept_content=['json'],
    result_serializer='json',
    timezone='UTC',
    enable_utc=True,
    task_track_started=True,
    task_time_limit=30 * 60,  # 30 minutes
    task_soft_time_limit=25 * 60,  # 25 minutes
    worker_prefetch_multiplier=4,
    worker_max_tasks_per_child=1000,
    
    # Result backend configuration
    result_expires=3600,  # Results expire after 1 hour
    result_extended=True,
    
    # Task routing
    task_routes={
        'tasks.send_email': {'queue': 'notifications'},
        'tasks.process_transaction': {'queue': 'transactions'},
        'tasks.generate_statement': {'queue': 'reports'},
    },
    
    # Beat schedule for periodic tasks
    beat_schedule={
        'cleanup-expired-sessions': {
            'task': 'tasks.cleanup_expired_sessions',
            'schedule': crontab(minute='0', hour='*'),  # Every hour
        },
        'generate-daily-report': {
            'task': 'tasks.generate_daily_report',
            'schedule': crontab(minute='0', hour='0'),  # Daily at midnight
        },
    },
)

# Auto-discover tasks
app.autodiscover_tasks(['workers'])

if __name__ == '__main__':
    app.start()
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
