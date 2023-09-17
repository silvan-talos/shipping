resource "aws_lb" "load_balancer" {
  name               = local.name
  load_balancer_type = "application"
  subnets            = var.vpc_subnets
  security_groups    = [var.security_group_id, aws_security_group.load_balancer.id]
}

resource "aws_lb_target_group" "shipping" {
  name        = local.name
  port        = var.output_port
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = var.vpc_id
  health_check {
    path     = "/ping"
    interval = 300
  }
}

resource "aws_lb_listener" "listener" {
  load_balancer_arn = aws_lb.load_balancer.arn
  port              = var.container_port
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.shipping.arn
  }
}

resource "aws_security_group" "load_balancer" {
  name   = local.name
  vpc_id = var.vpc_id
}

resource "aws_security_group_rule" "statistics" {
  type              = "ingress"
  from_port         = var.container_port
  to_port           = var.output_port
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.load_balancer.id
}

resource "aws_security_group_rule" "egress" {
  type              = "egress"
  from_port         = 0             # Allowing any incoming port
  to_port           = 0             # Allowing any outgoing port
  protocol          = "-1"          # Allowing any outgoing protocol
  cidr_blocks       = ["0.0.0.0/0"] # Allowing traffic out to all IP addresses
  security_group_id = aws_security_group.load_balancer.id
}
