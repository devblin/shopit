resource "aws_s3_bucket" "shopit" {
  bucket        = "shop-it"
  force_destroy = true
}

resource "aws_s3_bucket_acl" "shopit" {
  bucket = aws_s3_bucket.shopit.id
  acl    = "public-read"

  depends_on = [aws_s3_bucket.shopit, aws_s3_bucket_ownership_controls.shopit]
}

resource "aws_s3_bucket_ownership_controls" "shopit" {
  bucket = aws_s3_bucket.shopit.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }

  depends_on = [aws_s3_bucket_public_access_block.shopit]
}

resource "aws_s3_bucket_public_access_block" "shopit" {
  bucket = aws_s3_bucket.shopit.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false

  depends_on = [aws_s3_bucket.shopit]
}

locals {
  default_item_image_name = "default-product"
}

resource "aws_s3_object" "default-product400" {
  bucket = aws_s3_bucket.shopit.id
  key    = "default-product400.jpg"
  source = "frontend/public/${local.default_item_image_name}400.jpg"
  acl    = aws_s3_bucket_acl.shopit.acl

  depends_on = [aws_s3_bucket.shopit, aws_s3_bucket_acl.shopit]
}

resource "aws_s3_object" "default-product64" {
  bucket = aws_s3_bucket.shopit.id
  key    = "default-product64.jpg"
  source = "frontend/public/${local.default_item_image_name}64.jpg"
  acl    = aws_s3_bucket_acl.shopit.acl

  depends_on = [aws_s3_bucket.shopit, aws_s3_bucket_acl.shopit]
}

output "default_item_image_name" {
  value = "${local.default_item_image_name}.jpg"
}

output "shopit_bucket_name" {
  value = aws_s3_bucket.shopit.bucket

  depends_on = [aws_s3_bucket.shopit]
}

output "shopit_bucket_url" {
  value = aws_s3_bucket.shopit.bucket_domain_name

  depends_on = [aws_s3_bucket.shopit]
}

# get existing port
data "aws_s3_object" "shopit_port" {
  bucket = "terra-form"
  key    = "shopit_port"
}

locals {
  shopit_port = tonumber(coalesce(try(data.aws_s3_object.shopit_port.body, null), 5000)) == 5000 ? 5001 : 5000
}

resource "null_resource" "shopit_port" {
  provisioner "local-exec" {
    command = <<EOT
    aws s3api put-object \
    --bucket terra-form \
    --key shopit_port \
    --body ${local.shopit_port} \
    --content-type text/plain
    EOT

    on_failure = fail
  }

  triggers = {
    always_run = timestamp()
  }

  depends_on = [data.aws_s3_object.shopit_port]
}

output "shopit_port" {
  value = local.shopit_port
}
