export interface Residence {
  id: number
  property_id: number
  name: string
  phone: string
  email: string
  start_date: string
  end_date: string | null
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateResidenceRequest {
  property_id: number
  name: string
  phone: string
  email: string
  start_date: string
  end_date?: string | null
  is_active: boolean
}

export interface UpdateResidenceRequest {
  property_id: number
  name: string
  phone: string
  email: string
  start_date: string
  end_date?: string | null
  is_active: boolean
}

export interface ResidenceFilter {
  property_id?: number
  is_active?: boolean
}

