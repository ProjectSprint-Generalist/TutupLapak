#!/bin/bash

# TutupLapak API Docker Quick Start Script

set -e

echo "🐳 TutupLapak API Docker Setup"
echo "================================"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "✅ .env file created. You can edit it if needed."
fi

# Function to show help
show_help() {
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  dev     Start development environment (with admin tools)"
    echo "  prod    Start production environment"
    echo "  stop    Stop all services"
    echo "  logs    Show logs"
    echo "  status  Show service status"
    echo "  clean   Clean up everything"
    echo "  help    Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 dev     # Start development environment"
    echo "  $0 prod    # Start production environment"
    echo "  $0 stop    # Stop all services"
}

# Function to start development environment
start_dev() {
    echo "🚀 Starting development environment..."
    docker-compose -f docker-compose.dev.yml up -d
    
    echo ""
    echo "✅ Development environment started!"
    echo ""
    echo "🌐 Service URLs:"
    echo "  API:              http://localhost:8080"
    echo "  API Health:       http://localhost:8080/v1/health"
    echo "  MinIO Console:    http://localhost:9001 (admin/minioadmin123)"
    echo "  pgAdmin:          http://localhost:5050 (admin@tutuplapak.com/admin123)"
    echo "  Redis Commander:  http://localhost:8081"
    echo ""
    echo "📊 To view logs: $0 logs"
    echo "🛑 To stop: $0 stop"
}

# Function to start production environment
start_prod() {
    echo "🚀 Starting production environment..."
    docker-compose -f docker-compose.prod.yml up -d
    
    echo ""
    echo "✅ Production environment started!"
    echo ""
    echo "🌐 Service URLs:"
    echo "  API:              http://localhost:8080"
    echo "  API Health:       http://localhost:8080/v1/health"
    echo ""
    echo "📊 To view logs: $0 logs"
    echo "🛑 To stop: $0 stop"
}

# Function to stop services
stop_services() {
    echo "🛑 Stopping all services..."
    docker-compose down
    docker-compose -f docker-compose.dev.yml down
    docker-compose -f docker-compose.prod.yml down
    echo "✅ All services stopped."
}

# Function to show logs
show_logs() {
    echo "📊 Showing logs (Press Ctrl+C to exit)..."
    docker-compose logs -f
}

# Function to show status
show_status() {
    echo "📊 Service Status:"
    docker-compose ps
}

# Function to clean up
clean_up() {
    echo "🧹 Cleaning up Docker resources..."
    docker-compose down -v
    docker-compose -f docker-compose.dev.yml down -v
    docker-compose -f docker-compose.prod.yml down -v
    docker system prune -f
    docker volume prune -f
    echo "✅ Cleanup completed."
}

# Main script logic
case "${1:-dev}" in
    "dev")
        start_dev
        ;;
    "prod")
        start_prod
        ;;
    "stop")
        stop_services
        ;;
    "logs")
        show_logs
        ;;
    "status")
        show_status
        ;;
    "clean")
        clean_up
        ;;
    "help"|"-h"|"--help")
        show_help
        ;;
    *)
        echo "❌ Unknown option: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
