#!/bin/bash

# Set proper folder names (based on your current structure)
AUTH_DIR="auth_service"
USER_DIR="user_service"
CART_DIR="cart_service"
ORDER_DIR="order_service"
PAYMENT_DIR="payment_service"

# Ensure all directories exist
mkdir -p "$AUTH_DIR" "$USER_DIR" "$CART_DIR" "$ORDER_DIR" "$PAYMENT_DIR"

echo "🔐 Generating private key for $AUTH_DIR..."
openssl genrsa -out "$AUTH_DIR/private.key" 2048

echo "🔓 Generating public key from private key..."
openssl rsa -in "$AUTH_DIR/private.key" -pubout -out "$USER_DIR/public.key"

echo "📤 Distributing public key to other services..."
cp "$USER_DIR/public.key" "$CART_DIR/"
cp "$USER_DIR/public.key" "$ORDER_DIR/"
cp "$USER_DIR/public.key" "$PAYMENT_DIR/"

echo "✅ Keys generated and distributed successfully!"

