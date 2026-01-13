import logging
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from axe_selenium_python import Axe

logger = logging.getLogger(__name__)


class AccessibilityChecker:
    """
    Runs accessibility checks using Axe-core
    """
    
    def __init__(self):
        self.driver = None
    
    def check(self, url: str) -> dict:
        """
        Run accessibility check on a URL
        
        Args:
            url: The URL to check
            
        Returns:
            Dictionary with accessibility results
        """
        try:
            logger.info(f"Running accessibility check for: {url}")
            
            # Setup headless Chrome
            chrome_options = Options()
            chrome_options.add_argument("--headless")
            chrome_options.add_argument("--no-sandbox")
            chrome_options.add_argument("--disable-dev-shm-usage")
            
            self.driver = webdriver.Chrome(options=chrome_options)
            self.driver.get(url)
            
            # Run Axe
            axe = Axe(self.driver)
            axe.inject()
            results = axe.run()
            
            # Process results
            processed_results = self._process_results(results)
            
            logger.info(f"Accessibility check completed for: {url}")
            
            return processed_results
            
        except Exception as e:
            logger.error(f"Error running accessibility check: {str(e)}")
            return {'violations': [], 'passes': [], 'incomplete': []}
        
        finally:
            if self.driver:
                self.driver.quit()
    
    def _process_results(self, results: dict) -> dict:
        """
        Process Axe results into a simplified format
        """
        violations = []
        
        for violation in results.get('violations', []):
            violations.append({
                'id': violation.get('id'),
                'impact': violation.get('impact'),
                'description': violation.get('description'),
                'help': violation.get('help'),
                'helpUrl': violation.get('helpUrl'),
                'nodes': len(violation.get('nodes', []))
            })
        
        return {
            'violations': violations,
            'passes': len(results.get('passes', [])),
            'incomplete': len(results.get('incomplete', [])),
            'total_violations': len(violations)
        }
