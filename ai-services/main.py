#!/usr/bin/env python
"""
Hello Pulse AI Microservice
A modular, scalable AI service for Hello Pulse platform
"""

import time
from contextlib import asynccontextmanager
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from api.routes import generate, rag, search, agents
from core.config import settings
from core.logging import logger, setup_logging

# Set up logging
setup_logging()

# Global startup resources
@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup: Initialize connections and resources
    logger.info("Starting Hello Pulse AI Microservice")
    
    # Add any global resources here, like DB connections
    # that should be available for the lifetime of the app
    
    yield
    
    # Shutdown: Clean up resources
    logger.info("Shutting down Hello Pulse AI Microservice")

# Create FastAPI app
app = FastAPI(
    title=settings.PROJECT_NAME,
    description=settings.PROJECT_DESCRIPTION,
    version=settings.PROJECT_VERSION,
    lifespan=lifespan,
)

# Set up CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Request timing middleware
@app.middleware("http")
async def add_process_time_header(request: Request, call_next):
    start_time = time.time()
    try:
        response = await call_next(request)
        process_time = time.time() - start_time
        response.headers["X-Process-Time"] = str(process_time)
        
        # Log request completion
        logger.info(
            f"Request completed",
            extra={
                "path": request.url.path,
                "method": request.method,
                "process_time": process_time,
                "status_code": response.status_code
            }
        )
        return response
    except Exception as e:
        logger.exception(f"Request failed: {str(e)}")
        process_time = time.time() - start_time
        return JSONResponse(
            status_code=500,
            content={"detail": "Internal server error"}
        )

# Include routers
app.include_router(generate.router, prefix="/api")
app.include_router(rag.router, prefix="/api")
app.include_router(search.router, prefix="/api")
app.include_router(agents.router, prefix="/api")

# Health check endpoint
@app.get("/health")
async def health_check():
    return {"status": "healthy"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG,
        log_level="info",
    )