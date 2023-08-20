# always expected to be in ENV
variable "AWS_ACCESS_KEY_ID" {}
variable "AWS_SECRET_ACCESS_KEY" {}
variable "ENV" {}

# always use default or provide in ENV with "TF_VAR_" prefix
variable "AWS_REGION" { default = "ap-south-1" }
variable "DEFAULT_ITEM_IMAGE_NAME" { default = "DEFAULT_ITEM_IMAGE_NAME" }
