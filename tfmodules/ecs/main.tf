variable "envs" {
  type    = map(any)
  default = {}
}

variable "others" {
  type    = map(any)
  default = {}
}

locals {
  shopit_container_name      = "shopit"
  shopit_ec2_ami             = "ami-0205f72f24e39213b"
  shopit_ec2_ecs_role_policy = "AmazonEC2ContainerServiceforEC2Role"
}

data "aws_subnets" "shopit" {
  filter {
    name   = "vpc-id"
    values = [var.others.default_vpc_id]
  }
}

data "aws_iam_policy" "shopit" {
  name = local.shopit_ec2_ecs_role_policy
}

# create security groups
resource "aws_security_group" "shopit_all_to_alb" {
  name = "shopit_all_to_alb"

  ingress {
    from_port        = 80
    to_port          = 80
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "all"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_security_group" "shopit_alb_to_container" {
  name = "shopit_alb_to_container"

  ingress {
    from_port       = 0
    to_port         = 65535
    protocol        = "tcp"
    security_groups = [aws_security_group.shopit_all_to_alb.id]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "all"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

# create iam for ec2
data "aws_iam_policy_document" "shopit" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "shopit" {
  name               = "shopit"
  assume_role_policy = data.aws_iam_policy_document.shopit.json

  depends_on = [data.aws_iam_policy_document.shopit]
}

resource "aws_iam_role_policy_attachment" "shopit" {
  role       = aws_iam_role.shopit.name
  policy_arn = data.aws_iam_policy.shopit.arn

  depends_on = [aws_iam_role.shopit, data.aws_iam_policy.shopit]
}

resource "aws_iam_instance_profile" "shopit" {
  name = "shopit"
  role = aws_iam_role.shopit.name

  depends_on = [aws_iam_role.shopit]
}

# setup cluster
resource "aws_ecs_cluster" "shopit" {
  name = "shopit"
}

resource "aws_launch_template" "shopit" {
  name                   = "shopit"
  instance_type          = "t3.micro"
  image_id               = local.shopit_ec2_ami
  user_data              = base64encode("#!/bin/bash\necho ECS_CLUSTER=${aws_ecs_cluster.shopit.name} >> /etc/ecs/ecs.config")
  vpc_security_group_ids = [aws_security_group.shopit_alb_to_container.id]
  update_default_version = true

  iam_instance_profile {
    arn = aws_iam_instance_profile.shopit.arn
  }

  tag_specifications {
    resource_type = "instance"

    tags = {
      Name = aws_ecs_cluster.shopit.name
    }
  }

  depends_on = [aws_iam_instance_profile.shopit, aws_ecs_cluster.shopit, aws_security_group.shopit_alb_to_container]
}

resource "aws_autoscaling_group" "shopit" {
  name                      = "shopit"
  min_size                  = 1
  max_size                  = 1
  desired_capacity          = 1
  vpc_zone_identifier       = data.aws_subnets.shopit.ids
  force_delete              = true
  health_check_type         = "EC2"
  health_check_grace_period = 300

  launch_template {
    id = aws_launch_template.shopit.id
  }

  depends_on = [aws_launch_template.shopit, data.aws_subnets.shopit]
}

resource "aws_ecs_capacity_provider" "shopit_capacity" {
  name = "shopit_capacity"

  auto_scaling_group_provider {
    auto_scaling_group_arn = aws_autoscaling_group.shopit.arn
  }

  depends_on = [aws_autoscaling_group.shopit]
}

resource "aws_ecs_cluster_capacity_providers" "shopit" {
  cluster_name       = aws_ecs_cluster.shopit.name
  capacity_providers = [aws_ecs_capacity_provider.shopit_capacity.name]

  depends_on = [aws_ecs_cluster.shopit, aws_ecs_capacity_provider.shopit_capacity]
}

# create log group
resource "aws_cloudwatch_log_group" "shopit" {
  name              = "shopit"
  retention_in_days = 1
}

# create task definition
resource "aws_ecs_task_definition" "shopit" {
  family                   = "shopit"
  requires_compatibilities = ["EC2"]
  cpu                      = 1024
  memory                   = 350

  container_definitions = jsonencode([{
    name  = local.shopit_container_name
    image = var.others.shopit_image_url
    portMappings = [{
      containerPort = tonumber(var.envs.PORT)
      hostPort      = tonumber(var.envs.PORT)
    }]
    environment = [
      for key, value in var.envs : {
        name  = key
        value = value
      }
    ]
    memory = 300
    cpu    = 300
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.shopit.name
        "awslogs-region"        = var.envs.AWS_REGION,
        "awslogs-stream-prefix" = "container"
      }
    }
  }])

  depends_on = [aws_cloudwatch_log_group.shopit]
}

# create load balancer
resource "aws_lb" "shopit" {
  name               = "shopit"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.shopit_all_to_alb.id]
  subnets            = data.aws_subnets.shopit.ids

  depends_on = [aws_security_group.shopit_all_to_alb, data.aws_subnets.shopit]
}

resource "aws_lb_target_group" "shopit" {
  name_prefix      = "shopit"
  target_type      = "instance"
  port             = tonumber(var.envs.PORT)
  protocol         = "HTTP"
  protocol_version = "HTTP1"
  vpc_id           = var.others.default_vpc_id

  health_check {
    timeout             = 10
    interval            = 300
    unhealthy_threshold = 2
    healthy_threshold   = 2
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb_listener" "shopit" {
  load_balancer_arn = aws_lb.shopit.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "forward"

    forward {
      target_group {
        arn = aws_lb_target_group.shopit.arn
      }
    }
  }

  depends_on = [aws_lb_target_group.shopit, aws_lb.shopit]
}

# create ecs service
resource "aws_ecs_service" "shopit" {
  name                              = "shopit"
  cluster                           = aws_ecs_cluster.shopit.id
  task_definition                   = aws_ecs_task_definition.shopit.arn
  launch_type                       = "EC2"
  desired_count                     = 1
  scheduling_strategy               = "REPLICA"
  health_check_grace_period_seconds = 300
  triggers = {
    redeployment = timestamp()
  }

  deployment_controller {
    type = "ECS"
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.shopit.arn
    container_name   = local.shopit_container_name
    container_port   = var.envs.PORT
  }

  depends_on = [aws_ecs_cluster.shopit, aws_ecs_task_definition.shopit, aws_lb_target_group.shopit]
}

output "shopit_lb_dns" {
  value = aws_lb.shopit.dns_name

  depends_on = [aws_lb.shopit]
}
