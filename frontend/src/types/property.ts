export interface Property {
  id: number
  name: string
  address: string
  property_type: string
  default_mileage: number
  created_at: string
  updated_at: string
}

export interface CreatePropertyRequest {
  name: string
  address: string
  property_type: string
  default_mileage?: number
}
