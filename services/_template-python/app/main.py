#!/usr/bin/env python3
"""
FamGo Python FastAPI Service Template

This is the standard FastAPI service template for all ML/analytics services:
- demand-prediction-service
- eta-prediction-service
- surge-prediction-service
- fraud-detection-ml
- pooling-optimization-ml
- analytics-service
"""

import os
import logging
import asyncio
from contextlib import asynccontextmanager
from datetime import datetime

from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.responses import JSONResponse
import uvicorn
import aioredis
import asyncpg
from kafka import KafkaProducer, KafkaConsumer
import json

# Configuration
SERVICE_NAME = os.getenv("SERVICE_NAME", "template-service")
PORT = int(os.getenv("PORT", "8080"))
LOG_LEVEL = os.getenv("LOG_LEVEL", "INFO")

# Setup logging
logging.basicConfig(
    level=LOG_LEVEL,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(SERVICE_NAME)

# Dependencies
redis_client = None
db_pool = None
kafka_producer = None
kafka_consumer = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Manage service lifecycle - startup and shutdown"""
    # Startup
    logger.info(f"Starting {SERVICE_NAME}")
    
    global redis_client, db_pool, kafka_producer, kafka_consumer
    
    # Connect to Redis
    redis_client = await aioredis.create_redis_pool(
        os.getenv("REDIS_HOST", "localhost"),
        int(os.getenv("REDIS_PORT", 6379))
    )
    logger.info("Connected to Redis")
    
    # Connect to PostgreSQL
    db_pool = await asyncpg.create_pool(
        user=os.getenv("DB_USER", "postgres"),
        password=os.getenv("DB_PASSWORD", "postgres"),
        database=os.getenv("DB_NAME", "famgo"),
        host=os.getenv("DB_HOST", "localhost"),
        port=int(os.getenv("DB_PORT", 5432)),
        min_size=5,
        max_size=20,
    )
    logger.info("Connected to PostgreSQL")
    
    # Setup Kafka
    kafka_producer = KafkaProducer(
        bootstrap_servers=os.getenv("KAFKA_BROKERS", "localhost:9092").split(","),
        value_serializer=lambda v: json.dumps(v).encode('utf-8')
    )
    logger.info("Connected to Kafka")
    
    # Start Kafka consumer in background
    asyncio.create_task(kafka_consumer_loop())
    
    yield
    
    # Shutdown
    logger.info(f"Shutting down {SERVICE_NAME}")
    
    if redis_client:
        redis_client.close()
        await redis_client.wait_closed()
    
    if db_pool:
        await db_pool.close()
    
    if kafka_producer:
        kafka_producer.close()


# Create FastAPI app
app = FastAPI(
    title=SERVICE_NAME,
    description="FamGo Platform Service",
    version="1.0.0",
    lifespan=lifespan
)


# Health check endpoint
@app.get("/health")
async def health_check():
    """Health check endpoint for Kubernetes"""
    return {
        "status": "ok",
        "service": SERVICE_NAME,
        "timestamp": datetime.utcnow().isoformat()
    }


# Ready check endpoint
@app.get("/ready")
async def ready_check():
    """Readiness check - confirm dependencies are ready"""
    checks = {
        "redis": False,
        "postgres": False,
        "kafka": False,
    }
    
    try:
        if redis_client:
            await redis_client.execute('PING')
            checks["redis"] = True
    except Exception as e:
        logger.error(f"Redis check failed: {e}")
    
    try:
        if db_pool:
            async with db_pool.acquire() as conn:
                await conn.execute('SELECT 1')
            checks["postgres"] = True
    except Exception as e:
        logger.error(f"PostgreSQL check failed: {e}")
    
    try:
        if kafka_producer:
            kafka_producer.send('__health-check', {"test": True})
            checks["kafka"] = True
    except Exception as e:
        logger.error(f"Kafka check failed: {e}")
    
    all_ready = all(checks.values())
    return {
        "ready": all_ready,
        "checks": checks
    } if all_ready else JSONResponse(
        {"ready": False, "checks": checks},
        status_code=503
    )


# Example endpoints - replace with your service logic
@app.get("/api/v1/predict/{model_id}")
async def predict(model_id: str, input_data: dict):
    """
    Example prediction endpoint
    
    Replace this with your ML model inference logic
    """
    try:
        # Load model (cache in Redis)
        model_key = f"model:{model_id}"
        model = await redis_client.get(model_key)
        
        if not model:
            # Load from database or S3
            async with db_pool.acquire() as conn:
                model = await conn.fetchval(
                    "SELECT model_data FROM ml_models WHERE id = $1",
                    model_id
                )
                
                if not model:
                    raise HTTPException(status_code=404, detail="Model not found")
                
                # Cache in Redis for 1 hour
                await redis_client.setex(model_key, 3600, model)
        
        # Run inference
        prediction = await run_inference(model, input_data)
        
        # Publish event
        kafka_producer.send('prediction.completed', {
            "model_id": model_id,
            "prediction": prediction,
            "timestamp": datetime.utcnow().isoformat()
        })
        
        return {
            "model_id": model_id,
            "prediction": prediction,
            "timestamp": datetime.utcnow().isoformat()
        }
        
    except Exception as e:
        logger.error(f"Prediction failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/train/{model_id}")
async def train_model(model_id: str, background_tasks: BackgroundTasks):
    """
    Start model training in background
    
    Replace with your training pipeline
    """
    logger.info(f"Starting training for model {model_id}")
    
    # Run training in background
    background_tasks.add_task(
        train_model_background,
        model_id
    )
    
    return {
        "status": "training_started",
        "model_id": model_id
    }


async def kafka_consumer_loop():
    """Consume events from Kafka"""
    consumer = KafkaConsumer(
        os.getenv("KAFKA_TOPICS", "model.training.started").split(","),
        bootstrap_servers=os.getenv("KAFKA_BROKERS", "localhost:9092").split(","),
        group_id=os.getenv("KAFKA_GROUP_ID", SERVICE_NAME),
        value_deserializer=lambda m: json.loads(m.decode('utf-8'))
    )
    
    for message in consumer:
        try:
            await handle_kafka_event(message.topic, message.value)
        except Exception as e:
            logger.error(f"Error handling Kafka event: {e}")


async def handle_kafka_event(topic: str, event: dict):
    """Handle incoming Kafka events"""
    logger.info(f"Received event from topic {topic}: {event}")
    
    # Route to appropriate handler
    if topic == "model.training.started":
        await on_model_training_started(event)
    elif topic == "ride.created":
        await on_ride_created(event)


async def on_model_training_started(event: dict):
    """Handle model training started event"""
    logger.info(f"Model training started: {event}")
    # Implement your logic


async def on_ride_created(event: dict):
    """Handle ride created event"""
    logger.info(f"Ride created: {event}")
    # Implement your prediction logic


async def run_inference(model, input_data):
    """Run ML model inference"""
    # Implement your inference logic
    return {"score": 0.95}


async def train_model_background(model_id: str):
    """Background task for model training"""
    try:
        logger.info(f"Training model {model_id} in background")
        # Implement training logic
        
        # Save model
        async with db_pool.acquire() as conn:
            await conn.execute(
                "UPDATE ml_models SET model_data = $1, trained_at = NOW() WHERE id = $2",
                None,  # model_data serialized
                model_id
            )
        
        # Publish event
        kafka_producer.send('model.trained', {
            "model_id": model_id,
            "timestamp": datetime.utcnow().isoformat()
        })
        
        logger.info(f"Model {model_id} training completed")
        
    except Exception as e:
        logger.error(f"Background training failed: {e}")
        kafka_producer.send('model.training.failed', {
            "model_id": model_id,
            "error": str(e),
            "timestamp": datetime.utcnow().isoformat()
        })


if __name__ == "__main__":
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=PORT,
        log_level=LOG_LEVEL.lower()
    )
