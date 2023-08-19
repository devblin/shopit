# alwaays expected to be in ENV
variable "AWS_ACCESS_KEY" {}
variable "AWS_SECRET_KEY" {}
variable "ENV" {}

# always use default or provide in ENV with "TF_VAR_" prefix
variable "AWS_REGION" { default = "ap-south-1" }
variable "PORT" { default = 5000 }
variable "DEFAULT_ITEM_IMAGE_NAME" { default = "DEFAULT_ITEM_IMAGE_NAME" }
