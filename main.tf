variable "AWS_REGION" { default = "" }
variable "AWS_ACCESS_KEY" { default = "" }
variable "AWS_SECRET_KEY" { default = "" }

provider "aws" {
  access_key = var.AWS_ACCESS_KEY
  secret_key = var.AWS_SECRET_KEY
  region     = var.AWS_REGION
}

# S3 BUCKET SETUP
resource "aws_s3_bucket" "shopit" {
  bucket = "shopit"
}

resource "aws_s3_bucket_acl" "shopit" {
  bucket = aws_s3_bucket.shopit.bucket
  acl    = "public-read"
}

resource "aws_s3_object" "default-product400" {
  bucket = aws_s3_bucket.shopit.id
  key    = "default-product400.jpg"
  source = "frontend/public/default-product400.jpg"
}

resource "aws_s3_object" "default-product64" {
  bucket = aws_s3_bucket.shopit.id
  key    = "default-product64.jpg"
  source = "frontend/public/default-product64.jpg"
}

# DYNAMODB SETUP
# DYNAMODB ITEM TABLE SETUP
resource "aws_dynamodb_table" "shopit_item" {
  name           = "shopit_item"
  hash_key       = "Id"
  read_capacity  = 30
  write_capacity = 30

  attribute {
    name = "Id"
    type = "S"
  }
}

# ECR SETUP
resource "aws_ecr_repository" "shopit" {
  name = "shopit"
}

output "shopit_repo_url" {
  value = aws_ecr_repository.shopit.repository_url
}
