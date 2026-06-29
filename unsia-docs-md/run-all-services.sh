#!/bin/bash

# Script to run all UNSIA ERP services
# Usage: ./run-all-services.sh [IP_ADDRESS]
# Default IP: 10.10.20.56

export IP_ADDRESS=${1:-10.10.20.56}

echo "=========================================="
echo "Starting UNSIA ERP Services"
echo "IP Address: $IP_ADDRESS"
echo "=========================================="

# First generate env files
echo "Generating .env files..."
bash generate-env.sh $IP_ADDRESS

# Function to build and run a service
run_service() {
    local SERVICE_NAME=$1
    local SERVICE_DIR=$2
    local PORT=$3
    
    echo ""
    echo "Building and starting $SERVICE_NAME..."
    docker build -t $SERVICE_NAME $SERVICE_DIR
    docker run -d --name $SERVICE_NAME \
        --network unsia-erp-network \
        -p ${PORT}:${PORT} \
        --env-file $SERVICE_DIR/.env \
        $SERVICE_NAME
}

# Create network if not exists
docker network create unsia-erp-network 2>/dev/null || true

# Core Service (Port 8001)
run_service "unsia-core-service" "./services/unsia-core-service" 8001

# Reference Service (Port 8002)
run_service "unsia-reference-service" "./services/unsia-reference-service" 8002

# PMB Service (Port 8003)
run_service "unsia-pmb-service" "./services/unsia-pmb-service" 8003

# Academic Service (Port 8004)
run_service "unsia-academic-service" "./services/unsia-academic-service" 8004

# Finance Service (Port 8005)
run_service "unsia-finance-service" "./services/unsia-finance-service" 8005

# LMS Service (Port 8006)
run_service "unsia-lms-service" "./services/unsia-lms-service" 8006

# Assessment Service (Port 8007)
run_service "unsia-assessment-service" "./services/unsia-assessment-service" 8007

# HRIS Service (Port 8008)
run_service "unsia-hris-service" "./services/unsia-hris-service" 8008

# CRM Service (Port 8009)
run_service "unsia-crm-service" "./services/unsia-crm-service" 8009

# Portal Service (Port 8010)
run_service "unsia-portal-service" "./services/unsia-portal-service" 8010

echo ""
echo "=========================================="
echo "All services started!"
echo "=========================================="
echo ""
echo "Service Ports:"
echo "  - unsia-core-service:         8001"
echo "  - unsia-reference-service: 8002"
echo "  - unsia-pmb-service:       8003"
echo "  - unsia-academic-service:  8004"
echo "  - unsia-finance-service:  8005"
echo "  - unsia-lms-service:     8006"
echo "  - unsia-assessment-svc:   8007"
echo "  - unsia-hris-service:      8008"
echo "  - unsia-crm-service:       8009"
echo "  - unsia-portal-service:     8010"
echo ""
echo "To view logs: docker logs [SERVICE_NAME]"
echo "To stop: docker stop \$(docker ps -q)"
