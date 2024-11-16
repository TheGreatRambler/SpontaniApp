resource "aws_lambda_function" "this" {
  function_name = var.function_name
  role          = var.role_arn
  handler       = var.handler
  runtime       = var.runtime

  filename      = var.zip_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  environment {
    variables = var.env_vars
  }
}