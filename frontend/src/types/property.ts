export interface Property {
  id: number
  name: string
  address: string
  property_type: string
  default_mileage: number
  current_value: number
  value_updated_at: string | null
  created_at: string
  updated_at: string
}

export interface CreatePropertyRequest {
  name: string
  address: string
  property_type: string
  default_mileage?: number
}

export interface UpdatePropertyValueRequest {
  current_value: number
  value_updated_at?: string | null
}
