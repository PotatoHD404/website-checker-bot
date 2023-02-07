variable "app_name" {
  type    = string
  default = "checker-bot"
}

variable "region" {
  type    = string
  default = "eu-central-1"
}

variable "telegram_token" {
  type      = string
  sensitive = true
}

provider "aws" {
  region = var.region
}

terraform {
  backend "s3" {
    bucket = "terraform-credentials"
    key    = "checker-bot/terraform.tfstate"
    region = "eu-central-1"
  }
}

resource "random_password" "id" {
  length  = 8
  special = false
  numeric = true
  upper   = false
  lower   = true
}

resource "random_password" "random_path" {
  length  = 16
  special = false
  numeric = true
  upper   = false
  lower   = true
}

resource "aws_ssm_parameter" "bot_token" {
  name  = "bot-token"
  type  = "SecureString"
  value = var.telegram_token
}

data "archive_file" "bot_lambda_zip" {
  type        = "zip"
  output_path = "/tmp/bot-${random_password.id.result}.zip"
  source_dir  = "${path.root}/binaries/bot"
}

resource "aws_lambda_function" "bot_lambda" {
  function_name = "bot-${random_password.id.result}-function"

  filename         = data.archive_file.bot_lambda_zip.output_path
  source_code_hash = data.archive_file.bot_lambda_zip.output_base64sha256
  environment {
    variables = {
      TOKEN_PARAMETER = aws_ssm_parameter.bot_token.name
      REGION          = var.region
      domain          = sensitive(aws_apigatewayv2_stage.api.invoke_url)
      path_key        = random_password.random_path.result
    }
  }

  timeout = 30
  handler = "main"
  runtime = "go1.x"
  role    = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_event_rule" "every_five_minutes" {
  name = "every-five-minutes"
  description = "Fires every five minutes"
  schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "check_website_every_five_minutes" {
  rule = aws_cloudwatch_event_rule.every_five_minutes.name
  target_id = "check_website_every_five_minutes"
  arn = aws_lambda_function.bot_lambda.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_checker" {
  statement_id = "AllowExecutionFromCloudWatch"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.bot_lambda.function_name
  principal = "events.amazonaws.com"
  source_arn = aws_cloudwatch_event_rule.every_five_minutes.arn
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

  statement {
    actions = [
      "dynamodb:DeleteItem",
      "dynamodb:DescribeTable",
      "dynamodb:GetItem",
      "dynamodb:GetRecords",
      "dynamodb:PutItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:UpdateItem",
      "dynamodb:UpdateTable",
    ]

    resources = [
      aws_dynamodb_table.admins.arn,
      aws_dynamodb_table.websites.arn
    ]

    effect = "Allow"
  }
}

resource "aws_dynamodb_table" "admins" {
  name         = "checker-admins"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "chat_id"

  attribute {
    name = "chat_id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "websites" {
  name         = "checker-websites"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "name"

  attribute {
    name = "name"
    type = "S"
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

resource "aws_apigatewayv2_api" "api" {
  name          = "api-${random_password.id.result}"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "api" {
  api_id           = aws_apigatewayv2_api.api.id
  integration_type = "AWS_PROXY"

  integration_method     = "POST"
  integration_uri        = aws_lambda_function.bot_lambda.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "api" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "ANY /${random_password.random_path.result}/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.api.id}"
}

resource "aws_apigatewayv2_stage" "api" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_lambda_permission" "api" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.bot_lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}

resource "null_resource" "init_bot" {
  depends_on = [
    aws_lambda_permission.api,
    aws_apigatewayv2_stage.api,
    aws_lambda_function.bot_lambda,
    aws_iam_role_policy.lambda_exec_role,
    aws_cloudwatch_log_group.bot_log_group,
  ]
  triggers = {
    build_number = timestamp()
  }
  provisioner "local-exec" {
    command = "curl ${aws_apigatewayv2_stage.api.invoke_url}${random_password.random_path.result}/init-bot"
  }
}

