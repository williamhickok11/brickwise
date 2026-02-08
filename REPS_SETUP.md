# REPS Feature - Data & Persistence Setup

## ✅ What's Already Implemented

### Backend (Go)
- ✅ **Database Schema**: `time_entries` table created with all required fields
- ✅ **Models**: `TimeEntry`, `CreateTimeEntryRequest`, `UpdateTimeEntryRequest`, `TimeEntryFilter`
- ✅ **Repository Layer**: Full CRUD operations with SQLite/Postgres support
- ✅ **Service Layer**: Business logic, validation, auto-mileage calculation
- ✅ **Handlers**: All HTTP endpoints implemented
- ✅ **Routes**: Registered in `/api/v1/time-entries`
- ✅ **Validation**: Category, date, hours, description validation
- ✅ **Date Formatting**: Custom `Date` type ensures YYYY-MM-DD format in JSON

### Frontend (Vue/TypeScript)
- ✅ **Store**: Pinia store with all CRUD operations
- ✅ **API Client**: Axios client configured for all endpoints
- ✅ **Types**: TypeScript interfaces match backend models
- ✅ **UI Components**: REPSView and REPSEntryForm components
- ✅ **Filtering**: Property, category, date range filters

## 🔗 API Endpoints

All endpoints are available at `/api/v1/time-entries`:

- `GET /api/v1/time-entries` - List entries (with optional filters)
- `GET /api/v1/time-entries/{id}` - Get single entry
- `POST /api/v1/time-entries` - Create entry
- `PUT /api/v1/time-entries/{id}` - Update entry
- `DELETE /api/v1/time-entries/{id}` - Delete entry
- `GET /api/v1/time-entries/export` - Export to CSV

## 📋 Data Model

### TimeEntry Fields
- `id` (int) - Auto-generated
- `property_id` (int | null) - NULL for "General" activities
- `date` (string) - YYYY-MM-DD format
- `category` (string) - One of the valid categories
- `description` (string) - Required, max 2000 chars
- `hours` (float64) - Required, 0 < hours <= 24
- `notes` (string) - Optional
- `mileage` (float64) - Default 0
- `full_drive` (bool) - Auto-calculates mileage if true and property_id set
- `created_at` (timestamp)
- `updated_at` (timestamp)

### Valid Categories
- Property Management
- Maintenance & Repairs
- Contractor Oversight
- Accounting & Admin
- Deal Sourcing
- Construction Oversight
- Software Management

## 🔧 Key Features

1. **Auto-Mileage Calculation**: When `full_drive` is true and a property is selected, mileage is automatically set from the property's `default_mileage` field.

2. **Property Filtering**: 
   - `property_id: null` or `property_id: 0` = "General" activities
   - Specific property ID = Filter to that property
   - No filter = All entries

3. **Date Formatting**: Dates are serialized as "YYYY-MM-DD" strings in JSON responses (not RFC3339).

## ✅ What's Working

The entire data flow is connected:
1. Frontend store → API client → Backend handlers
2. Backend handlers → Service → Repository → Database
3. Database → Repository → Service → Handlers → JSON → Frontend

## 🧪 Testing Checklist

To verify everything works:

1. **Create Entry**
   ```bash
   curl -X POST http://localhost:8080/api/v1/time-entries \
     -H "Content-Type: application/json" \
     -d '{
       "property_id": null,
       "date": "2024-01-15",
       "category": "Property Management",
       "description": "Test entry",
       "hours": 2.5,
       "notes": "",
       "mileage": 0,
       "full_drive": false
     }'
   ```

2. **List Entries**
   ```bash
   curl http://localhost:8080/api/v1/time-entries
   ```

3. **Filter by Property**
   ```bash
   curl "http://localhost:8080/api/v1/time-entries?property_id=null"
   ```

4. **Update Entry**
   ```bash
   curl -X PUT http://localhost:8080/api/v1/time-entries/1 \
     -H "Content-Type: application/json" \
     -d '{"property_id": null, "date": "2024-01-15", "category": "Property Management", "description": "Updated", "hours": 3.0, "notes": "", "mileage": 0, "full_drive": false}'
   ```

5. **Export CSV**
   ```bash
   curl "http://localhost:8080/api/v1/time-entries/export" -o entries.csv
   ```

## 🚀 Next Steps

1. **Start the backend server**:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. **Start the frontend** (if not already running):
   ```bash
   cd frontend
   npm run dev
   ```

3. **Test the UI**: Navigate to `/reps` route and try:
   - Creating a new entry
   - Filtering by property/category/date
   - Editing an entry
   - Deleting an entry
   - Exporting to CSV

## 📝 Notes

- The database will be created automatically at `./brickwise.db` (SQLite) or use Postgres if `DATABASE_URL` is set
- The `time_entries` table is created automatically on server startup
- Property relationships: `property_id` can be NULL for "General" activities
- Date handling: All dates are stored as DATE type and serialized as "YYYY-MM-DD" strings

## 🐛 Potential Issues to Watch

1. **Date Format**: Dates are now serialized as "YYYY-MM-DD" - verify frontend handles this correctly
2. **Property ID Null Handling**: When `property_id` is `null` in JSON, it should work. When filtering, use `property_id=null` in query string
3. **Boolean Handling**: SQLite stores booleans as integers (0/1), Postgres as BOOLEAN - both are handled correctly
