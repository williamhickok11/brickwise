export interface TimeEntry {
  id: number
  property_id: number | null
  date: string
  category: string
  description: string
  hours: number
  notes: string
  mileage: number
  full_drive: boolean
  created_at: string
  updated_at: string
}

export interface CreateTimeEntryRequest {
  property_id: number | null
  date: string
  category: string
  description: string
  hours: number
  notes: string
  mileage: number
  full_drive: boolean
}

export interface UpdateTimeEntryRequest {
  property_id: number | null
  date: string
  category: string
  description: string
  hours: number
  notes: string
  mileage: number
  full_drive: boolean
}

export interface TimeEntryFilter {
  property_id?: number | null
  start_date?: string
  end_date?: string
  category?: string
}

export const CATEGORIES = [
  'Property Management',
  'Maintenance & Repairs',
  'Contractor Oversight',
  'Accounting & Admin',
  'Deal Sourcing',
  'Construction Oversight',
  'Software Management',
] as const

export type Category = typeof CATEGORIES[number]
