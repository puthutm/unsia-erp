#!/bin/bash
# ==========================================
# UNSIA ERP - Jalankan di WSL
# ==========================================

# Configure UTF-8
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8

echo "=========================================="
echo "  UNSIA ERP - Running on WSL"
echo "=========================================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker belum terinstall!"
    echo "   Install Docker Desktop for Windows dengan WSL2 integration"
    echo "  atau jalankan: wsl --install"
    exit 1
fi

# Check Docker daemon
if ! docker info &> /dev/null; then
    echo "❌ Docker daemon tidak berjalan!"
    echo "   Start Docker Desktop"
    exit 1
fi

# Check docker-compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose belum terinstall!"
    exit 1
fi

echo "✅ Docker: $(docker --version)"
echo "✅ Docker Compose: $(docker-compose --version)"
echo ""

# Navigate to project directory
cd "$(dirname "$0")" || exit 1

# Build and start services
echo "=========================================="
echo "  Building & Starting Services..."
echo "=========================================="
echo ""

# Pull base images first
echo "📦 Pulling base images..."
docker-compose pull postgres redis rabbitmq || true

# Build all services
echo "🔨 Building services..."
docker-compose build

# Start infrastructure first
echo "🚀 Starting infrastructure services..."
docker-compose up -d postgres redis rabbitmq

# Wait for PostgreSQL
echo "⏳ Waiting for PostgreSQL..."
for i in {1..30}; do
    if docker-compose exec -T postgres pg_isready -U postgres &> /dev/null; then
        echo "✅ PostgreSQL ready!"
        break
    fi
    sleep 1
done

# Wait for RabbitMQ
echo "⏳ Waiting for RabbitMQ..."
for i in {1..30}; do
    if docker-compose exec -T rabbitmq rabbitmq-diagnostics -q ping &> /dev/null; then
        echo "✅ RabbitMQ ready!"
        break
    fi
    sleep 1
done

# Start all services
echo "🚀 Starting all services..."
docker-compose up -d

echo ""
echo "=========================================="
echo "  Services Started!"
echo "=========================================="
echo ""
echo "📋 Service URLs:"
echo "   - Core Service:        http://localhost:8001"
echo "   - Reference Service:  http://localhost:8002"
echo "   - PMB Service:        http://localhost:8003"
echo "   - Academic Service:  http://localhost:8004"
echo "   - Finance Service:   http://localhost:8005"
echo "   - LMS Service:       http://localhost:8006"
echo "   - Assessment Svc:    http://localhost:8007"
echo "   - HRIS Service:      http://localhost:8008"
echo "   - CRM Service:       http://localhost:8009"
echo "   - Portal Service:    http://localhost:8010"
echo "   - Frontend Web:      http://localhost:3000"
echo ""
echo "   - RabbitMQ UI:       http://localhost:15672"
echo "                      User: unsia / Pass: unsia"
echo ""
echo "=========================================="
echo "  Useful Commands:"
echo "=========================================="
echo "   View logs:    docker-compose logs -f"
echo "   Stop:        docker-compose down"
echo "   Restart:     docker-compose restart"
echo "   Status:      docker-compose ps"
echo ""
