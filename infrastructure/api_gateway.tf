resource "aws_apigatewayv2_api" "shipping" {
  name          = local.name
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "default" {
  api_id             = aws_apigatewayv2_api.shipping.id
  integration_type   = "HTTP_PROXY"
  integration_uri    = var.integration_uri
  integration_method = "ANY"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.shipping.id
  name        = "$default"
  auto_deploy = true
}
resource "aws_apigatewayv2_route" "default" {
  api_id    = aws_apigatewayv2_api.shipping.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.default.id}"
}
