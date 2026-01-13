import os
import logging
from typing import Dict, List
from openai import OpenAI
from langchain.prompts import PromptTemplate
from langchain_openai import ChatOpenAI
from langchain.chains import LLMChain

logger = logging.getLogger(__name__)


class SuggestionGenerator:
    """
    Generates AI-powered suggestions for fixing SEO and performance issues
    """
    
    def __init__(self):
        self.openai_api_key = os.getenv('OPENAI_API_KEY')
        if not self.openai_api_key:
            logger.warning("OpenAI API key not found. AI suggestions will be disabled.")
            self.enabled = False
        else:
            self.enabled = True
            self.llm = ChatOpenAI(
                model="gpt-4",
                temperature=0.7,
                openai_api_key=self.openai_api_key
            )
    
    def generate(self, url: str, lighthouse_results: dict, 
                 accessibility_results: dict, page_data: dict) -> str:
        """
        Generate AI suggestions based on audit results
        
        Args:
            url: The URL being audited
            lighthouse_results: Results from Lighthouse
            accessibility_results: Results from accessibility check
            page_data: Additional page data from crawler
            
        Returns:
            String with AI-generated suggestions
        """
        if not self.enabled:
            return "AI suggestions are disabled (no OpenAI API key configured)"
        
        try:
            # Prepare context
            context = self._prepare_context(
                url, lighthouse_results, accessibility_results, page_data
            )
            
            # Create prompt
            prompt = self._create_prompt(context)
            
            # Generate suggestions
            chain = LLMChain(llm=self.llm, prompt=prompt)
            suggestions = chain.run(context=context)
            
            logger.info(f"Generated AI suggestions for: {url}")
            return suggestions
            
        except Exception as e:
            logger.error(f"Error generating AI suggestions: {str(e)}")
            return f"Error generating suggestions: {str(e)}"
    
    def _prepare_context(self, url: str, lighthouse_results: dict,
                        accessibility_results: dict, page_data: dict) -> str:
        """
        Prepare context for the AI prompt
        """
        context_parts = [
            f"URL: {url}",
            f"\n## Page Information",
            f"Title: {page_data.get('title', 'N/A')}",
            f"Description: {page_data.get('description', 'N/A')}",
            f"H1 Tags: {', '.join(page_data.get('h1_tags', []))}",
        ]
        
        # Add Lighthouse metrics
        if lighthouse_results:
            context_parts.extend([
                f"\n## Performance Metrics",
                f"LCP: {lighthouse_results.get('lcp', 'N/A')}s",
                f"FID: {lighthouse_results.get('fid', 'N/A')}ms",
                f"CLS: {lighthouse_results.get('cls', 'N/A')}",
                f"Performance Score: {lighthouse_results.get('performance_score', 'N/A')}/100",
                f"SEO Score: {lighthouse_results.get('seo_score', 'N/A')}/100",
            ])
            
            # Add SEO issues
            seo_issues = lighthouse_results.get('seo_issues', [])
            if seo_issues:
                context_parts.append(f"\n## SEO Issues Found ({len(seo_issues)}):")
                for issue in seo_issues[:5]:  # Limit to top 5
                    context_parts.append(f"- {issue.get('title')}: {issue.get('description')}")
            
            # Add performance issues
            perf_issues = lighthouse_results.get('performance_issues', [])
            if perf_issues:
                context_parts.append(f"\n## Performance Issues Found ({len(perf_issues)}):")
                for issue in perf_issues[:5]:  # Limit to top 5
                    context_parts.append(f"- {issue.get('title')}: {issue.get('description')}")
        
        # Add accessibility issues
        accessibility_violations = accessibility_results.get('violations', [])
        if accessibility_violations:
            context_parts.append(f"\n## Accessibility Issues Found ({len(accessibility_violations)}):")
            for violation in accessibility_violations[:5]:  # Limit to top 5
                context_parts.append(f"- {violation.get('description')} (Impact: {violation.get('impact')})")
        
        return "\n".join(context_parts)
    
    def _create_prompt(self, context: str) -> PromptTemplate:
        """
        Create the prompt template for AI generation
        """
        template = """
You are an expert SEO and web performance consultant. Analyze the following website audit results and provide actionable, prioritized recommendations for improvement.

{context}

Please provide:
1. **Priority Issues**: The top 3-5 most critical issues that should be fixed first
2. **Quick Wins**: Easy fixes that can be implemented immediately
3. **Technical Recommendations**: Specific code changes or configurations needed
4. **Long-term Strategy**: Broader improvements for sustained optimization

Format your response in Markdown with clear sections and bullet points. Be specific and provide code examples where applicable.
"""
        
        return PromptTemplate(
            input_variables=["context"],
            template=template
        )
