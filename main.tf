terraform {
  backend "s3" {
    bucket         = "terra-form"
    dynamodb_table = "terra-form"
    key            = "shopit"
    region         = "ap-south-1"
  }
}

resource "aws_s3_object" "shopit_port" {
  bucket = "terra-form"
  key    = "shopit_port"
  source = var.PORT
}

provider "aws" {
  access_key = var.AWS_ACCESS_KEY_ID
  secret_key = var.AWS_SECRET_ACCESS_KEY
  region     = var.AWS_REGION
}

data "aws_vpc" "default_vpc" {
  default = true
}

locals {
  default_vpc_id = data.aws_vpc.default_vpc.id
}

module "s3" {
  source = "./tfmodules/s3"
}

module "dynamodb" {
  source = "./tfmodules/dynamodb"
}

module "ecr" {
  source = "./tfmodules/ecr"

  others = {
    "shopit_bucket_url"       = module.s3.shopit_bucket_url
    "default_item_image_name" = module.s3.default_item_image_name
  }

  envs = {
    "AWS_REGION" = var.AWS_REGION
    "ENV"        = var.ENV
  }
}

module "ecs" {
  source = "./tfmodules/ecs"

  others = {
    "default_vpc_id"   = local.default_vpc_id
    "shopit_image_url" = module.ecr.shopit_image_url
  }

  envs = {
    "AWS_ACCESS_KEY_ID"        = var.AWS_ACCESS_KEY_ID
    "AWS_SECRET_ACCESS_KEY"    = var.AWS_SECRET_ACCESS_KEY
    "AWS_REGION"               = var.AWS_REGION
    "ENV"                      = var.ENV
    "PORT"                     = var.PORT
    "DEFAULT_ITEM_IMAGE_NAME"  = module.s3.default_item_image_name
    "AWS_S3_BUCKET"            = module.s3.shopit_bucket_name
    "AWS_DYNAMO_DB_ITEM_TABLE" = module.dynamodb.shopit_item_table_name
  }
}
