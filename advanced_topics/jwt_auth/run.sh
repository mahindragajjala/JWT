#!/bin/bash

# Define services and their ports
declare -A SERVICES
SERVICES=(
    ["auth_service"]=8000
    ["user_service"]=8001
    ["cart_service"]=8002
    ["order_service"]=8003
    ["payment_service"]=8004
)

# Build and run each service
for SERVICE in "${!SERVICES[@]}"; do
    PORT=${SERVICES[$SERVICE]}

    echo "🔧 Building $SERVICE..."
    go build -o "$SERVICE/service" "$SERVICE/main.go"

    if [ $? -ne 0 ]; then
        echo "❌ Build failed for $SERVICE, skipping..."
        continue
    fi

    echo "🚀 Running $SERVICE on port $PORT..."
    PORT=$PORT "./$SERVICE/service" > "$SERVICE/output.log" 2>&1 &

    echo "✅ $SERVICE started on port $PORT (log: $SERVICE/output.log)"
done

echo "🎉 All services are running in the background!"

