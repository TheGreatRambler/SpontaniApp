variable "function_name" {
  type        = string
  description = "Name of Lambda"
}

variable "role_arn" {
  type        = string
  description = "IAM role ARN to be used by Lambda"
}

variable "handler" {
  type = string
  description = "The entry point for the Lambda. Format: \"filename.function\""
}

variable "runtime" {
  type = string
  description = "Runtime to be used by Lambda"
  default = "python3.12"
}

variable "source_dir" {
  type = string
  description = "Directory containing source code for lambda"
}

variable "env_vars" {
  type = map(string)
  description = "Environment variables to be passed to Lambda"
  default = {}
}