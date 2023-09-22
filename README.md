# Shopit

A inventory management system written in **Golang**, **DynamoDB** for database, **S3** for file storage and **ReactJS** for UI.

[![Deploy](https://github.com/devblin/shopit/actions/workflows/deploy.yml/badge.svg?branch=production)](https://github.com/devblin/shopit/actions/workflows/deploy.yml)

## Features:

-   Basic CRUD funtionalities:
    -   Create inventory items
    -   Edit Them
    -   Delete Them
    -   View a list of them
-   Allow image uploads and storing image with generated thumbnails.

## Usage:

### Local Development:

To test the application locally, follow below steps:
- Install [localstack](https://docs.localstack.cloud/getting-started/installation/), [awslocal](https://github.com/localstack/awscli-local) and [tflocal](https://docs.localstack.cloud/user-guide/integrations/terraform/).
- Start localstack `localstack start -d`
- Execute below commands to setup basic infra for local terraform:
    - S3 bucket to handle terraform state: 
      ```
      awslocal s3api create-bucket --bucket terra-form --region ap-south-1  --create-bucket-configuration LocationConstraint=ap-south-1
      ```
    - DynamoDB table to handle state locking:
      ```
      awslocal dynamodb create-table --table-name terra-form --region ap-south-1 --key-schema AttributeName=LockID,KeyType=HASH --attribute-definitions AttributeName=LockID,AttributeType=S --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
      ```

- Provision aws in localstack using tflocal:
    - Create `dev.tfvars` with below contents in root of dir:
      ```sh
      AWS_ACCESS_KEY_ID     = "test"
      AWS_SECRET_ACCESS_KEY = "test"
      ENV                   = "dev"
      ```
    - Run `tflocal init -var-file=dev.tfvars`
    - Run `tflocal apply -var-file=dev.tfvars`
    - Use the `shopit_lb_dns` output's value to open the application.
