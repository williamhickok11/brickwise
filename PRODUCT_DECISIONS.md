# Product Decisions & UX Guidelines

**Last Updated:** February 6, 2026

This document captures key product decisions, UX philosophy, and architectural choices made during the development of Brickwise. It serves as a living reference for maintaining consistency and understanding the "why" behind our decisions.

---

## Product Vision

**Target User:** Small-to-mid-scale landlords who self-manage and care about cash flow, time, and sanity.

**Core Philosophy:** Optimize for clarity and speed, not flash. Default to "fewer clicks, fewer screens." Assume users are busy, distracted, and managing real money.

**Long-term Goal:** Eventually be better than apartments.com, but start with a solid MVP to build on.

---

## MVP Scope: 4-Tab Navigation

The MVP consists of four core tabs in a mobile-first bottom navigation:

1. **Dashboard** - Cash flow summary, property overview (placeholder for MVP)
2. **Properties** - CRUD functionality for properties
3. **Residences** - Lease/tenant management (placeholder for MVP)
4. **REPS** - Time tracking for Real Estate Professional Status tax compliance

**Decision Rationale:** 
- Mobile-first design because primary use case is on phone
- Bottom tabs provide thumb-friendly navigation
- REPS is a core feature, not an afterthought
- Placeholder tabs establish structure for future expansion

---

## REPS Time Tracking Feature

### The Problem
- Daily logging is too much friction
- Weekly logging loses accuracy
- Need REPS-compliant documentation for tax filing
- Current Excel/Numbers workflow is manual and error-prone

### The Solution
**Mobile-first time entry with voice input support**

**Key Decisions:**

1. **Single Entry Workflow** (not bulk add)
   - "Save & Add Another" button keeps property/date context
   - Easier to maintain consistency
   - Faster than switching between modes

