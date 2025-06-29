package main
/*
# Login and store cookie
curl -X POST http://localhost:8080/login \
   -H "Content-Type: application/json" \
   -c cookies.txt \
   -d '{"username":"mahindra", "password":"123456"}'

# Use stored cookie to access profile
curl -X GET http://localhost:8080/profile \
   -b cookies.txt
*/
