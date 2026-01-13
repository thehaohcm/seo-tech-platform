import os
import logging
from sqlalchemy import create_engine, Column, Integer, String, Float, Text, TIMESTAMP, JSON, Boolean
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from datetime import datetime

logger = logging.getLogger(__name__)

Base = declarative_base()


class PageAudit(Base):
    __tablename__ = 'page_audits'
    
    id = Column(Integer, primary_key=True)
    run_id = Column(Integer, nullable=False)
    url = Column(Text, nullable=False)
    status_code = Column(Integer)
    load_time_ms = Column(Integer)
    
    # Core Web Vitals
    lcp_score = Column(Float)
    fid_score = Column(Float)
    cls_score = Column(Float)
    fcp_score = Column(Float)
    ttfb_score = Column(Float)
    
    # SEO
    title = Column(String(500))
    meta_description = Column(Text)
    h1_tags = Column(JSON)
    canonical_url = Column(Text)
    has_robots_meta = Column(Boolean)
    
    # Issues
    seo_issues = Column(JSON)
    accessibility_issues = Column(JSON)
    performance_issues = Column(JSON)
    ai_suggestions = Column(Text)
    
    # Storage
    html_snapshot_path = Column(String(255))
    screenshot_path = Column(String(255))
    
    created_at = Column(TIMESTAMP, default=datetime.utcnow)


class Database:
    """
    Database connection and operations
    """
    
    def __init__(self):
        db_host = os.getenv('DB_HOST', 'localhost')
        db_port = os.getenv('DB_PORT', '5432')
        db_name = os.getenv('DB_NAME', 'seo_platform')
        db_user = os.getenv('DB_USER', 'seouser')
        db_password = os.getenv('DB_PASSWORD', 'seopass')
        
        connection_string = f"postgresql://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}"
        
        self.engine = create_engine(connection_string)
        Session = sessionmaker(bind=self.engine)
        self.session = Session()
    
    def save_page_audit(self, audit_data: dict):
        """
        Save page audit results to database
        """
        try:
            page_audit = PageAudit(**audit_data)
            self.session.add(page_audit)
            self.session.commit()
            logger.info(f"Saved audit for URL: {audit_data.get('url')}")
        except Exception as e:
            self.session.rollback()
            logger.error(f"Error saving page audit: {str(e)}")
            raise
    
    def close(self):
        """
        Close database connection
        """
        self.session.close()
