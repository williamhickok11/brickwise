# Brickwise

Landlord/investor property management application.

> 📋 **Product Decisions:** See [PRODUCT_DECISIONS.md](./PRODUCT_DECISIONS.md) for UX philosophy, feature decisions, and architectural choices.

## Stack

- **Frontend**: Vue 3 + TypeScript + Vite
- **Backend**: Go REST API
- **Database**: SQLite (dev) → PostgreSQL (production)

## Quick Start

### Prerequisites

- Node.js 18+ and npm/yarn/pnpm
- Go 1.21+

### Development

**Backend:**
```bash
cd backend
go mod download
go run cmd/server/main.go
```
Server runs on `http://localhost:8080`

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```
Frontend runs on `http://localhost:5173`

## Project Structure

```
├── frontend/          # Vue 3 + TypeScript app
│   ├── src/
│   │   ├── api/      # API client
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/   # Pinia stores
│   │   └── types/    # TypeScript types
│   └── ...
├── backend/           # Go REST API
│   ├── cmd/          # Application entrypoints
│   ├── internal/     # Private application code
│   │   ├── handlers/ # HTTP handlers
│   │   ├── models/   # Data models
│   │   ├── db/       # Database layer
│   │   └── middleware/
│   └── ...
└── README.md
```

## Architecture Notes

- **Simple & Boring**: Standard patterns, minimal abstractions
- **Fast Iteration**: Hot reload on both sides, minimal build steps
- **Clean Seams**: Easy to swap components (DB, auth, etc.) later
