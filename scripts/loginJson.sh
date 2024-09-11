#!/bin/bash

curl -X POST http://localhost:9090/superuser/login \
-H "Content-Type: application/json" \
-H "Accept: application/json" \
-d '{
  "email": "testuser@example.com",
  "password": "password123"
}'
