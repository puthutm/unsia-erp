#!/bin/bash

# ==========================================
# UNSIA ERP - Run All Services Script
# Server IP: 10.10.20.56
# ==========================================

SERVER_IP="10.10.20.56"
POSTGRES_HOST="$SERVER_IP"
REDIS_HOST="$SERVER_IP"
RABBITMQ_HOST="$SERVER_IP"

echo "Starting UNSIA ERP Services..."
echo "Server IP: $SERVER_IP"
echo "=========================================="

# ==========================================
# CORE SERVICE - Port 8001
# ==========================================
echo "[1/10] Starting unsia-core-service (port 8001)..."
cd services/unsia-core-service
PORT=8001 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=core_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
RABBITMQ_URL=amqp://unsia:unsia@$RABBITMQ_HOST:5672 \
go run ./cmd/core-service &
CORE_PID=$!
cd ../..
echo "unsia-core-service started (PID: $CORE_PID)"

sleep 2

# ==========================================
# REFERENCE SERVICE - Port 8002
# ==========================================
echo "[2/10] Starting unsia-reference-service (port 8002)..."
cd services/unsia-reference-service
PORT=8002 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=reference_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
CORE_SERVICE_URL=http://$SERVER_IP:8001 \
go run ./cmd/reference-service &
REFERENCE_PID=$!
cd ../..
echo "unsia-reference-service started (PID: $REFERENCE_PID)"

sleep 2

# ==========================================
# PMB SERVICE - Port 8003
# ==========================================
echo "[3/10] Starting unsia-pmb-service (port 8003)..."
cd services/unsia-pmb-service
PORT=8003 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=pmb_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
RABBITMQ_URL=amqp://unsia:unsia@$RABBITMQ_HOST:5672 \
CORE_SERVICE_URL=http://$SERVER_IP:8001 \
REFERENCE_SERVICE_URL=http://$SERVER_IP:8002 \
go run ./cmd/pmb-service &
PMB_PID=$!
cd ../..
echo "unsia-pmb-service started (PID: $PMB_PID)"

sleep 2

# ==========================================
# ACADEMIC SERVICE - Port 8004
# ==========================================
echo "[4/10] Starting unsia-academic-service (port 8004)..."
cd services/unsia-academic-service
PORT=8004 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=academic_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
RABBITMQ_URL=amqp://unsia:unsia@$RABBITMQ_HOST:5672 \
CORE_SERVICE_URL=http://$SERVER_IP:8001 \
REFERENCE_SERVICE_URL=http://$SERVER_IP:8002 \
FINANCE_SERVICE_URL=http://$SERVER_IP:8005 \
go run ./cmd/academic-service &
ACADEMIC_PID=$!
cd ../..
echo "unsia-academic-service started (PID: $ACADEMIC_PID)"

sleep 2

# ==========================================
# FINANCE SERVICE - Port 8005
# ==========================================
echo "[5/10] Starting unsia-finance-service (port 8005)..."
cd services/unsia-finance-service
PORT=8005 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=finance_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
RABBITMQ_URL=amqp://unsia:unsia@$RABBITMQ_HOST:5672 \
CORE_SERVICE_URL=http://$SERVER_IP:8001 \
REFERENCE_SERVICE_URL=http://$SERVER_IP:8002 \
go run ./cmd/finance-service &
FINANCE_PID=$!
cd ../..
echo "unsia-finance-service started (PID: $FINANCE_PID)"

sleep 2

# ==========================================
# LMS SERVICE - Port 8006
# ==========================================
echo "[6/10] Starting unsia-lms-service (port 8006)..."
cd services/unsia-lms-service
PORT=8006 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=lms_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
RABBITMQ_URL=amqp://unsia:unsia@$RABBITMQ_HOST:5672 \
CORE_SERVICE_URL=http://$SERVER_IP:8001 \
REFERENCE_SERVICE_URL=http://$SERVER_IP:8002 \
ACADEMIC_SERVICE_URL=http://$SERVER_IP:8004 \
go run ./cmd/lms-service &
LMS_PID=$!
cd ../..
echo "unsia-lms-service started (PID: $LMS_PID)"

