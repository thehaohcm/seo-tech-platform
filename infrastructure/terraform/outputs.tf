output "rds_endpoint" {
  description = "RDS endpoint"
  value       = aws_db_instance.postgres.endpoint
}

output "redis_endpoint" {
  description = "Redis endpoint"
  value       = aws_elasticache_cluster.redis.cache_nodes[0].address
}

output "s3_bucket_name" {
  description = "S3 bucket name"
  value       = aws_s3_bucket.snapshots.id
}

output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}
