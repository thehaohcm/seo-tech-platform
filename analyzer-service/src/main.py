import os
import json
import logging
import threading
from dotenv import load_dotenv

from audits.lighthouse_runner import LighthouseRunner
from audits.accessibility_checker import AccessibilityChecker
from ai_agent.suggestion_generator import SuggestionGenerator
from models.database import Database
from queue_listener import QueueListener
from test_runner import AutoTestRunner

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


def process_page_audit(page_data: dict):
    """
    Process a single page audit
    """
    try:
        url = page_data.get('url')
        run_id = page_data.get('run_id')
        
        logger.info(f"Processing audit for URL: {url}")
        
        # Run Lighthouse audit
        lighthouse_runner = LighthouseRunner()
        lighthouse_results = lighthouse_runner.run_audit(url)
        
        # Run accessibility checks
        accessibility_checker = AccessibilityChecker()
        accessibility_results = accessibility_checker.check(url)
        
        # Generate AI suggestions
        suggestion_generator = SuggestionGenerator()
        ai_suggestions = suggestion_generator.generate(
            url=url,
            lighthouse_results=lighthouse_results,
            accessibility_results=accessibility_results,
            page_data=page_data
        )
        
        # Save results to database
        db = Database()
        
        # Check if run_id exists before saving to avoid foreign key violation
        if not db.check_audit_run_exists(run_id):
            logger.warning(f"Audit run {run_id} does not exist. Skipping page audit for {url}")
            return
        
        db.save_page_audit({
            'run_id': run_id,
            'url': url,
            'lcp_score': lighthouse_results.get('lcp'),
            'fid_score': lighthouse_results.get('fid'),
            'cls_score': lighthouse_results.get('cls'),
            'fcp_score': lighthouse_results.get('fcp'),
            'ttfb_score': lighthouse_results.get('ttfb'),
            'status_code': lighthouse_results.get('status_code', 200),
            'load_time_ms': lighthouse_results.get('load_time_ms', 0),
            'title': page_data.get('title', ''),
            'meta_description': page_data.get('meta_description', ''),
            'h1_tags': page_data.get('h1_tags', []),
            'canonical_url': page_data.get('canonical_url', ''),
            'has_robots_meta': bool(page_data.get('has_robots_meta', False)),
            'seo_issues': lighthouse_results.get('seo_issues', []),
            'accessibility_issues': accessibility_results.get('violations', []),
            'performance_issues': lighthouse_results.get('performance_issues', []),
            'ai_suggestions': ai_suggestions
        })
        
        logger.info(f"Successfully processed audit for URL: {url}")
        
    except Exception as e:
        logger.error(f"Error processing page audit: {str(e)}", exc_info=True)
        raise


def main():
    """
    Main entry point for the analyzer service
    """
    logger.info("Starting Analyzer Service...")
    
    # Initialize queue listener
    redis_url = os.getenv('REDIS_URL', 'localhost:6379')
    queue_listener = QueueListener(redis_url)
    
    # Start analysis queue listener in main thread
    logger.info("Listening for analysis jobs...")
    
    # Start test queue listener in separate thread
    def listen_test_queue():
        logger.info("Listening for test jobs...")
        queue_listener.listen('test_queue', process_test_job)
    
    test_thread = threading.Thread(target=listen_test_queue, daemon=True)
    test_thread.start()
    
    # Main thread handles analysis queue
    queue_listener.listen('analysis_queue', process_page_audit)


def process_test_job(test_data: dict):
    """
    Process an automated test job
    """
    try:
        url = test_data.get('url')
        page_id = test_data.get('page_id')
        
        logger.info(f"Running auto tests for URL: {url}, Page ID: {page_id}")
        
        # Run automated tests
        test_runner = AutoTestRunner()
        test_results = test_runner.run_tests(url, page_id)
        
        # Store results in Redis for frontend to retrieve
        import redis
        redis_url = os.getenv('REDIS_URL', 'localhost:6379')
        host, port = redis_url.split(':')
        r = redis.Redis(host=host, port=int(port), decode_responses=False)
        
        # Store as JSON with 1 hour expiry
        result_key = f"test_result:{page_id}"
        r.setex(result_key, 3600, json.dumps(test_results))
        
        logger.info(f"Test completed for {url}: {test_results['status']}")
        logger.info(f"Results: {test_results['passed']}/{test_results['total_tests']} tests passed")
        logger.info(f"Results stored in Redis with key: {result_key}")
        
    except Exception as e:
        logger.error(f"Error processing test job: {str(e)}", exc_info=True)


if __name__ == "__main__":
    main()
