/** Keep in sync with backend ValidateProperty validTypes. */
export type PropertyType = 'residential' | 'commercial' | 'industrial' | 'land' | 'mixed'

export const PROPERTY_TYPES: { label: string; value: PropertyType }[] = [
  { label: 'Residential', value: 'residential' },
  { label: 'Commercial', value: 'commercial' },
  { label: 'Industrial', value: 'industrial' },
  { label: 'Land', value: 'land' },
  { label: 'Mixed use', value: 'mixed' },
]

export interface Property {
  id: number
  name: string
  address: string
  property_type: PropertyType | string
  default_mileage: number
  current_value: number
  value_updated_at: string | null
  created_at: string
  updated_at: string
}

export interface CreatePropertyRequest {
  name: string
  address: string
  property_type: PropertyType | ''
  default_mileage?: number
}

export interface UpdatePropertyValueRequest {
  current_value: number
  value_updated_at?: string | null
}
