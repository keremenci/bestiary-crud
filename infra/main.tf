terraform {
  backend "s3" {
    bucket = "bestiary-tfstate"
    key    = "terraform/state.tfstate"
    region = "eu-central-1"
    dynamodb_table = "bestiary-tfstate-lock"
    encrypt = true
  }
}

provider "aws" {
  region = "eu-central-1"
}

# Networking stuff

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
}

resource "aws_subnet" "main" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
}

resource "aws_route_table" "main" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
}

resource "aws_route_table_association" "a" {
  subnet_id      = aws_subnet.main.id
  route_table_id = aws_route_table.main.id
}


resource "aws_security_group" "instance" {
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # SSH access from anywhere. Restrict in the future
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Allow HTTP
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"] # Allow outgoing
  }
}

# Key pairs for the ec2 instance

resource "aws_key_pair" "deployer" {
  key_name   = "bestiary-deployer-ed25519"
  public_key = file("./secrets/deployer_ed25519.pub")
}

resource "aws_key_pair" "sysadmin" {
  key_name   = "bestiary-sysadmin-ed25519"
  public_key = file("./secrets/sysadmin_ed25519.pub")
}

# IAM role for EC2 to access ECR

resource "aws_iam_role" "ec2_role" {
  name = "bestiary-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_instance_profile" "ec2_instance_profile" {
  name = "bestiary-ec2-instance-profile"
  role = aws_iam_role.ec2_role.name
}

resource "aws_iam_policy" "ecr_policy" {
  name        = "bestiary-ecr-policy"
  description = "Policy to allow EC2 to pull images from ECR"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "attach_ecr_policy" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = aws_iam_policy.ecr_policy.arn
}


# EC2 instance that will run the app

resource "aws_instance" "app" {
  ami                    = "ami-0e872aee57663ae2d" # Ubuntu Server 24.04 LTS
  instance_type          = "t2.micro"
  subnet_id              = aws_subnet.main.id
  vpc_security_group_ids = [aws_security_group.instance.id]

  key_name = aws_key_pair.sysadmin.key_name

  tags = {
    Name = "bestiary-app"
  }

  user_data = <<-EOF
              #!/bin/bash
              for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done
              # Add Docker's official GPG key:
              sudo apt-get -qq update
              sudo apt-get -qq --yes --force-yes install ca-certificates curl
              sudo install -m 0755 -d /etc/apt/keyrings
              sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
              sudo chmod a+r /etc/apt/keyrings/docker.asc

              # Add the repository to Apt sources:
              echo \
                "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
                $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
                sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
              sudo apt-get -qq update
              sudo apt-get -qq --yes --force-yes install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
              sudo systemctl enable docker.service
              sudo systemctl enable containerd.service
              sudo groupadd docker

              # Create user, password generated with python3 crypt
              sudo useradd -m -s /bin/bash -U -G sudo,docker,adm,dip,lxd,cdrom bestiary
              sudo mkdir -p /home/bestiary/.ssh
              sudo chmod 700 /home/bestiary/.ssh

              # Add sysadmin key to authorized_keys for ubuntu user
              echo "${file("./secrets/sysadmin_ed25519.pub")}" | sudo tee -a /home/ubuntu/.ssh/authorized_keys

              # Add both sysadmin and deployer keys to authorized_keys for bestiary user
              echo "${file("./secrets/sysadmin_ed25519.pub")}" | sudo tee -a /home/bestiary/.ssh/authorized_keys
              echo "${file("./secrets/deployer_ed25519.pub")}" | sudo tee -a /home/bestiary/.ssh/authorized_keys

              sudo chmod 600 /home/bestiary/.ssh/authorized_keys
              sudo chown -R bestiary:bestiary /home/bestiary

              # Install credential helper
              sudo apt-get -qq --yes --force-yes install amazon-ecr-credential-helper
              sudo mkdir -p /home/ubuntu/.docker
              sudo mkdir -p /home/bestiary/.docker
              sudo bash -c 'echo "{\"credsStore\": \"ecr-login\"}" > /home/ubuntu/.docker/config.json'
              sudo bash -c 'echo "{\"credsStore\": \"ecr-login\"}" > /home/bestiary/.docker/config.json'
              sudo chown -R ubuntu:ubuntu /home/ubuntu
              sudo chown -R bestiary:bestiary /home/bestiary
              EOF
}

# Elastic IP

resource "aws_eip" "app_ip" {
  instance = aws_instance.app.id
}

# Instance IP

output "instance_ip" {
  value = aws_eip.app_ip.public_ip
}

resource "aws_ecr_repository" "bestiary_repo" {
  name = "bestiary-repo"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "bestiary-repo"
  }
}

output "ecr_repository_url" {
  value = aws_ecr_repository.bestiary_repo.repository_url
}