"""
Logging configuration for the Hello Pulse AI Microservice
"""

import logging
import sys
import json
from datetime import datetime
from typing import Any, Dict, Optional

from core.config import settings

class JsonFormatter(logging.Formatter):
    """
    Formatter that outputs JSON strings after parsing the log record.
    """
    def format(self, record: logging.LogRecord) -> str:
        """Format log record as JSON"""
        log_record: Dict[str, Any] = {
            "timestamp": datetime.utcnow().isoformat(),
            "level": record.levelname,
            "message": record.getMessage(),
            "module": record.module,
            "function": record.funcName,
        }
        
        # Add exception info if available
        if record.exc_info:
            log_record["exception"] = self.formatException(record.exc_info)
        
        # Add extra fields from the record
        if hasattr(record, "organization_id"):
            log_record["organization_id"] = record.organization_id
            
        if hasattr(record, "user_id"):
            log_record["user_id"] = record.user_id
            
        if hasattr(record, 'request_id'):
            log_record['request_id'] = record.request_id
            
        # Add any extra attributes from the record
        for key, value in record.__dict__.items():
            if key not in {
                "args", "asctime", "created", "exc_info", "exc_text", "filename",
                "funcName", "id", "levelname", "levelno", "lineno", "module",
                "msecs", "message", "msg", "name", "pathname", "process",
                "processName", "relativeCreated", "stack_info", "thread", "threadName"
            }:
                log_record[key] = value
                
        return json.dumps(log_record)

def setup_logging() -> None:
    """Configure logging for the application"""
    # Get log level from settings
    log_level = getattr(logging, settings.LOG_LEVEL.upper(), logging.INFO)
    
    # Configure root logger
    root_logger = logging.getLogger()
    root_logger.setLevel(log_level)
    
    # Remove existing handlers
    for handler in root_logger.handlers[:]:
        root_logger.removeHandler(handler)
    
    # Create console handler with JSON formatter
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(JsonFormatter())
    root_logger.addHandler(console_handler)
    
    # Suppress unwanted logs from third-party libraries
    logging.getLogger("uvicorn.access").setLevel(logging.WARNING)
    logging.getLogger("uvicorn.error").setLevel(logging.ERROR)
    
    # Create logger for this application
    logger = logging.getLogger("hello_pulse_ai")
    logger.setLevel(log_level)

# Create a global logger instance
logger = logging.getLogger("hello_pulse_ai")

class LoggingContext:
    """
    Context manager for adding context to logs
    """
    def __init__(
        self, 
        organization_id: Optional[str] = None,
        user_id: Optional[str] = None,
        request_id: Optional[str] = None,
        **extra: Any
    ):
        self.organization_id = organization_id
        self.user_id = user_id
        self.request_id = request_id
        self.extra = extra
        self.old_factory = None
    
    def __enter__(self):
        """Add context to logging"""
        old_factory = logging.getLogRecordFactory()
        
        def record_factory(*args, **kwargs):
            record = old_factory(*args, **kwargs)
            if self.organization_id:
                record.organization_id = self.organization_id
            if self.user_id:
                record.user_id = self.user_id
            if self.request_id:
                record.request_id = self.request_id
            for key, value in self.extra.items():
                setattr(record, key, value)
            return record
        
        self.old_factory = old_factory
        logging.setLogRecordFactory(record_factory)
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Reset logging context"""
        if self.old_factory:
            logging.setLogRecordFactory(self.old_factory)