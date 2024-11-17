terraform {
  required_providers {
    dotenv = {
      source  = "germanbrew/dotenv"
      version = "~> 1.0"
    }
  }
}

provider "aws" {
  region = "us-east-1" 
}

provider "dotenv" {}

data "dotenv" "env" {
  filename = "./.env"
}

locals {
  env_vars = data.dotenv.env.entries
}

resource "aws_iam_role" "lambda_role" {
  name = "lambda_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "s3_policy" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "null_resource" "run_build_script" {
  provisioner "local-exec" {
    command = "cd ../backend/get_lambda;./build.sh"
  }
}

module "get_lambda" {
    source = "./lambda"
    function_name = "get"
    zip_path = "../backend/get_lambda/getLambda.zip"
    env_vars = local.env_vars
    handler = "bootstrap"
    runtime = "provided.al2"
    role_arn = aws_iam_role.lambda_role.arn
}

module "api_gateway" {
    source = "./api_gateway"
    api_gateway_name = "spontaniapp"
    lambda_configs= {
        "get" = {
            lambda_arn = module.get_lambda.lambda_arn
            methods = ["GET"]
            prefix = "get"
        },
    }
}