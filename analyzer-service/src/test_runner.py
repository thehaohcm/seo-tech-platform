import os
import json
import logging
import time
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from code_generator import PlaywrightCodeGenerator

logger = logging.getLogger(__name__)


class AutoTestRunner:
    """
    Automated test runner that generates and executes tests for web pages
    """
    
    def __init__(self):
        self.chrome_path = os.getenv('CHROME_BIN', '/usr/bin/chromium')
        
    def run_tests(self, url, page_id):
        """
        Run automated tests on a URL and return results with screenshot
        """
        logger.info(f"Running automated tests for: {url}")
        
        # Setup Chrome options
        chrome_options = Options()
        chrome_options.add_argument('--headless')
        chrome_options.add_argument('--no-sandbox')
        chrome_options.add_argument('--disable-dev-shm-usage')
        chrome_options.add_argument('--disable-gpu')
        chrome_options.add_argument('--window-size=1920,1080')
        chrome_options.binary_location = self.chrome_path
        
        results = {
            'status': 'passed',
            'total_tests': 0,
            'passed': 0,
            'failed': 0,
            'test_details': [],
            'screenshot_url': None,
            'execution_time': 0,
            'python_code': None  # Will store generated Playwright code
        }
        
        driver = None
        start_time = time.time()
        html_content = ""
        
        try:
            driver = webdriver.Chrome(options=chrome_options)
            driver.set_page_load_timeout(30)
            
            # Load the page
            driver.get(url)
            time.sleep(2)  # Wait for page to stabilize
            
            # Capture HTML for code generation
            try:
                html_content = driver.page_source
            except Exception as e:
                logger.warning(f"Could not capture HTML: {str(e)}")
            
            # Run test suite
            test_methods = [
                self._test_page_loads,
                self._test_title_exists,
                self._test_meta_description,
                self._test_h1_tags,
                self._test_images_have_alt,
                self._test_forms_have_labels,
                self._test_links_are_valid,
                self._test_no_console_errors
            ]
            
            for test_method in test_methods:
                try:
                    test_result = test_method(driver, url)
                    results['test_details'].append(test_result)
                    results['total_tests'] += 1
                    
                    if test_result['passed']:
                        results['passed'] += 1
                    else:
                        results['failed'] += 1
                        results['status'] = 'failed' if results['status'] == 'passed' else results['status']
                        
                except Exception as e:
                    logger.error(f"Test method {test_method.__name__} failed: {str(e)}")
                    results['test_details'].append({
                        'name': test_method.__name__.replace('_test_', '').replace('_', ' ').title(),
                        'passed': False,
                        'message': f'Test error: {str(e)}'
                    })
                    results['total_tests'] += 1
                    results['failed'] += 1
                    results['status'] = 'failed'
            
            # Take screenshot
            screenshot_path = f"/tmp/screenshot_{page_id}_{int(time.time())}.png"
            driver.save_screenshot(screenshot_path)
            
            # In production, upload to S3 or similar storage
            # For now, use a relative path
            results['screenshot_url'] = f"data:image/png;base64,{self._get_screenshot_base64(screenshot_path)}"
            
            # Clean up screenshot file
            if os.path.exists(screenshot_path):
                os.remove(screenshot_path)
            
            # Generate Playwright Python code
            if html_content:
                try:
                    code_generator = PlaywrightCodeGenerator()
                    python_code = code_generator.generate_page_object(html_content, url)
                    results['python_code'] = python_code
                    logger.info(f"Successfully generated Python code for {url}")
                except Exception as e:
                    logger.error(f"Failed to generate Python code: {str(e)}")
                    results['python_code'] = None
                
        except TimeoutException:
            results['status'] = 'failed'
            results['test_details'].append({
                'name': 'Page Load',
                'passed': False,
                'message': 'Page load timeout after 30 seconds'
            })
            results['total_tests'] += 1
            results['failed'] += 1
            
        except Exception as e:
            logger.error(f"Error running tests: {str(e)}", exc_info=True)
            results['status'] = 'failed'
            results['test_details'].append({
                'name': 'Test Execution',
                'passed': False,
                'message': f'Test execution error: {str(e)}'
            })
            results['total_tests'] += 1
            results['failed'] += 1
            
        finally:
            if driver:
                driver.quit()
            
            results['execution_time'] = round(time.time() - start_time, 2)
            
        # Determine final status
        if results['failed'] == 0:
            results['status'] = 'passed'
        elif results['failed'] < results['total_tests'] / 2:
            results['status'] = 'warning'
        else:
            results['status'] = 'failed'
            
        logger.info(f"Tests completed for {url}: {results['passed']}/{results['total_tests']} passed")
        return results
    
    def _test_page_loads(self, driver, url):
        """Test if page loads successfully"""
        try:
            # Check if we can find body element
            body = driver.find_element(By.TAG_NAME, 'body')
            return {
                'name': 'Page Loads Successfully',
                'passed': body is not None,
                'message': 'Page loaded and body element found'
            }
        except:
            return {
                'name': 'Page Loads Successfully',
                'passed': False,
                'message': 'Failed to find body element'
            }
    
    def _test_title_exists(self, driver, url):
        """Test if page has a title"""
        title = driver.title
        passed = bool(title and len(title) > 0)
        return {
            'name': 'Page Has Title',
            'passed': passed,
            'message': f'Title: "{title}"' if passed else 'No title found'
        }
    
    def _test_meta_description(self, driver, url):
        """Test if page has meta description"""
        try:
            meta = driver.find_element(By.XPATH, '//meta[@name="description"]')
            content = meta.get_attribute('content')
            passed = bool(content and len(content) > 0)
            return {
                'name': 'Meta Description Present',
                'passed': passed,
                'message': f'Meta description: {len(content)} characters' if passed else 'No meta description'
            }
        except NoSuchElementException:
            return {
                'name': 'Meta Description Present',
                'passed': False,
                'message': 'Meta description tag not found'
            }
    
    def _test_h1_tags(self, driver, url):
        """Test if page has H1 tags"""
        h1_tags = driver.find_elements(By.TAG_NAME, 'h1')
        count = len(h1_tags)
        passed = count > 0 and count <= 3
        
        if count == 0:
            message = 'No H1 tags found'
        elif count > 3:
            message = f'Too many H1 tags ({count}), should be 1-3'
        else:
            message = f'{count} H1 tag(s) found'
            
        return {
            'name': 'H1 Tags',
            'passed': passed,
            'message': message
        }
    
    def _test_images_have_alt(self, driver, url):
        """Test if images have alt attributes"""
        images = driver.find_elements(By.TAG_NAME, 'img')
        total = len(images)
        
        if total == 0:
            return {
                'name': 'Images Have Alt Text',
                'passed': True,
                'message': 'No images found on page'
            }
        
        missing_alt = 0
        for img in images:
            alt = img.get_attribute('alt')
            if not alt:
                missing_alt += 1
        
        passed = missing_alt == 0
        return {
            'name': 'Images Have Alt Text',
            'passed': passed,
            'message': f'{total - missing_alt}/{total} images have alt text' if not passed else 'All images have alt text'
        }
    
    def _test_forms_have_labels(self, driver, url):
        """Test if form inputs have associated labels"""
        inputs = driver.find_elements(By.TAG_NAME, 'input')
        
        # Filter out hidden inputs and buttons
        visible_inputs = [inp for inp in inputs if inp.get_attribute('type') not in ['hidden', 'submit', 'button']]
        total = len(visible_inputs)
        
        if total == 0:
            return {
                'name': 'Form Inputs Have Labels',
                'passed': True,
                'message': 'No form inputs found'
            }
        
        missing_labels = 0
        for inp in visible_inputs:
            inp_id = inp.get_attribute('id')
            aria_label = inp.get_attribute('aria-label')
            
            # Check if input has id and corresponding label
            has_label = False
            if inp_id:
                try:
                    driver.find_element(By.XPATH, f'//label[@for="{inp_id}"]')
                    has_label = True
                except NoSuchElementException:
                    pass
            
            if not has_label and not aria_label:
                missing_labels += 1
        
        passed = missing_labels == 0
        return {
            'name': 'Form Inputs Have Labels',
            'passed': passed,
            'message': f'{total - missing_labels}/{total} inputs have labels' if not passed else 'All inputs have labels'
        }
    
    def _test_links_are_valid(self, driver, url):
        """Test if links have href attributes"""
        links = driver.find_elements(By.TAG_NAME, 'a')
        total = len(links)
        
        if total == 0:
            return {
                'name': 'Links Are Valid',
                'passed': True,
                'message': 'No links found'
            }
        
        invalid = 0
        for link in links:
            href = link.get_attribute('href')
            if not href or href == '#' or href == 'javascript:void(0)':
                invalid += 1
        
        passed = invalid < total * 0.2  # Allow up to 20% placeholder links
        return {
            'name': 'Links Are Valid',
            'passed': passed,
            'message': f'{total - invalid}/{total} links have valid hrefs'
        }
    
    def _test_no_console_errors(self, driver, url):
        """Test if page has JavaScript console errors"""
        try:
            logs = driver.get_log('browser')
            errors = [log for log in logs if log['level'] == 'SEVERE']
            error_count = len(errors)
            
            passed = error_count == 0
            return {
                'name': 'No Console Errors',
                'passed': passed,
                'message': f'{error_count} console error(s) found' if not passed else 'No console errors'
            }
        except Exception as e:
            # Some browsers don't support console logs
            return {
                'name': 'No Console Errors',
                'passed': True,
                'message': 'Console logs not available'
            }
    
    def _get_screenshot_base64(self, screenshot_path):
        """Convert screenshot to base64"""
        import base64
        
        try:
            with open(screenshot_path, 'rb') as f:
                return base64.b64encode(f.read()).decode('utf-8')
        except Exception as e:
            logger.error(f"Error reading screenshot: {str(e)}")
            return None
