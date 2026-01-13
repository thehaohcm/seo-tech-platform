import subprocess
import json
import logging

logger = logging.getLogger(__name__)


class LighthouseRunner:
    """
    Wrapper for running Lighthouse CLI audits
    """
    
    def __init__(self):
        self.lighthouse_cmd = "lighthouse"
    
    def run_audit(self, url: str) -> dict:
        """
        Run Lighthouse audit on a URL
        
        Args:
            url: The URL to audit
            
        Returns:
            Dictionary with Lighthouse results
        """
        try:
            logger.info(f"Running Lighthouse audit for: {url}")
            
            # Run Lighthouse CLI
            cmd = [
                self.lighthouse_cmd,
                url,
                "--output=json",
                "--output-path=stdout",
                "--chrome-flags=--headless",
                "--quiet"
            ]
            
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                timeout=60
            )
            
            if result.returncode != 0:
                logger.error(f"Lighthouse failed: {result.stderr}")
                return {}
            
            # Parse JSON output
            lighthouse_data = json.loads(result.stdout)
            
            # Extract key metrics
            metrics = self._extract_metrics(lighthouse_data)
            
            logger.info(f"Lighthouse audit completed for: {url}")
            return metrics
            
        except subprocess.TimeoutExpired:
            logger.error(f"Lighthouse audit timed out for: {url}")
            return {}
        except Exception as e:
            logger.error(f"Error running Lighthouse: {str(e)}")
            return {}
    
    def _extract_metrics(self, lighthouse_data: dict) -> dict:
        """
        Extract relevant metrics from Lighthouse results
        """
        audits = lighthouse_data.get('audits', {})
        categories = lighthouse_data.get('categories', {})
        
        return {
            # Core Web Vitals
            'lcp': audits.get('largest-contentful-paint', {}).get('numericValue', 0) / 1000,
            'fid': audits.get('max-potential-fid', {}).get('numericValue', 0),
            'cls': audits.get('cumulative-layout-shift', {}).get('numericValue', 0),
            'fcp': audits.get('first-contentful-paint', {}).get('numericValue', 0) / 1000,
            'ttfb': audits.get('server-response-time', {}).get('numericValue', 0),
            
            # Scores
            'performance_score': categories.get('performance', {}).get('score', 0) * 100,
            'seo_score': categories.get('seo', {}).get('score', 0) * 100,
            'accessibility_score': categories.get('accessibility', {}).get('score', 0) * 100,
            'best_practices_score': categories.get('best-practices', {}).get('score', 0) * 100,
            
            # SEO Issues
            'seo_issues': self._extract_seo_issues(audits),
            'performance_issues': self._extract_performance_issues(audits)
        }
    
    def _extract_seo_issues(self, audits: dict) -> list:
        """
        Extract SEO-related issues
        """
        issues = []
        
        seo_audits = [
            'document-title',
            'meta-description',
            'http-status-code',
            'link-text',
            'is-crawlable',
            'robots-txt',
            'canonical'
        ]
        
        for audit_id in seo_audits:
            audit = audits.get(audit_id, {})
            if audit.get('score', 1) < 1:
                issues.append({
                    'type': audit_id,
                    'title': audit.get('title', ''),
                    'description': audit.get('description', ''),
                    'score': audit.get('score', 0)
                })
        
        return issues
    
    def _extract_performance_issues(self, audits: dict) -> list:
        """
        Extract performance-related issues
        """
        issues = []
        
        performance_audits = [
            'render-blocking-resources',
            'unused-css-rules',
            'unused-javascript',
            'modern-image-formats',
            'uses-text-compression',
            'uses-responsive-images'
        ]
        
        for audit_id in performance_audits:
            audit = audits.get(audit_id, {})
            if audit.get('score', 1) < 1:
                issues.append({
                    'type': audit_id,
                    'title': audit.get('title', ''),
                    'description': audit.get('description', ''),
                    'savings': audit.get('details', {}).get('overallSavingsMs', 0)
                })
        
        return issues