sleep 2

# ==========================================
# ASSESSMENT SERVICE - Port 8007
# ==========================================
echo "[7/10] Starting unsia-assessment-service (port 8007)..."
cd services/unsia-assessment-service
PORT=8007 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=assessment_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
RABBITMQ_URL=amqp://unsia:unsia@$RABBITMQ_HOST:5672 \
CORE_SERVICE_URL=http://$SERVER_IP:8001 \
LMS_SERVICE_URL=http://$SERVER_IP:8006 \
go run ./cmd/assessment-service &
ASSESSMENT_PID=$!
cd ../..
echo "unsia-assessment-service started (PID: $ASSESSMENT_PID)"

sleep 2

# ==========================================
# HRIS SERVICE - Port 8008
# ==========================================
echo "[8/10] Starting unsia-hris-service (port 8008)..."
cd services/unsia-hris-service
PORT=8008 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=hris_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
RABBITMQ_URL=amqp://unsia:unsia@$RABBITMQ_HOST:5672 \
CORE_SERVICE_URL=http://$SERVER_IP:8001 \
REFERENCE_SERVICE_URL=http://$SERVER_IP:8002 \
FINANCE_SERVICE_URL=http://$SERVER_IP:8005 \
go run ./cmd/hris-service &
HRIS_PID=$!
cd ../..
echo "unsia-hris-service started (PID: $HRIS_PID)"

sleep 2

# ==========================================
# CRM SERVICE - Port 8009
# ==========================================
echo "[9/10] Starting unsia-crm-service (port 8009)..."
cd services/unsia-crm-service
PORT=8009 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=crm_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
RABBITMQ_URL=amqp://unsia:unsia@$RABBITMQ_HOST:5672 \
CORE_SERVICE_URL=http://$SERVER_IP:8001 \
REFERENCE_SERVICE_URL=http://$SERVER_IP:8002 \
go run ./cmd/crm-service &
CRM_PID=$!
cd ../..
echo "unsia-crm-service started (PID: $CRM_PID)"

sleep 2

# ==========================================
# PORTAL SERVICE - Port 8010
# ==========================================
echo "[10/10] Starting unsia-portal-service (port 8010)..."
cd services/unsia-portal-service
PORT=8010 \
DB_HOST=$POSTGRES_HOST \
DB_PORT=5432 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=portal_db \
REDIS_URL=redis://$REDIS_HOST:6379 \
RABBITMQ_URL=amqp://unsia:unsia@$RABBITMQ_HOST:5672 \
CORE_SERVICE_URL=http://$SERVER_IP:8001 \
REFERENCE_SERVICE_URL=http://$SERVER_IP:8002 \
go run ./cmd/portal-service &
PORTAL_PID=$!
cd ../..
echo "unsia-portal-service started (PID: $PORTAL_PID)"

echo "=========================================="
echo "All services started!"
echo ""
echo "Service Ports:"
echo "  - Core Service:        8001"
echo "  - Reference Service:   8002"
echo "  - PMB Service:          8003"
echo "  - Academic Service:    8004"
echo "  - Finance Service:     8005"
echo "  - LMS Service:         8006"
echo "  - Assessment Service: 8007"
echo "  - HRIS Service:        8008"
echo "  - CRM Service:         8009"
echo "  - Portal Service:      8010"
echo ""
echo "Infrastructure:"
echo "  - PostgreSQL: $POSTGRES_HOST:5432"
echo "  - Redis:     $REDIS_HOST:6379"
echo "  - RabbitMQ:  $RABBITMQ_HOST:5672"
echo ""
echo "Press Ctrl+C to stop all services"

# Wait for all background processes
wait
