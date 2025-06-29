package main
/*
# Login and store token
TOKEN=$(curl -s -X POST http://localhost:8081/login \
   -H "Content-Type: application/json" \
   -d '{"username":"mahindra", "password":"123456"}' | jq -r '.token')

# Use the token in Authorization header
curl -X GET http://localhost:8081/profile \
   -H "Authorization: Bearer $TOKEN"
*/
