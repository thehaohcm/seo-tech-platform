import json
import logging
import redis
import time

logger = logging.getLogger(__name__)


class QueueListener:
    """
    Redis queue listener for processing jobs
    """
    
    def __init__(self, redis_url: str):
        host, port = redis_url.split(':')
        self.redis_client = redis.Redis(
            host=host,
            port=int(port),
            db=0,
            decode_responses=True
        )
        
        # Test connection
        try:
            self.redis_client.ping()
            logger.info("Connected to Redis successfully")
        except Exception as e:
            logger.error(f"Failed to connect to Redis: {str(e)}")
            raise
    
    def listen(self, queue_name: str, handler):
        """
        Listen for jobs in the queue and process them
        
        Args:
            queue_name: Name of the Redis queue to listen to
            handler: Function to handle each job
        """
        logger.info(f"Listening to queue: {queue_name}")
        
        while True:
            try:
                # Block and wait for job (timeout after 5 seconds)
                result = self.redis_client.blpop(queue_name, timeout=5)
                
                if result is None:
                    # No jobs, continue waiting
                    continue
                
                # Extract job data
                _, job_data = result
                
                # Parse JSON
                job = json.loads(job_data)
                
                # Process the job
                logger.info(f"Processing job: {job.get('url', 'Unknown URL')}")
                handler(job)
                
            except json.JSONDecodeError as e:
                logger.error(f"Invalid JSON in queue: {str(e)}")
            except Exception as e:
                logger.error(f"Error processing job: {str(e)}", exc_info=True)
                # Sleep a bit before retrying to avoid tight loop on persistent errors
                time.sleep(1)
    
    def close(self):
        """
        Close Redis connection
        """
        self.redis_client.close()
