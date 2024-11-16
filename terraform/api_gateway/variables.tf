variable "api_gateway_name" {
  description = "The name of the API Gateway."
  type        = string
}

variable "lambda_configs" {
  description = "A map of Lambda configuration, including ARN, HTTP methods, and path prefix."
  type = map(object({
    lambda_arn = string
    methods    = list(string)
    authorized_methods = list(string)
    prefix     = string
  }))
}