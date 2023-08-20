# s3 outputs
output "default_item_image_name" {
  value = module.s3.default_item_image_name
}

output "shopit_bucket_name" {
  value = module.s3.shopit_bucket_name
}

output "shopit_bucket_url" {
  value = module.s3.shopit_bucket_url
}

output "shopit_port" {
  value = module.s3.shopit_port
}

# ecr outputs
output "shopit_repo_url" {
  value = module.ecr.shopit_repo_url
}

output "shopit_image_url" {
  value = module.ecr.shopit_image_url
}

# dynamodb outputs
output "shopit_item_table_name" {
  value = module.dynamodb.shopit_item_table_name
}

# ecs outputs
output "shopit_lb_dns" {
  value = module.ecs.shopit_lb_dns
}
