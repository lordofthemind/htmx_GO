#!/bin/bash

curl -X POST http://localhost:9090/superuser/register \
-H "Content-Type: application/json" \
-H "Accept: application/json" \
-d '{
  "username": "testuser",
  "email": "testuser@example.com",
  "password": "password123"
}'
