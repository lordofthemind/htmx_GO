#!/bin/bash

curl -X POST http://localhost:9090/superuser/register \
-H "Content-Type: application/json" \
-H "Accept": "text/html" \
-d '{
  "username": "testuser",
  "email": "testuser2@example.com",
  "password": "password123"
}'
