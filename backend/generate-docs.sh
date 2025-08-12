#!/bin/bash

# Generate Swagger documentation
echo "🔄 Generating Swagger documentation..."
swag init -g cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation generated successfully!"
    echo "📁 Files generated:"
    echo "   - docs/docs.go"
    echo "   - docs/swagger.json" 
    echo "   - docs/swagger.yaml"
    echo ""
    echo "🚀 Start the server and visit http://localhost:8080/swagger/ to view the API documentation"
else
    echo "❌ Failed to generate Swagger documentation"
    exit 1
fi
