resource "aws_lambda_function" "this" {
  function_name = var.function_name
  role          = var.role_arn
  handler       = var.handler
  runtime       = var.runtime

  filename      = var.zip_path
  source_code_hash = filesha256(var.zip_path)

  environment {
    variables = var.env_vars
  }
}