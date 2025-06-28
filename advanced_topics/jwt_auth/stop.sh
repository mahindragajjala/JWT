#!/bin/bash

# Define the service names (used during build in run.sh)
SERVICES=(
    "auth_service"
    "user_service"
    "cart_service"
    "order_service"
    "payment_service"
)

echo "🛑 Stopping all services..."

for SERVICE in "${SERVICES[@]}"; do
    echo "🔍 Searching for $SERVICE..."

    # Find and kill the running process for each service
    PID=$(pgrep -f "./$SERVICE/service")

    if [ -z "$PID" ]; then
        echo "⚠️  $SERVICE not running."
    else
        kill "$PID"
        echo "✅ $SERVICE (PID $PID) stopped."
    fi
done

echo "🧹 Cleanup done. All specified services stopped."

