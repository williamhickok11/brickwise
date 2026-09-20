export interface Mortgage {
  id: number
  property_id: number
  original_loan: number
  current_balance: number
  interest_rate: number
  monthly_payment: number
  monthly_pi: number
  balance_as_of: string | null
  lender: string
  notes: string
  created_at: string
  updated_at: string
}

export interface UpsertMortgageRequest {
  original_loan: number
  current_balance: number
  interest_rate: number
  monthly_payment: number
  monthly_pi: number
  balance_as_of?: string | null
  lender: string
  notes: string
}
