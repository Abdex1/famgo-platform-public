terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC
resource "aws_vpc" "famgo" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "famgo-vpc"
  }
}

# Public Subnets
resource "aws_subnet" "public_1" {
  vpc_id                  = aws_vpc.famgo.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "${var.aws_region}a"
  map_public_ip_on_launch = true

  tags = {
    Name = "famgo-public-1"
  }
}

resource "aws_subnet" "public_2" {
  vpc_id                  = aws_vpc.famgo.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = "${var.aws_region}b"
  map_public_ip_on_launch = true

  tags = {
    Name = "famgo-public-2"
  }
}

# Private Subnets
resource "aws_subnet" "private_1" {
  vpc_id            = aws_vpc.famgo.id
  cidr_block        = "10.0.10.0/24"
  availability_zone = "${var.aws_region}a"

  tags = {
    Name = "famgo-private-1"
  }
}

resource "aws_subnet" "private_2" {
  vpc_id            = aws_vpc.famgo.id
  cidr_block        = "10.0.11.0/24"
  availability_zone = "${var.aws_region}b"

  tags = {
    Name = "famgo-private-2"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "famgo" {
  vpc_id = aws_vpc.famgo.id

  tags = {
    Name = "famgo-igw"
  }
}

# Route Table
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.famgo.id

  route {
    cidr_block      = "0.0.0.0/0"
    gateway_id      = aws_internet_gateway.famgo.id
  }

  tags = {
    Name = "famgo-public-rt"
  }
}

resource "aws_route_table_association" "public_1" {
  subnet_id      = aws_subnet.public_1.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_2" {
  subnet_id      = aws_subnet.public_2.id
  route_table_id = aws_route_table.public.id
}

# RDS PostgreSQL
resource "aws_db_instance" "famgo" {
  identifier            = "famgo-db"
  engine                = "postgres"
  engine_version        = "16.0"
  instance_class        = "db.t3.small"
  allocated_storage     = 20
  db_name               = "famgo_platform"
  username              = "famgo"
  password              = random_password.db_password.result
  skip_final_snapshot   = true
  publicly_accessible   = false
  db_subnet_group_name  = aws_db_subnet_group.famgo.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  tags = {
    Name = "famgo-db"
  }
}

# DB Subnet Group
resource "aws_db_subnet_group" "famgo" {
  name       = "famgo-db-subnet"
  subnet_ids = [aws_subnet.private_1.id, aws_subnet.private_2.id]

  tags = {
    Name = "famgo-db-subnet-group"
  }
}

# Security Groups
resource "aws_security_group" "rds" {
  name        = "famgo-rds-sg"
  description = "Security group for RDS"
  vpc_id      = aws_vpc.famgo.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "famgo-rds-sg"
  }
}

# ElastiCache Redis
resource "aws_elasticache_cluster" "famgo" {
  cluster_id           = "famgo-redis"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379
  security_group_ids   = [aws_security_group.redis.id]
  subnet_group_name    = aws_elasticache_subnet_group.famgo.name

  tags = {
    Name = "famgo-redis"
  }
}

# ElastiCache Subnet Group
resource "aws_elasticache_subnet_group" "famgo" {
  name       = "famgo-cache-subnet"
  subnet_ids = [aws_subnet.private_1.id, aws_subnet.private_2.id]
}

# Security Group for Redis
resource "aws_security_group" "redis" {
  name        = "famgo-redis-sg"
  description = "Security group for Redis"
  vpc_id      = aws_vpc.famgo.id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "famgo-redis-sg"
  }
}

# Random Password for RDS
resource "random_password" "db_password" {
  length  = 16
  special = true
}

# Outputs
output "rds_endpoint" {
  value       = aws_db_instance.famgo.endpoint
  description = "RDS endpoint"
}

output "redis_endpoint" {
  value       = aws_elasticache_cluster.famgo.cache_nodes[0].address
  description = "Redis endpoint"
}
