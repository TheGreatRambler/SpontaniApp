# Create the API Gateway
resource "aws_api_gateway_rest_api" "api_gateway" {
  name = var.api_gateway_name
  binary_media_types = ["image/*", "application/octet-stream"] 
}

# Create a resource for each Lambda with its specified prefix
resource "aws_api_gateway_resource" "lambda_resource" {
  for_each = var.lambda_configs

  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  parent_id   = aws_api_gateway_rest_api.api_gateway.root_resource_id
  path_part   = each.value.prefix
}

# Local variable to generate combinations of lambda_keys and methods
locals {
  lambda_method_combinations = flatten([
    for lambda_key, lambda_value in var.lambda_configs : [
      for method in lambda_value.methods : {
        lambda_key = lambda_key
        method     = method
      }
    ]
  ])

  lambda_method_map = {
    for combo in local.lambda_method_combinations :
    "${combo.lambda_key}-${combo.method}" => combo
  }
}

# Create methods for each Lambda and method combination
resource "aws_api_gateway_method" "lambda_method" {
  for_each = local.lambda_method_map

  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  resource_id   = aws_api_gateway_resource.lambda_resource[each.value.lambda_key].id
  http_method   = upper(each.value.method)
  authorization = "NONE"
}

# Integrate Lambda functions with API Gateway methods
resource "aws_api_gateway_integration" "lambda_integration" {
  for_each = local.lambda_method_map

  rest_api_id             = aws_api_gateway_rest_api.api_gateway.id
  resource_id             = aws_api_gateway_resource.lambda_resource[each.value.lambda_key].id
  http_method             = upper(each.value.method)
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/${var.lambda_configs[each.value.lambda_key].lambda_arn}/invocations"

  depends_on = [aws_api_gateway_method.lambda_method]
}

# Create OPTIONS method for CORS for each resource
resource "aws_api_gateway_method" "options" {
  for_each = aws_api_gateway_resource.lambda_resource

  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  resource_id   = each.value.id
  http_method   = "OPTIONS"
  authorization = "NONE"
}

# Integration for OPTIONS method
resource "aws_api_gateway_integration" "options_integration" {
  for_each = aws_api_gateway_resource.lambda_resource

  rest_api_id          = aws_api_gateway_rest_api.api_gateway.id
  resource_id          = each.value.id
  http_method          = aws_api_gateway_method.options[each.key].http_method
  type                 = "MOCK"
  passthrough_behavior = "WHEN_NO_MATCH"

  request_templates = {
    "application/json" = "{\"statusCode\": 200}"
  }

  # Ensure the method exists before creating the integration
  depends_on = [aws_api_gateway_method.options]
}

# Method response for OPTIONS method
resource "aws_api_gateway_method_response" "options_response" {
  for_each = aws_api_gateway_resource.lambda_resource

  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  resource_id = each.value.id
  http_method = aws_api_gateway_method.options[each.key].http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers"     = true
    "method.response.header.Access-Control-Allow-Methods"     = true
    "method.response.header.Access-Control-Allow-Origin"      = true
    "method.response.header.Access-Control-Allow-Credentials" = true
  }
}

# Integration response for OPTIONS method
resource "aws_api_gateway_integration_response" "options_integration_response" {
  for_each = aws_api_gateway_resource.lambda_resource

  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  resource_id = each.value.id
  http_method = aws_api_gateway_method.options[each.key].http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers"     = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods"     = "'GET,POST,PUT,DELETE,OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin"      = "'*'"
    "method.response.header.Access-Control-Allow-Credentials" = "'true'"
  }

  # Ensure both the integration and method response are created first
  depends_on = [
    aws_api_gateway_integration.options_integration,
    aws_api_gateway_method_response.options_response
  ]
}



# Give API Gateway permission to invoke the Lambda functions
resource "aws_lambda_permission" "api_gateway_lambda" {
  for_each = var.lambda_configs

  statement_id  = "AllowAPIGatewayInvoke-${each.key}"
  action        = "lambda:InvokeFunction"
  function_name = each.value.lambda_arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.api_gateway.execution_arn}/*/*"
}

# Update API deployment
resource "aws_api_gateway_deployment" "api_deployment" {
  depends_on = [
    aws_api_gateway_integration.lambda_integration,
    aws_api_gateway_method.lambda_method,
    aws_api_gateway_integration.options_integration,
    aws_api_gateway_method.options,
    aws_api_gateway_integration_response.options_integration_response,
    aws_api_gateway_method_response.options_response
  ]
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  stage_name  = "prod"

  triggers = {
    redeployment = sha1(jsonencode(var.lambda_configs))
  }
}