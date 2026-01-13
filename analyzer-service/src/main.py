import os
import json
import logging
from dotenv import load_dotenv

from audits.lighthouse_runner import LighthouseRunner
from audits.accessibility_checker import AccessibilityChecker
from ai_agent.suggestion_generator import SuggestionGenerator
from models.database import Database
from queue_listener import QueueListener

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
        db.save_page_audit({
            'run_id': run_id,
            'url': url,
            'lcp_score': lighthouse_results.get('lcp'),
            'fid_score': lighthouse_results.get('fid'),
            'cls_score': lighthouse_results.get('cls'),
            'seo_issues': json.dumps(lighthouse_results.get('seo_issues', [])),
            'accessibility_issues': json.dumps(accessibility_results.get('violations', [])),
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
    
    # Start listening for analysis jobs
    logger.info("Listening for analysis jobs...")
    queue_listener.listen('analysis_queue', process_page_audit)


if __name__ == "__main__":
    main()
