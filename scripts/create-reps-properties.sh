#!/bin/bash
# Quick script to create properties for REPS V2
# Make sure your backend server is running on http://localhost:8080

API_URL="${API_URL:-http://localhost:8080/api/v1/properties}"

echo "Creating properties for REPS V2..."
echo ""

# Create 3522 Wood Bridge Dr
echo "Creating: 3522 Wood Bridge Dr"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "3522 Wood Bridge Dr",
    "address": "3522 Wood Bridge Dr",
    "property_type": "residential",
    "default_mileage": 0
  }' | jq '.' 2>/dev/null || echo "Created (response not JSON formatted)"
echo ""

# Create 3500 Saindon St
echo "Creating: 3500 Saindon St"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "3500 Saindon St",
    "address": "3500 Saindon St",
    "property_type": "residential",
    "default_mileage": 0
  }' | jq '.' 2>/dev/null || echo "Created (response not JSON formatted)"
echo ""

echo "✓ Properties created! They should now appear in the REPS V2 Property dropdown."
echo ""
echo "To verify, run: curl $API_URL"
