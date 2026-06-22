#!/bin/bash

# FamGo Platform - Service Bootstrap Script
# Creates a new microservice from the template
# Usage: ./scripts/bootstrap-service.sh auth-service "Authentication Service"

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Validate arguments
if [ $# -lt 1 ]; then
    echo -e "${RED}Error: Service name required${NC}"
    echo "Usage: $0 <service-name> [description]"
    echo "Example: $0 auth-service 'Authentication Service'"
    exit 1
fi

SERVICE_NAME=$1
SERVICE_DESC=${2:-"FamGo $SERVICE_NAME"}
SERVICE_DIR="services/$SERVICE_NAME"
TEMPLATE_DIR="services/_template"

# Check if template exists
if [ ! -d "$TEMPLATE_DIR" ]; then
    echo -e "${RED}Error: Template not found at $TEMPLATE_DIR${NC}"
    exit 1
fi

# Check if service already exists
if [ -d "$SERVICE_DIR" ]; then
    echo -e "${YELLOW}Warning: Service directory already exists at $SERVICE_DIR${NC}"
    read -p "Do you want to continue? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}FamGo Service Bootstrap${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}Service Name:${NC} $SERVICE_NAME"
echo -e "${BLUE}Description:${NC} $SERVICE_DESC"
echo -e "${BLUE}Target Directory:${NC} $SERVICE_DIR"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"

# Create service directory
echo -e "\n${YELLOW}[1/5]${NC} Creating service directory..."
mkdir -p "$SERVICE_DIR"
cp -r "$TEMPLATE_DIR"/* "$SERVICE_DIR/"
echo -e "${GREEN}✓ Directory created${NC}"

# Update package.json
echo -e "\n${YELLOW}[2/5]${NC} Updating package.json..."
cd "$SERVICE_DIR"

# Use sed to update name and description in package.json
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/\"name\": \"@famgo\/service-template\"/\"name\": \"@famgo\/$SERVICE_NAME\"/" package.json
    sed -i '' "s/\"description\": \"FamGo Enterprise Service Template\"/\"description\": \"$SERVICE_DESC\"/" package.json
else
    # Linux
    sed -i "s/\"name\": \"@famgo\/service-template\"/\"name\": \"@famgo\/$SERVICE_NAME\"/" package.json
    sed -i "s/\"description\": \"FamGo Enterprise Service Template\"/\"description\": \"$SERVICE_DESC\"/" package.json
fi

echo -e "${GREEN}✓ package.json updated${NC}"

# Create .env file
echo -e "\n${YELLOW}[3/5]${NC} Creating .env file..."
cp .env.example .env
echo -e "${GREEN}✓ .env created (update with your values)${NC}"

# Install dependencies
echo -e "\n${YELLOW}[4/5]${NC} Installing dependencies..."
npm install --silent
echo -e "${GREEN}✓ Dependencies installed${NC}"

# Create directory structure
echo -e "\n${YELLOW}[5/5]${NC} Setting up project structure..."
rm -rf src/modules/example
mkdir -p src/modules/{mymodule,health}

# Create health module
mkdir -p src/modules/health
cat > src/modules/health/health.controller.ts << 'EOF'
import { Controller, Get } from '@nestjs/common';
import { ApiTags, ApiOperation } from '@nestjs/swagger';

@ApiTags('health')
@Controller('health')
export class HealthController {
  @Get()
  @ApiOperation({ summary: 'Health check' })
  check() {
    return {
      status: 'ok',
      timestamp: new Date().toISOString(),
      service: process.env.SERVICE_NAME,
    };
  }
}
EOF

cat > src/modules/health/health.module.ts << 'EOF'
import { Module } from '@nestjs/common';
import { HealthController } from './health.controller';

@Module({
  controllers: [HealthController],
})
export class HealthModule {}
EOF

# Update app.module.ts
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "/import { ExampleModule }/s/ExampleModule/HealthModule/" src/app.module.ts
    sed -i '' "s/ExampleModule/HealthModule/g" src/app.module.ts
else
    sed -i "/import { ExampleModule }/s/ExampleModule/HealthModule/" src/app.module.ts
    sed -i "s/ExampleModule/HealthModule/g" src/app.module.ts
fi

echo -e "${GREEN}✓ Project structure created${NC}"

# Print next steps
echo -e "\n${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}✓ Service created successfully!${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "\n${YELLOW}Next steps:${NC}"
echo -e "  1. cd $SERVICE_DIR"
echo -e "  2. Update .env with your configuration"
echo -e "  3. Create your business modules:"
echo -e "     npm run generate module modules/my-module"
echo -e "  4. Start development server:"
echo -e "     npm run start:dev"
echo -e "  5. View API docs at http://localhost:3000/api/docs"
echo -e "\n${YELLOW}Useful commands:${NC}"
echo -e "  make dev              - Start development server"
echo -e "  make test             - Run tests"
echo -e "  make db-generate      - Generate migration"
echo -e "  make db-migrate       - Run migrations"
echo -e "  make docker-build     - Build Docker image"
echo -e "\n${BLUE}═══════════════════════════════════════════${NC}"

cd - > /dev/null
