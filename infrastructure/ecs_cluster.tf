resource "aws_ecs_cluster" "shipping" {
  name = local.name
}

resource "aws_ecs_task_definition" "shipping" {
  family = local.name
  container_definitions = jsonencode([
    {
      name      = local.name
      image     = aws_ecr_repository.shipping.repository_url
      cpu       = 256
      memory    = 512
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
      }]
    }
  ])
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  memory                   = 512
  cpu                      = 256
  execution_role_arn       = aws_iam_role.ecsTaskExecRole.arn
}

resource "aws_iam_role" "ecsTaskExecRole" {
  name = "ecsTaskExecRole"
  assume_role_policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = [
        {
          Sid    = ""
          Effect = "Allow"
          Principal = {
            Service = "ecs-tasks.amazonaws.com"
          }
          Action = "sts:AssumeRole"
        }
      ]
    }
  )
}

resource "aws_iam_role_policy_attachment" "ecsTaskExecRole_policy" {
  role       = aws_iam_role.ecsTaskExecRole.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_ecs_service" "shipping" {
  health_check_grace_period_seconds = 30
  name                              = local.name
  cluster                           = aws_ecs_cluster.shipping.id
  task_definition                   = aws_ecs_task_definition.shipping.arn
  launch_type                       = "FARGATE"
  desired_count                     = 1
  depends_on                        = [aws_iam_role_policy_attachment.ecsTaskExecRole_policy]

  load_balancer {
    container_name   = aws_ecs_task_definition.shipping.family
    container_port   = var.container_port
    target_group_arn = aws_lb_target_group.shipping.arn
  }

  network_configuration {
    subnets          = var.vpc_subnets
    security_groups  = [var.security_group_id, aws_security_group.load_balancer.id]
    assign_public_ip = true
  }
}

