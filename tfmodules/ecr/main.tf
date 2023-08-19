variable "envs" {
  type    = map(any)
  default = {}
}

variable "others" {
  type    = map(any)
  default = {}
}

resource "aws_ecr_repository" "shopit" {
  name         = "shopit"
  force_delete = true
}

locals {
  shopit_repo_url  = aws_ecr_repository.shopit.repository_url
  shopit_image_url = "${local.shopit_repo_url}:latest"

  docker_build = <<-EOT
    echo "Building frontend..."
    cd frontend
    npm ci --force
    export REACT_APP_ENV=${var.envs.ENV}
    export REACT_APP_S3_URL=${var.others.shopit_bucket_url}
    export REACT_APP_DEFAULT_ITEM_IMAGE_NAME=${var.others.default_item_image_name}
    npm run build -- --ignore-path .env
    rm -rf ../backend/public
    mv build ../backend/public
    cd ..

    echo "Building backend..."
    aws ecr get-login-password --region ${var.envs.AWS_REGION} | docker login --username AWS --password-stdin ${local.shopit_repo_url}
    docker build -t shopit backend
    pwd
    docker tag shopit:latest ${local.shopit_image_url}
    docker push ${local.shopit_image_url}
  EOT
}

resource "null_resource" "build_push_docker_image" {
  provisioner "local-exec" {
    command = local.docker_build
  }
}

output "shopit_repo_url" {
  value = local.shopit_repo_url
}

output "shopit_image_url" {
  value = local.shopit_image_url
}

