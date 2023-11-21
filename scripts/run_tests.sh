#!/bin/bash

go test ./... -coverprofile=coverage.out

go tool cover -html=coverage.out -o coverage.html

echo "Generated coverage report at ./coverage.html"
