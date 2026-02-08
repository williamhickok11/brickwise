# Creating Properties for REPS V2

I've created scripts to add the two properties you requested. Here are several ways to create them:

## Option 1: Using the API (Recommended)

Once your backend server is running, use these curl commands:

```bash
# Create 3522 Wood Bridge Dr
curl -X POST http://localhost:8080/api/v1/properties \
  -H "Content-Type: application/json" \
  -d '{
    "name": "3522 Wood Bridge Dr",
    "address": "3522 Wood Bridge Dr",
    "property_type": "residential",
    "default_mileage": 0
  }'

# Create 3500 Saindon St
curl -X POST http://localhost:8080/api/v1/properties \
  -H "Content-Type: application/json" \
  -d '{
    "name": "3500 Saindon St",
    "address": "3500 Saindon St",
    "property_type": "residential",
    "default_mileage": 0
  }'
```

## Option 2: Using SQL Script

If you have SQLite3 installed:

```bash
cd backend
sqlite3 brickwise.db < scripts/create-properties.sql
```

Or for Postgres:

```bash
cd backend
psql $DATABASE_URL -f scripts/create-properties.sql
```

## Option 3: Using the Go Script

```bash
cd backend
go run cmd/create-properties/main.go
```

## Option 4: Via the Frontend UI

1. Navigate to the Properties page in your app
2. Click "Add Property"
3. Fill in:
   - Name: `3522 Wood Bridge Dr`
   - Address: `3522 Wood Bridge Dr`
   - Type: `residential`
   - Default Mileage: `0`
4. Repeat for `3500 Saindon St`

## Verification

After creating the properties, they should automatically appear in the REPS V2 Property dropdown. The frontend fetches properties on mount via `propertyStore.fetchProperties()`.

To verify via API:

```bash
curl http://localhost:8080/api/v1/properties
```

You should see both properties in the response.
