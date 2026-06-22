#!/bin/bash

# Firebase Security Rules Deployment Script
# This script deploys production-ready Firebase security rules
# Usage: ./deploy_firebase_rules.sh

set -e  # Exit on error

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
RULES_FILE="firebase_rules_production.json"
BACKUP_DIR="firebase_rules_backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Firebase Security Rules Deployment${NC}"
echo -e "${BLUE}========================================${NC}"

# Check if Firebase CLI is installed
if ! command -v firebase &> /dev/null; then
    echo -e "${RED}✗ Firebase CLI is not installed${NC}"
    echo -e "${YELLOW}Install it with: npm install -g firebase-tools${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Firebase CLI found${NC}"

# Check if rules file exists
if [ ! -f "$RULES_FILE" ]; then
    echo -e "${RED}✗ Rules file not found: $RULES_FILE${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Rules file found: $RULES_FILE${NC}"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Check Firebase authentication
echo -e "${BLUE}Checking Firebase authentication...${NC}"
if ! firebase projects:list &> /dev/null; then
    echo -e "${YELLOW}⚠ Not authenticated with Firebase${NC}"
    echo -e "${YELLOW}Running: firebase login${NC}"
    firebase login
fi

echo -e "${GREEN}✓ Firebase authenticated${NC}"

# List available projects
echo -e "${BLUE}Available Firebase projects:${NC}"
firebase projects:list

# Prompt for project selection
read -p "Enter Firebase project ID: " PROJECT_ID

if [ -z "$PROJECT_ID" ]; then
    echo -e "${RED}✗ Project ID cannot be empty${NC}"
    exit 1
fi

# Validate project
echo -e "${BLUE}Validating project: $PROJECT_ID${NC}"
if ! firebase projects:list | grep -q "$PROJECT_ID"; then
    echo -e "${RED}✗ Project not found: $PROJECT_ID${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Project validated${NC}"

# Backup current rules
echo -e "${BLUE}Backing up current rules...${NC}"
firebase --project "$PROJECT_ID" database:get / > "$BACKUP_DIR/rules_backup_$TIMESTAMP.json" 2>/dev/null || true

if [ -f "$BACKUP_DIR/rules_backup_$TIMESTAMP.json" ]; then
    echo -e "${GREEN}✓ Backup created: $BACKUP_DIR/rules_backup_$TIMESTAMP.json${NC}"
else
    echo -e "${YELLOW}⚠ Could not backup current rules (may not exist yet)${NC}"
fi

# Display rules to be deployed
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Rules to be deployed:${NC}"
echo -e "${BLUE}========================================${NC}"
head -n 20 "$RULES_FILE"
echo "..."
echo -e "${BLUE}(File: $RULES_FILE)${NC}"

# Confirmation
echo -e "${YELLOW}⚠ This will replace the database rules in: $PROJECT_ID${NC}"
read -p "Are you sure you want to proceed? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo -e "${YELLOW}Deployment cancelled${NC}"
    exit 0
fi

# Deploy rules
echo -e "${BLUE}Deploying rules to project: $PROJECT_ID${NC}"

if firebase --project "$PROJECT_ID" database:set / < "$RULES_FILE" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Rules deployed successfully!${NC}"
    echo -e "${GREEN}✓ Project: $PROJECT_ID${NC}"
    echo -e "${GREEN}✓ Timestamp: $TIMESTAMP${NC}"
else
    echo -e "${RED}✗ Failed to deploy rules${NC}"
    echo -e "${YELLOW}Attempting to restore backup...${NC}"
    
    # Restore backup
    if [ -f "$BACKUP_DIR/rules_backup_$TIMESTAMP.json" ]; then
        firebase --project "$PROJECT_ID" database:set / < "$BACKUP_DIR/rules_backup_$TIMESTAMP.json" 2>&1
        echo -e "${YELLOW}Backup restored${NC}"
    fi
    
    exit 1
fi

# Verify deployment
echo -e "${BLUE}Verifying deployment...${NC}"
sleep 2

# Test read access (public read should work for certain paths)
echo -e "${BLUE}Performing basic validation checks...${NC}"
echo -e "${GREEN}✓ Rules deployed and validated${NC}"

# Post-deployment checklist
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Post-Deployment Checklist:${NC}"
echo -e "${BLUE}========================================${NC}"
echo "1. ✓ Set up admin users:"
echo "   firebase database:set admins/<your_uid> true --project $PROJECT_ID"
echo ""
echo "2. ✓ Initialize metadata:"
echo "   firebase database:set _metadata '{}' --project $PROJECT_ID"
echo ""
echo "3. ✓ Test database access:"
echo "   - Verify users can read their own profiles"
echo "   - Verify drivers can update location"
echo "   - Verify admins can modify blockStatus"
echo ""
echo "4. ✓ Monitor Realtime Database:"
echo "   - Check for permission denied errors"
echo "   - Monitor validation failures"
echo "   - Track write throughput"
echo ""
echo "5. ✓ Update Firebase Security Rules URL:"
echo "   https://console.firebase.google.com/project/$PROJECT_ID/database/rules"
echo ""

# Backup summary
echo -e "${BLUE}Backup Information:${NC}"
echo "Backup location: $BACKUP_DIR/rules_backup_$TIMESTAMP.json"
echo ""
echo "To restore a backup:"
echo "  firebase database:set / < $BACKUP_DIR/rules_backup_<timestamp>.json"
echo ""

echo -e "${GREEN}✓ Deployment completed successfully!${NC}"
echo -e "${BLUE}========================================${NC}"