2. **Property-Based Organization**
   - One consolidated log with property filter (not separate sheets)
   - Simpler data model
   - Still REPS-compliant (structure doesn't matter for tax filing)
   - Easier to search/filter across all activities

3. **Voice Input for Descriptions**
   - Web Speech API for description field
   - Structured fields (date, property, category, hours) use quick taps
   - Balances speed with accuracy

4. **Mileage Tracking**
   - "Full Drive" checkbox/toggle
   - Auto-fills mileage from property's `default_mileage` when checked
   - Manual override allowed
   - Example: WoodBridge Dr = 22 miles round trip (stored as default)

5. **Export Format**
   - CSV export matching Numbers format
   - Columns: Date, Activity Category, Description, Time Spent (hrs), Notes, Mileage
   - Filterable by property, date range, category

### Categories (Hardcoded for MVP)
- Property Management
- Maintenance & Repairs
- Contractor Oversight
- Accounting & Admin
- Deal Sourcing
- Construction Oversight
- Software Management

---

## Data Model Decisions

### Properties
- Added `default_mileage` field to store round-trip distance per property
- Used for auto-filling mileage when "Full Drive" is checked

### Time Entries
- `property_id` is nullable (NULL = "General" activities)
- Property-based, not Residence-based (work is on the property regardless of current tenants)
- Residences will be for lease/tenant management in future features

### Residences vs Properties
- **Residences** = Tenants/Leases (people you create leases with)
- A Property can have multiple Residences over time (different tenants, same property)
- REPS time tracking is Property-based (work is on the property, not tied to specific tenants/leases)
- Residences tab will be for lease/tenant management in future features

---

## UX Philosophy & Guardrails

### Core Principles
1. **Clarity over flash** - Optimize for clarity and speed, not visual effects
2. **Fewer clicks, fewer screens** - Default to streamlined workflows
3. **Mobile-first** - Primary use case is on phone, optimize for thumb navigation
4. **Workflows over dashboards** - Prefer actionable workflows to passive displays

### Guardrails
- ❌ Do not invent features for large property management companies
- ❌ Avoid VC-bait features unless explicitly asked
- ✅ Prefer workflows over dashboards
- ✅ Balance power with simplicity (avoid "enterprise PM bloat")

### Mobile-First Design Decisions
- Bottom tab navigation (thumb-friendly)
- Large touch targets (min 44px)
- Sticky action buttons at bottom
- Voice input for free-form text
- Quick taps for structured data
- Number pad for hours input (mobile keyboard)

---

## Architecture Decisions

### Navigation
- **Mobile:** Bottom tab bar (sticky, always visible)
- **Desktop:** Can adapt later (currently hidden on desktop via media query)
- Replaced desktop sidebar with mobile-first approach

### Backend
- Go REST API with clean layered architecture
- SQLite for dev, Postgres for production
- Database abstraction handles both SQLite and Postgres differences

### Frontend
- Vue 3 + TypeScript + Vite
- Pinia for state management
- Mobile-first responsive design
- Web Speech API for voice input (with fallback)

---

## Feature Prioritization Framework

When proposing features, include:
1. **The landlord pain it solves**
2. **Rough UX shape** (not visual design)
3. **Whether it's v1, v2, or "nice to have"**

### REPS Feature Classification
- **v1 MVP:** Basic time entry, property-based organization, voice input, CSV export
- **v2:** Reports, charts, activity summaries, Excel export with formatting
- **Nice to have:** Natural language parsing, recurring task templates, mobile app

---

## Future Considerations

### Dashboard Tab
- **Future:** Cash flow summary, property overview, key metrics
- **MVP:** Placeholder only

### Residences Tab
- **Future:** Lease/tenant management, lease history, tenant contact info
- **MVP:** Placeholder only
- **Note:** Residences are people you create leases with. Same property can have different tenants over time.

### Properties Tab
- **Future:** Add "Quick Add REPS Entry" button from property detail
- **Current:** Full CRUD functionality

---

## Key Workflows

### Daily REPS Logging
1. Open REPS tab
2. Tap "Add Entry"
3. Select property (or "General")
4. Tap category button
5. Enter hours (number pad)
6. Voice input description OR type manually
7. Optional: Check "Full Drive" (auto-fills mileage)
8. Tap "Save & Add Another" (keeps property/date)
9. Repeat for multiple entries
10. Tap "Save & Done" when finished

### Weekly Catch-Up
1. Open REPS tab
2. Set date filters for past week
3. Use "Save & Add Another" to batch entries
4. Property/date context preserved between entries

### Tax Filing Export
1. Open REPS tab
2. Apply filters (property, date range, category) if needed
3. Tap "Export"
4. CSV downloads matching Numbers format
5. Import into Numbers/Excel for tax filing

---

## Technical Notes

### Voice Input
- Uses Web Speech API (`webkitSpeechRecognition` or `SpeechRecognition`)
- Fallback to manual typing if voice not available
- Visual feedback during recording
- Only used for description field (structured fields use taps)

### Export Format
- CSV format (not Excel) for MVP
- Matches Numbers column structure
- Can be imported directly into Numbers/Excel
- Future: Consider Excel formatting for v2

### Database Migrations
- Schema changes handled in `backend/internal/db/db.go`
- Supports both SQLite (dev) and Postgres (production)
- `default_mileage` column added to existing `properties` table
- `time_entries` table created with foreign key to `properties`

---

## Decision Log

### 2026-02-06: Initial MVP Decisions
- ✅ Mobile-first bottom tab navigation
- ✅ REPS as core feature (not add-on)
- ✅ Single entry workflow with "Save & Add Another"
- ✅ Property-based organization (not separate sheets)
- ✅ Voice input for descriptions
- ✅ Mileage tracking with "Full Drive" auto-fill
- ✅ CSV export matching Numbers format
- ✅ 4-tab structure: Dashboard, Properties, Residences, REPS

---

## Questions for Future Consideration

1. Should we add property photos/attachments?
2. Should we add recurring task templates?
3. Should we add natural language parsing for voice input?
4. Should we add charts/reports for REPS hours tracking?
5. Should we add tenant/lease management features?
6. Should we add cash flow tracking?
7. Should we add maintenance request tracking?
8. Should we add document storage?

**Note:** These are future considerations. Focus on MVP first.

---

## How to Update This Document

When making significant product decisions:
1. Add entry to "Decision Log" with date
2. Update relevant sections if decisions change
3. Add new questions to "Questions for Future Consideration"
4. Update "Last Updated" date at top

This document should reflect the "why" behind decisions, not just the "what."
