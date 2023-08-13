variable "AWS_REGION" { description = "The AWS region" }
variable "AWS_ACCESS_KEY" { description = "AWS access key" }
variable "AWS_SECRET_KEY" { description = "AWS secret key" }

provider "aws" {
  access_key = var.AWS_ACCESS_KEY
  secret_key = var.AWS_SECRET_KEY
  region     = var.AWS_REGION
}

resource "aws_s3_bucket" "shopit" {
  bucket = "shopit"
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
