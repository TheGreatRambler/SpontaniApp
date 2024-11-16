output "lambda_arn" {
  description = "The ARN for the lambda"
  value       = aws_lambda_function.this.arn
}