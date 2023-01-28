provider "aws" {
  region = "eu-central-1"
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

data "external" "prebuild" {
  program = [
    "bash", "-c", <<EOT
mkdir -p bin
EOT
  ]
  working_dir = "${path.module}/binaries"
}

data "external" "checker_build" {
  program = [
    "bash", "-c", <<EOT
env GOOS=linux GOARCH=amd64 go build -o ../../binaries/checker
EOT
  ]
  working_dir = "${path.module}/src/checker"
}

data "external" "bot_build" {
  program = [
    "bash", "-c", <<EOT
env GOOS=linux GOARCH=amd64 go build -o ../../binaries/bot
EOT
  ]
  working_dir = "${path.module}/src/bot"
}


data "archive_file" "checker_lambda_zip" {
  type        = "zip"
  output_path = "/tmp/lambda-${random_id.id.hex}.zip"
  source_dir  = "${data.external.checker_build.working_dir}/${data.external.checker_build.result.dest}"
}

data "archive_file" "bot_lambda_zip" {
  type        = "zip"
  output_path = "/tmp/bot-${random_id.id.hex}.zip"
  source_dir  = "${data.external.bot_build.working_dir}/${data.external.bot_build.result.dest}"
}

resource "aws_lambda_function" "checker_lambda" {
  function_name = "checker-${random_id.id.hex}-function"

  filename         = data.archive_file.checker_lambda_zip.output_path
  source_code_hash = data.archive_file.checker_lambda_zip.output_base64sha256
  environment {
    variables = {
      domain          = aws_apigatewayv2_api.api.api_endpoint
      path_key        = random_id.random_path.hex
      token_parameter = aws_ssm_parameter.bot_token.name
    }
  }

  timeout = 30
  handler = "checker"
  runtime = "go1.x"
  role    = aws_iam_role.lambda_exec.arn
}

resource "aws_lambda_function" "bot_lambda" {
  function_name = "bot-${random_id.id.hex}-function"

  filename         = data.archive_file.bot_lambda_zip.output_path
  source_code_hash = data.archive_file.bot_lambda_zip.output_base64sha256
  environment {
    variables = {
      token_parameter = aws_ssm_parameter.bot_token.name
    }
  }

  timeout = 30
  handler = "bot"
  runtime = "go1.x"
  role    = aws_iam_role.lambda_exec.arn
}

data "aws_lambda_invocation" "set_webhook" {
  function_name = aws_lambda_function.bot_lambda.function_name

  input = <<JSON
{
	"setWebhook": true
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

resource "aws_cloudwatch_log_group" "checker_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.checker_lambda.function_name}"
  retention_in_days = 14
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

resource "aws_apigatewayv2_integration" "checker_api" {
  api_id           = aws_apigatewayv2_api.api.id
  integration_type = "AWS_PROXY"

  integration_method     = "POST"
  integration_uri        = aws_lambda_function.checker_lambda.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_integration" "bot_api" {
  api_id           = aws_apigatewayv2_api.api.id
  integration_type = "AWS_PROXY"

  integration_method     = "POST"
  integration_uri        = aws_lambda_function.bot_lambda.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "checker_api" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "ANY /checker/${random_id.random_path.hex}/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.checker_api.id}"
}

resource "aws_apigatewayv2_route" "bot_api" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "ANY /bot/${random_id.random_path.hex}/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.bot_api.id}"
}

resource "aws_apigatewayv2_stage" "api" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_lambda_permission" "checker_apigw" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.checker_lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}

resource "aws_lambda_permission" "bot_apigw" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.bot_lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}