locals {
  region = "eu-central-1"
}

provider "aws" {
  region = local.region
}

terraform {
  backend "s3" {
    bucket = "terraform-credentials"
    key    = "checker-bot/terraform.tfstate"
    region = "eu-central-1"
  }
}

resource "random_id" "id" {
  byte_length = 8
}

resource "random_id" "random_path" {
  byte_length = 16
}

variable "telegram_token" {
  type      = string
  sensitive = true
}

resource "aws_ssm_parameter" "bot_token" {
  name  = "bot-token"
  type  = "SecureString"
  value = var.telegram_token
}

#resource "null_resource" "bot_build" {
#  triggers = {
#    build_number = timestamp()
#  }
#  provisioner "local-exec" {
#    command = "cd ${path.root} && mkdir binaries && mkdir binaries/bot && cd ${path.root}/src && env GOOS=linux GOARCH=amd64 go build -o ${path.root}/binaries/bot/main ."
#  }
#}

data "archive_file" "bot_lambda_zip" {
#  depends_on  = [null_resource.bot_build]
  type        = "zip"
  output_path = "/tmp/bot-${random_id.id.hex}.zip"
  source_dir  = "${path.root}/binaries/bot"
}

resource "aws_lambda_function" "bot_lambda" {
  function_name = "bot-${random_id.id.hex}-function"

  filename         = data.archive_file.bot_lambda_zip.output_path
  source_code_hash = data.archive_file.bot_lambda_zip.output_base64sha256
  environment {
    variables = {
      TOKEN_PARAMETER = aws_ssm_parameter.bot_token.name
      REGION          = local.region
    }
  }

  timeout = 30
  handler = "main"
  runtime = "go1.x"
  role    = aws_iam_role.lambda_exec.arn
}

resource "aws_lambda_invocation" "set_webhook" {
  function_name = aws_lambda_function.bot_lambda.function_name

  input = <<JSON
{
    "version": "2.0",
    "routeKey": "$default",
    "rawPath": "/init-bot",
    "rawQueryString": "",
    "headers": {
    },
    "requestContext": {
        "accountId": "anonymous",
        "apiId": "zkpmstc7celktxmj24j4frgeeq0tnsbi",
        "domainName": "zkpmstc7celktxmj24j4frgeeq0tnsbi.lambda-url.eu-central-1.on.aws",
        "domainPrefix": "zkpmstc7celktxmj24j4frgeeq0tnsbi",
        "http": {
            "method": "GET",
            "path": "/init-bot",
            "protocol": "HTTP/1.1",
            "sourceIp": "213.108.105.86",
            "userAgent": "PostmanRuntime/7.30.0"
        },
        "requestId": "c83417f3-9c82-4780-b9f4-fb5697cae66c",
        "routeKey": "$default",
        "stage": "$default",
        "time": "31/Jan/2023:17:29:35 +0000",
        "timeEpoch": 1675186175633
    },
    "isBase64Encoded": false
}
JSON
}

data "aws_iam_policy_document" "lambda_exec_role_policy" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:*:*:*"
    ]
  }
  statement {
    actions = [
      "ssm:GetParameter",
    ]
    resources = [
      aws_ssm_parameter.bot_token.arn
    ]
  }
}

resource "aws_cloudwatch_log_group" "bot_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.bot_lambda.function_name}"
  retention_in_days = 14
}

resource "aws_iam_role_policy" "lambda_exec_role" {
  role   = aws_iam_role.lambda_exec.id
  policy = data.aws_iam_policy_document.lambda_exec_role_policy.json
}

resource "aws_iam_role" "lambda_exec" {
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

# api gw

resource "aws_apigatewayv2_api" "api" {
  name          = "api-${random_id.id.hex}"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "bot_api" {
  api_id           = aws_apigatewayv2_api.api.id
  integration_type = "AWS_PROXY"

  integration_method     = "POST"
  integration_uri        = aws_lambda_function.bot_lambda.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "bot_api" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "ANY /${random_id.random_path.hex}/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.bot_api.id}"
}

resource "aws_apigatewayv2_stage" "api" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_lambda_permission" "bot_apigw" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.bot_lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}