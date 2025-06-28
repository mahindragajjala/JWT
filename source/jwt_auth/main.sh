#!/bin/bash

# List of your service directories
SERVICES=("auth_service" "user_service" "cart_service" "order_service" "payment_service")

# Iterate through each service directory
for SERVICE in "${SERVICES[@]}"; do
    if [ -d "$SERVICE" ]; then
        echo "📁 Creating empty main.go in $SERVICE..."

        # Create the file if not exists
        touch "$SERVICE/main.go"

        # Set full read/write/execute permissions (777) for all users
        chmod 777 "$SERVICE/main.go"

        echo "✅ main.go created and chmod 777 set in $SERVICE"
    else
        echo "❌ Directory $SERVICE does not exist, skipping..."
    fi
done

echo "🎉 All done!"

