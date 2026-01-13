-- Database initialization script for SEO Tech Platform
-- PostgreSQL 15+

-- Create database (if not exists via env)
-- CREATE DATABASE seo_platform;

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Bảng quản lý User
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Bảng quản lý Project (Website cần Audit)
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    domain VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    settings JSONB, -- Cấu hình crawl (UserAgent, Cookies, Timeout, Max Pages)
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Bảng lưu kết quả từng lần chạy (Audit Run)
CREATE TABLE IF NOT EXISTS audit_runs (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'queued', -- 'queued', 'crawling', 'analyzing', 'completed', 'failed'
    started_at TIMESTAMP DEFAULT NOW(),
    finished_at TIMESTAMP,
    overall_score INT, -- Điểm trung bình (0-100)
    total_pages INT DEFAULT 0,
    error_message TEXT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Bảng chi tiết từng URL trong 1 lần chạy
CREATE TABLE IF NOT EXISTS page_audits (
    id SERIAL PRIMARY KEY,
    run_id INT NOT NULL,
    url TEXT NOT NULL,
    status_code INT,
    load_time_ms INT,
    
    -- Core Web Vitals
    lcp_score FLOAT, -- Largest Contentful Paint
    fid_score FLOAT, -- First Input Delay
    cls_score FLOAT, -- Cumulative Layout Shift
    fcp_score FLOAT, -- First Contentful Paint
    ttfb_score FLOAT, -- Time to First Byte
    
    -- SEO Metrics
    title VARCHAR(500),
    meta_description TEXT,
    h1_tags JSONB,
    canonical_url TEXT,
    has_robots_meta BOOLEAN DEFAULT FALSE,
    
    -- Issues and Suggestions
    seo_issues JSONB, -- Chứa mảng các lỗi tìm thấy
    accessibility_issues JSONB,
    performance_issues JSONB,
    ai_suggestions TEXT, -- Nội dung AI gợi ý cách fix
    
    -- Storage
    html_snapshot_path VARCHAR(255), -- Link tới S3 chứa file HTML lúc crawl
    screenshot_path VARCHAR(255),
    
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (run_id) REFERENCES audit_runs(id) ON DELETE CASCADE
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_runs_project_id ON audit_runs(project_id);
CREATE INDEX IF NOT EXISTS idx_audit_runs_status ON audit_runs(status);
CREATE INDEX IF NOT EXISTS idx_page_audits_run_id ON page_audits(run_id);
CREATE INDEX IF NOT EXISTS idx_page_audits_url ON page_audits(url);

-- Insert sample data for development
INSERT INTO users (email, password_hash, full_name) 
VALUES 
    ('admin@example.com', '$2a$10$YourHashedPasswordHere', 'Admin User'),
    ('developer@example.com', '$2a$10$YourHashedPasswordHere', 'Developer User')
ON CONFLICT (email) DO NOTHING;

-- Sample project
INSERT INTO projects (user_id, domain, name, settings)
VALUES 
    (1, 'https://example.com', 'Example Website', 
     '{"userAgent": "SEO-Bot/1.0", "maxPages": 100, "timeout": 30000}')
ON CONFLICT DO NOTHING;

-- Comments
COMMENT ON TABLE projects IS 'Stores website projects to be audited';
COMMENT ON TABLE audit_runs IS 'Tracks each audit execution for a project';
COMMENT ON TABLE page_audits IS 'Detailed audit results for each page crawled';
COMMENT ON COLUMN page_audits.seo_issues IS 'JSON array of SEO issues found (missing title, meta, etc)';
COMMENT ON COLUMN page_audits.ai_suggestions IS 'AI-generated recommendations for fixing issues';
