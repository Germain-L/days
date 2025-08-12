#!/bin/bash

# Test runner script for the Days backend

set -e

echo "🧪 Running Days Backend Test Suite"
echo "=================================="

# Set test environment
export GO_ENV=test
export JWT_SECRET=test-secret-for-all-tests

echo
echo "📦 Installing dependencies..."
go mod tidy

echo
echo "🔍 Running linting checks..."
go vet ./...

echo
echo "🏗️ Building all packages..."
go build ./...

echo
echo "🧪 Running unit tests..."
echo

echo "  🔐 Testing JWT auth module..."
go test -v ./internal/auth

echo
echo "  📊 Testing services..."
go test -v ./internal/services

echo
echo "  🌐 Testing handlers..."
go test -v ./internal/handlers

echo
echo "🚀 Running all tests with coverage..."
go test -v -race -coverprofile=coverage.out ./...

echo
echo "📊 Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo
echo "📈 Coverage summary:"
go tool cover -func=coverage.out | tail -1

echo
echo "✅ All tests completed successfully!"
echo
echo "📊 Coverage report generated: coverage.html"
echo "🎯 To run integration tests: go test -v ./... -tags=integration"
echo "⚡ To run only unit tests: go test -short ./..."
